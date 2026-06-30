package pqr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/thealanphipps-del/pqr/internal/service"
)

const Version = "v1.08"

type ServerConfig struct {
	GemmaEndpoint    string
	GemmaModel       string
	LMStudioEndpoint string
	LMStudioModel    string
	GeminiAPIKey     string
}

func loadConfig() *ServerConfig {
	cfg := &ServerConfig{
		GemmaEndpoint:    os.Getenv("GEMMA_ENDPOINT"),
		GemmaModel:       os.Getenv("GEMMA_MODEL"),
		LMStudioEndpoint: os.Getenv("LMSTUDIO_ENDPOINT"),
		LMStudioModel:    os.Getenv("LMSTUDIO_MODEL"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
	}
	if cfg.GemmaEndpoint == "" {
		cfg.GemmaEndpoint = "http://192.168.12.169:11434"
	}
	if cfg.GemmaModel == "" {
		cfg.GemmaModel = "gemma2:2b"
	}
	if cfg.LMStudioEndpoint == "" {
		cfg.LMStudioEndpoint = "http://host.docker.internal:1234"
	}
	if cfg.LMStudioModel == "" {
		cfg.LMStudioModel = "gemma-2-9b-it"
	}
	return cfg
}

var (
	rateMu  sync.Mutex
	clients = make(map[string]*clientRate)
)

type clientRate struct {
	requests     int
	lastReset    time.Time
	backoffUntil time.Time
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rateMu.Lock()
		client, exists := clients[ip]
		if !exists {
			client = &clientRate{lastReset: now}
			clients[ip] = client
		}

		if now.Before(client.backoffUntil) {
			rateMu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Quota exceeded. 300s backoff active."})
			return
		}

		if now.Sub(client.lastReset) >= time.Second {
			client.requests = 0
			client.lastReset = now
		}

		client.requests++
		if client.requests > 30 {
			client.backoffUntil = now.Add(300 * time.Second)
			rateMu.Unlock()
			log.Printf("[SECURITY] Rate limit exceeded for IP %s. Triggering 300s backoff.", ip)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded (30 req/sec). 300s backoff triggered."})
			return
		}
		rateMu.Unlock()
		c.Next()
	}
}

type Server struct {
	Service *service.SwarmService
	Healing *service.HealingService
	Auth    *service.AuthService
	AI      *service.AIService
	Router  *gin.Engine
	Config  *ServerConfig
}

func NewServer(svc *service.SwarmService, healing *service.HealingService, auth *service.AuthService, ai *service.AIService) *Server {
	r := gin.Default()
	s := &Server{
		Service: svc,
		Healing: healing,
		Auth:    auth,
		AI:      ai,
		Router:  r,
		Config:  loadConfig(),
	}

	// Apply Security Headers Middleware
	r.Use(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
		c.Next()
	})

	// Apply Rate Limiting Middleware (30 req/sec, 300s backoff)
	r.Use(rateLimitMiddleware())

	// Static UI serving
	r.StaticFile("/", "./web/dashboard.html")
	r.StaticFile("/legacy", "./web/index.html")
	r.StaticFile("/wiki", "./web/wiki.html")
	r.StaticFile("/hud", "./web/hud.html")
	r.Static("/static", "./web")

	api := r.Group("/REST/2.0")
	{
		// Ticket CRUD
		api.POST("/ticket", s.handleCreateTicket)
		api.GET("/ticket/:id", s.handleGetTicket)
		api.PUT("/ticket/:id", s.handleUpdateTicket)
		api.GET("/tickets", s.handleSearchTickets)

		// Agent memory operations
		api.POST("/agent/:agentID/memory/:ticketID", s.handleStoreMemory)
		api.GET("/agent/:agentID/memory/:ticketID", s.handleGetMemory)
		api.GET("/agent/:agentID/context", s.handleGetAgentContext)

		// Audit and relationships
		api.GET("/ticket/:id/audit", s.handleGetAuditTrail)
		api.GET("/ticket/:id/links", s.handleGetLinks)
		api.POST("/ticket/:parentID/link/:childID", s.handleLinkTickets)

		// Health
		api.GET("/health", s.handleHealth)
		api.GET("/health/gemma", s.handleGemmaHealth)

		// Chat & Swarm Balancing
		api.POST("/chat/gemma", s.handleGemmaChat)
		api.POST("/chat/lmstudio", s.handleLMStudioChat)
		api.POST("/chat/swarm", s.handleSwarmChat)
		api.GET("/health/lmstudio", s.handleLMStudioHealth)

		// Self-healing
		api.POST("/healing/ticket", s.handleCreateHealingTicket)
		api.POST("/healing/iterate/:id", s.handleProcessHealingIteration)
		api.POST("/healing/failure", s.handleRecordHealingFailure)
		api.POST("/healing/resolve", s.handleResolveHealingTicket)

		// Metrics
		api.GET("/metrics/tokens", s.handleGetMetrics)

		// Initialize schema
		api.POST("/init", s.handleInitSchema)

		// Documentation
		api.GET("/docs/:name", s.handleGetDoc)

		// Gemini Emergency Bridge
		api.POST("/emergency/bridge", s.handleEmergencyBridge)

		// Legacy Sovereign API (S25 Compatibility)
		api.GET("/status", s.handleStatus)
		api.GET("/bridge", s.handleBridge)
		api.GET("/files", s.handleListFiles)
		api.GET("/wiki", s.handleWiki)
	}

	// SAML Endpoints
	if s.Auth != nil {
		r.GET("/saml/metadata", gin.WrapH(http.HandlerFunc(s.Auth.HandleMetadata)))
		r.POST("/saml/sso", gin.WrapH(http.HandlerFunc(s.Auth.HandleSSO)))
		r.GET("/saml/sso", gin.WrapH(http.HandlerFunc(s.Auth.HandleSSO)))
	}

	return s
}

func (s *Server) handleCreateTicket(c *gin.Context) {
	var req struct {
		Subject string                 `json:"Subject"`
		Queue   string                 `json:"Queue"`
		Content string                 `json:"Text"`
		AgentID string                 `json:"AgentID"`
		Layer   int                    `json:"Layer"`
		Intent  map[string]interface{} `json:"Intent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AgentID == "" {
		req.AgentID = "REST-API-USER"
	}
	if req.Layer == 0 {
		req.Layer = 2
	}

	fabricContent := domain.FabricContent{
		IntentBlob: req.Intent,
		RawContent: []byte(req.Content),
	}
	if fabricContent.IntentBlob == nil {
		fabricContent.IntentBlob = map[string]interface{}{"subject": req.Subject, "queue": req.Queue}
	}

	ticketID, err := s.Service.CreateFabricTicket(c.Request.Context(), req.Layer, req.AgentID, fabricContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Link to Genesis if it's a new chain
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	s.Service.LinkTicketsWithAudit(c.Request.Context(), genesisID, ticketID, domain.RelEvolution, req.AgentID)

	c.JSON(http.StatusCreated, gin.H{
		"id":      ticketID.String(),
		"message": fmt.Sprintf("Ticket %s created", ticketID),
	})
}

func (s *Server) handleGetTicket(c *gin.Context) {
	idStr := c.Param("id")
	ticketID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	ticket, content, err := s.Service.GetTicketWithContent(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         ticket.ID.String(),
		"layer":      ticket.LayerID,
		"creator":    ticket.CreatorAgentID,
		"status":     ticket.Status,
		"created_at": ticket.CreatedAt,
		"intent":     content.IntentBlob,
		"content":    string(content.RawContent),
	})
}

func (s *Server) handleUpdateTicket(c *gin.Context) {
	idStr := c.Param("id")
	ticketID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	var req struct {
		Status  string `json:"Status"`
		Title   string `json:"Title"`
		Creator string `json:"Creator"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Permission Check: Only creator or assigned agent can update
	ticket, _, err := s.Service.GetTicketWithContent(c.Request.Context(), ticketID)
	if err == nil {
		if ticket.CreatorAgentID != req.Creator && ticket.AssignedTo != req.Creator {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied: Ticket not assigned to agent"})
			return
		}
	}

	err = s.Service.UpdateExtended(c.Request.Context(), ticketID, req.Status, req.Title, "", req.Creator)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (s *Server) handleSearchTickets(c *gin.Context) {
	tickets, err := s.Service.GetRecentTickets(c.Request.Context(), 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to consistent format with content
	var response []gin.H
	for _, t := range tickets {
		// Fetch full content for each ticket to show subject in UI
		_, content, err := s.Service.GetTicketWithContent(c.Request.Context(), t.ID)

		item := gin.H{
			"id":          t.ID.String(),
			"layer":       t.LayerID,
			"creator":     t.CreatorAgentID,
			"status":      t.Status,
			"created_at":  t.CreatedAt,
			"assigned_to": t.AssignedTo,
		}

		if err == nil && content != nil {
			item["intent"] = content.IntentBlob
			item["content"] = string(content.RawContent)
		}

		response = append(response, item)
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

func (s *Server) handleStoreMemory(c *gin.Context) {
	agentID := c.Param("agentID")
	ticketID := c.Param("ticketID")

	var req struct {
		MemType        string                 `json:"memory_type"`
		Data           map[string]interface{} `json:"data"`
		RelevanceScore float64                `json:"relevance_score"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(ticketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	if err := s.Service.StoreAgentMemory(c.Request.Context(), agentID, id, req.MemType, req.Data, req.RelevanceScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "memory stored", "agent": agentID, "ticket": ticketID})
}

func (s *Server) handleGetMemory(c *gin.Context) {
	agentID := c.Param("agentID")
	ticketID := c.Param("ticketID")
	memType := c.Query("type")

	if memType == "" {
		memType = "context"
	}

	id, err := uuid.Parse(ticketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	data, err := s.Service.GetAgentMemory(c.Request.Context(), agentID, id, memType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "memory not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) handleGetAgentContext(c *gin.Context) {
	agentID := c.Param("agentID")

	tickets, err := s.Service.GetAgentContext(c.Request.Context(), agentID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []gin.H
	for _, t := range tickets {
		_, content, err := s.Service.GetTicketWithContent(c.Request.Context(), t.ID)
		item := gin.H{
			"id":         t.ID.String(),
			"layer":      t.LayerID,
			"creator":    t.CreatorAgentID,
			"status":     t.Status,
			"created_at": t.CreatedAt,
		}
		if err == nil && content != nil {
			item["intent"] = content.IntentBlob
			item["content"] = string(content.RawContent)
		}
		response = append(response, item)
	}

	c.JSON(http.StatusOK, gin.H{"agent": agentID, "context_tickets": response})
}

func (s *Server) handleGetAuditTrail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	trail, err := s.Service.GetAuditTrail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": idStr, "audit_trail": trail})
}

func (s *Server) handleLinkTickets(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")

	var req struct {
		RelationType string `json:"relationship_type"`
		AgentID      string `json:"agent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pID, err := uuid.Parse(parentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent id"})
		return
	}

	cID, err := uuid.Parse(childID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child id"})
		return
	}

	relType := domain.RelationshipType(req.RelationType)
	if relType != domain.RelEvolution && relType != domain.RelConsequence && relType != domain.RelContext && relType != domain.RelGenesis {
		relType = domain.RelEvolution
	}

	if err := s.Service.LinkTicketsWithAudit(c.Request.Context(), pID, cID, relType, req.AgentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tickets linked", "parent": parentID, "child": childID})
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "PQR-ticketing",
		"status":  "healthy",
		"version": Version,
	})
}

func (s *Server) handleGetMetrics(c *gin.Context) {
	used, quota, err := s.Service.GetMetric(c.Request.Context(), "tokens_used")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	percent := (used / quota) * 100.0
	c.JSON(http.StatusOK, gin.H{
		"tokens_used":      used,
		"token_quota":      quota,
		"usage_percentage": percent,
	})
}

func (s *Server) handleInitSchema(c *gin.Context) {
	if err := s.Service.InitSchema(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schema initialized"})
}

func (s *Server) handleCreateHealingTicket(c *gin.Context) {
	var req struct {
		Issue      string `json:"issue"`
		LogSnippet string `json:"logSnippet"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := s.Healing.CreateHealingTicket(c.Request.Context(), req.Issue, req.LogSnippet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

func (s *Server) handleProcessHealingIteration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	if err := s.Healing.ProcessHealingLoop(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "iteration processed"})
}

func (s *Server) handleRecordHealingFailure(c *gin.Context) {
	var req struct {
		TicketID string `json:"ticketID"`
		Failure  string `json:"failure"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(req.TicketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID format"})
		return
	}
	if err := s.Healing.RecordFailure(c.Request.Context(), id, req.Failure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "failure recorded"})
}

func (s *Server) handleResolveHealingTicket(c *gin.Context) {
	var req struct {
		TicketID   string `json:"ticketID"`
		Resolution string `json:"resolution"`
		AgentID    string `json:"agentID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := uuid.Parse(req.TicketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID format"})
		return
	}
	if err := s.Healing.MarkResolved(c.Request.Context(), id, req.Resolution, req.AgentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket resolved and added to knowledge base"})
}
func (s *Server) handleGetDoc(c *gin.Context) {
	name := c.Param("name")

	// Securely construct path and prevent traversal
	cleanName := filepath.Clean(name)
	if cleanName == "." || cleanName == ".." || strings.Contains(cleanName, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doc name"})
		return
	}

	path := filepath.Join("docs", cleanName+".md")
	// Additional check to ensure we are still inside docs directory
	absDocs, err := filepath.Abs("docs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve docs directory"})
		return
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve file path"})
		return
	}
	if !strings.HasPrefix(absPath, absDocs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doc path"})
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doc not found"})
		return
	}

	c.String(http.StatusOK, string(content))
}

func (s *Server) handleGemmaHealth(c *gin.Context) {
	gemmaURL := s.Config.GemmaEndpoint

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// Ollama responds to /api/tags or just /
	resp, err := client.Get(gemmaURL + "/api/tags")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "OFFLINE", "error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "ONLINE", "endpoint": gemmaURL})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "code": resp.StatusCode})
	}
}
func (s *Server) handleGetLinks(c *gin.Context) {
	idStr := c.Param("id")
	ticketID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	links, err := s.Service.GetTicketLinks(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket_id": ticketID.String(),
		"links":     links,
	})
}

func (s *Server) handleGemmaChat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
		Model   string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gemmaURL := s.Config.GemmaEndpoint

	modelName := req.Model
	if modelName == "" {
		modelName = s.AI.GetBestOllamaModel()
	}
	modelAgentID := "model-" + modelName

	// 1. Retrieval Augmented Context (RAG)
	contextTickets, _ := s.Service.GetRecentTickets(c.Request.Context(), 3)
	contextText := "Sovereign Mesh Context:\n"
	for _, t := range contextTickets {
		contextText += fmt.Sprintf("- Ticket %s: status is %s\n", t.ID, t.Status)
	}

	// Add model memory
	memories, _ := s.Service.GetAgentContext(c.Request.Context(), modelAgentID, 3)
	contextText += "\nYour Operational Memory:\n"
	for _, t := range memories {
		_, content, err := s.Service.GetTicketWithContent(c.Request.Context(), t.ID)
		if err == nil && content != nil {
			contextText += fmt.Sprintf("- %s\n", string(content.RawContent))
		}
	}

	prompt := fmt.Sprintf("%s\nUser: %s\nAssistant:", contextText, req.Message)

	log.Printf("[GEMMA] Requesting model %s with prompt length %d", modelName, len(prompt))

	performRequest := func(m string) (map[string]interface{}, error) {
		ollamaReq := map[string]interface{}{
			"model": m,
			"messages": []map[string]interface{}{
				{"role": "user", "content": contextText + "\n\nUser Question: " + req.Message},
			},
			"stream": false,
		}

		body, err := json.Marshal(ollamaReq)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		reqObj, err := http.NewRequest("POST", gemmaURL+"/api/chat", bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		reqObj.Header.Set("Content-Type", "application/json")
		reqObj.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(reqObj)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		// log.Printf("[GEMMA] Raw Response: %s", string(respBytes)) // Masked for security

		var result map[string]interface{}
		if err := json.Unmarshal(respBytes, &result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return result, nil
	}

	result, err := performRequest(modelName)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemma offline", "details": err.Error()})
		return
	}

	if errMsg, ok := result["error"].(string); ok {
		log.Printf("[GEMMA] Error from node: %s", errMsg)

		// Create a ticket for the failure (Layer 4)
		ticketContent := domain.FabricContent{
			IntentBlob: map[string]interface{}{
				"type":  "CHAT_FAILURE",
				"query": req.Message,
				"error": errMsg,
				"model": modelName,
			},
			RawContent: []byte("ERROR: " + errMsg),
		}
		s.Service.CreateFabricTicket(c.Request.Context(), 4, "gemma-ai", ticketContent)

		c.JSON(http.StatusOK, gin.H{"response": "ERROR: " + errMsg, "context": contextText})
		return
	}

	// 3. Extract Chat Response
	var respText string
	if msg, ok := result["message"].(map[string]interface{}); ok {
		if content, ok := msg["content"].(string); ok {
			respText = content
		}
	}

	if respText == "" {
		log.Printf("[GEMMA] Empty response from node. Raw: %+v", result)
		respText = "No response from model."
	}

	log.Printf("[GEMMA] Response received (%d bytes). Creating ticket...", len(respText))

	ticketContent := domain.FabricContent{
		IntentBlob: map[string]interface{}{
			"type":  "CHAT_VOLLEY",
			"query": req.Message,
			"model": modelName,
		},
		RawContent: []byte(respText),
	}
	ticketID, _ := s.Service.CreateFabricTicket(c.Request.Context(), 4, modelAgentID, ticketContent)
	_ = s.Service.StoreAgentMemory(c.Request.Context(), modelAgentID, ticketID, "conversation", map[string]interface{}{
		"query":    req.Message,
		"response": respText,
	}, 0.9)

	// Estimate tokens (chars / 4 as a heuristic)
	tokenEstimate := float64(len(req.Message)+len(respText)) / 4.0
	_ = s.Service.IncrementMetric(c.Request.Context(), "tokens_used", tokenEstimate)

	c.JSON(http.StatusOK, gin.H{
		"response": respText,
		"tokens":   tokenEstimate,
		"context":  contextText,
	})
}

func (s *Server) handleLMStudioChat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lmURL := s.Config.LMStudioEndpoint
	lmModel := s.AI.GetBestLMModel()
	modelAgentID := "model-" + lmModel

	// Fetch Agent Memory for the Model
	memories, _ := s.Service.GetAgentContext(c.Request.Context(), modelAgentID, 5)

	var sysMsg strings.Builder
	sysMsg.WriteString("You are the PQR Sovereign Node AI. Here is your operational memory:\n")
	for _, t := range memories {
		_, content, err := s.Service.GetTicketWithContent(c.Request.Context(), t.ID)
		if err == nil && content != nil {
			sysMsg.WriteString(fmt.Sprintf("- %s\n", string(content.RawContent)))
		}
	}

	ollamaReq := map[string]interface{}{
		"model": lmModel,
		"messages": []map[string]interface{}{
			{"role": "system", "content": sysMsg.String()},
			{"role": "user", "content": req.Message},
		},
		"stream": false,
	}
	body, err := json.Marshal(ollamaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode request"})
		return
	}

	reqObj, err := http.NewRequest("POST", lmURL+"/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	reqObj.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(reqObj)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "LM Studio offline", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response from LM Studio"})
		return
	}

	var respText string
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				respText, _ = msg["content"].(string)
			}
		}
	}

	// Store the interaction in the Model's Agent Memory
	ticketContent := domain.FabricContent{
		IntentBlob: map[string]interface{}{
			"type":  "MODEL_MEMORY",
			"query": req.Message,
		},
		RawContent: []byte(fmt.Sprintf("User: %s\nResponse: %s", req.Message, respText)),
	}
	ticketID, _ := s.Service.CreateFabricTicket(c.Request.Context(), 3, modelAgentID, ticketContent)
	_ = s.Service.StoreAgentMemory(c.Request.Context(), modelAgentID, ticketID, "conversation", map[string]interface{}{
		"query":    req.Message,
		"response": respText,
	}, 0.9)

	// Estimate tokens
	tokenEstimate := float64(len(req.Message)+len(respText)) / 4.0
	_ = s.Service.IncrementMetric(c.Request.Context(), "tokens_used", tokenEstimate)

	c.JSON(http.StatusOK, gin.H{
		"response": respText,
		"tokens":   tokenEstimate,
	})
}

func (s *Server) handleLMStudioHealth(c *gin.Context) {
	client := http.Client{Timeout: 1 * time.Second}

	lmURL := s.Config.LMStudioEndpoint

	resp, err := client.Get(lmURL + "/v1/models")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "OFFLINE"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "ERROR", "code": resp.StatusCode})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ONLINE"})
}
func (s *Server) handleSwarmChat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Give the Swarm Coordinator memory
	agentID := "swarm-coordinator"
	memories, _ := s.Service.GetAgentContext(c.Request.Context(), agentID, 5)

	var enhancedPrompt strings.Builder
	enhancedPrompt.WriteString("System Memory / Knowledge Base:\n")
	for _, t := range memories {
		_, content, err := s.Service.GetTicketWithContent(c.Request.Context(), t.ID)
		if err == nil && content != nil {
			enhancedPrompt.WriteString(fmt.Sprintf("- %s\n", string(content.RawContent)))
		}
	}
	enhancedPrompt.WriteString("\nUser Message: " + req.Message)

	resp, engine, err := s.AI.QuerySwarm(c.Request.Context(), enhancedPrompt.String())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Swarm AI nodes are offline"})
		return
	}

	// Pipe the Swarm's solution back into the Ticketing Fabric as a Knowledge Base entry
	knowledgeIntent := map[string]interface{}{
		"type":           "KNOWLEDGE_BASE_ENTRY",
		"original_query": req.Message,
		"solving_engine": engine,
	}
	knowledgeContent := domain.FabricContent{
		IntentBlob: knowledgeIntent,
		RawContent: []byte(resp),
	}

	// Create a Layer 3 (Knowledge) ticket
	kbTicketID, createErr := s.Service.CreateFabricTicket(c.Request.Context(), 3, "swarm-coordinator", knowledgeContent)
	if createErr == nil {
		// Store it in the local-triage-agent's explicit memory so local models can retrieve it via RAG
		_ = s.Service.StoreAgentMemory(c.Request.Context(), "local-triage-agent", kbTicketID, "knowledge", map[string]interface{}{
			"query":    req.Message,
			"solution": resp,
		}, 0.95)
		log.Printf("[SWARM] Piped solution back to knowledge base (Ticket: %s)", kbTicketID.String())
	} else {
		log.Printf("[SWARM] Warning: Failed to store knowledge base entry: %v", createErr)
	}

	// Estimate tokens
	tokenEstimate := float64(len(req.Message)+len(resp)) / 4.0
	_ = s.Service.IncrementMetric(c.Request.Context(), "tokens_used", tokenEstimate)

	c.JSON(http.StatusOK, gin.H{
		"response": resp,
		"engine":   engine,
		"tokens":   tokenEstimate,
	})
}

func (s *Server) handleEmergencyBridge(c *gin.Context) {
	// Verify Gemini API Key for Emergency Access
	apiKey := c.GetHeader("X-Gemini-Key")
	expectedKey := s.Config.GeminiAPIKey

	if apiKey == "" || apiKey != expectedKey {
		log.Printf("[EMERGENCY] ⚠️ Unauthorized bridge attempt from %s", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Denied: Invalid Emergency Key"})
		return
	}

	var req struct {
		Command string                 `json:"command"`
		Params  map[string]interface{} `json:"params"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[EMERGENCY] ⚡ Gemini Command Received: %s", req.Command)

	switch req.Command {
	case "GET_SYSTEM_HEALTH":
		status := "HEALTHY"
		if s.Auth == nil {
			status = "AUTH_DEGRADED"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"node":    "pqr-sovereign-001",
			"uptime":  time.Now().Format(time.RFC3339),
			"version": Version,
		})

	case "LIST_RECENT_TICKETS":
		tickets, _ := s.Service.GetRecentTickets(c.Request.Context(), 10)
		c.JSON(http.StatusOK, tickets)

	case "TRIGGER_HEALING":
		issue, _ := req.Params["issue"].(string)
		logSnippet, _ := req.Params["logSnippet"].(string)
		id, err := s.Healing.CreateHealingTicket(c.Request.Context(), issue, logSnippet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"healing_ticket_id": id.String()})

	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown emergency command"})
	}
}

func (s *Server) handleStatus(c *gin.Context) {
	used, quota, _ := s.Service.GetMetric(c.Request.Context(), "tokens_used")
	vitality := 100.0
	if quota > 0 {
		vitality = 100.0 - ((used / quota) * 100.0)
	}

	c.JSON(http.StatusOK, gin.H{
		"node_id":  "ΩX9R2#",
		"status":   "SINGULARITY",
		"vitality": vitality,
		"up_time":  time.Now().Format(time.RFC3339),
		"logic":    "AELLOK-V10",
	})
}

func (s *Server) handleBridge(c *gin.Context) {
	cmd := c.Query("cmd")
	if cmd == "" {
		c.String(http.StatusBadRequest, "No command provided")
		return
	}
	output, err := s.Healing.ExecuteDiagnostic(c.Request.Context(), cmd)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, output)
}

func (s *Server) handleListFiles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"files": []string{"server.go", "Dockerfile", "docs/ARCHITECTURE.md"},
	})
}

func (s *Server) handleWiki(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"sections": []string{"Overview", "Identity", "Fabric", "Swarm"},
	})
}
