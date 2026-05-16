package pqr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/thealanphipps-del/pqr/internal/service"
)

const Version = "v1.07"

type Server struct {
	Service *service.SwarmService
	Healing *service.HealingService
	Auth    *service.AuthService
	AI      *service.AIService
	Router  *gin.Engine
}

func NewServer(svc *service.SwarmService, healing *service.HealingService, auth *service.AuthService, ai *service.AIService) *Server {
	r := gin.Default()
	s := &Server{
		Service: svc,
		Healing: healing,
		Auth:    auth,
		AI:      ai,
		Router:  r,
	}

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
	c.JSON(http.StatusOK, tickets)
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

func (s *Server) handleStoreMemory(c *gin.Context) {
	agentID := c.Param("agentID")
	ticketID := c.Param("ticketID")
	
	var req struct {
		MemType         string                 `json:"memory_type"`
		Data            map[string]interface{} `json:"data"`
		RelevanceScore  float64                `json:"relevance_score"`
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
	
	c.JSON(http.StatusOK, gin.H{"agent": agentID, "context_tickets": tickets})
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
	id, _ := uuid.Parse(req.TicketID)
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
	id, _ := uuid.Parse(req.TicketID)
	if err := s.Healing.MarkResolved(c.Request.Context(), id, req.Resolution, req.AgentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket resolved and added to knowledge base"})
}
func (s *Server) handleGetDoc(c *gin.Context) {
	name := c.Param("name")
	// Sanitize name to prevent path traversal
	if name == "" || name == ".." || name == "." {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doc name"})
		return
	}

	path := fmt.Sprintf("./docs/%s.md", name)
	content, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doc not found"})
		return
	}

	c.String(http.StatusOK, string(content))
}

func (s *Server) handleGemmaHealth(c *gin.Context) {
	gemmaURL := os.Getenv("GEMMA_ENDPOINT")
	if gemmaURL == "" {
		gemmaURL = "http://192.168.12.169:11434"
	}

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

	// For now we'll just return an empty list or fetch from DB if implemented
	// In a real SG-DAO this would query the ticket_relationships table
	c.JSON(http.StatusOK, gin.H{
		"ticket_id": ticketID.String(),
		"links":     []string{},
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

	gemmaURL := os.Getenv("GEMMA_ENDPOINT")
	if gemmaURL == "" {
		gemmaURL = "http://192.168.12.169:11434"
	}

	modelName := req.Model
	if modelName == "" {
		modelName = os.Getenv("GEMMA_MODEL")
		if modelName == "" {
			modelName = "gemma2:2b"
		}
	}

	// 1. Retrieval Augmented Context (RAG)
	contextTickets, _ := s.Service.GetRecentTickets(c.Request.Context(), 3)
	contextText := "Sovereign Mesh Context:\n"
	for _, t := range contextTickets {
		contextText += fmt.Sprintf("- Ticket %s: status is %s\n", t.ID, t.Status)
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

		body, _ := json.Marshal(ollamaReq)
		
		reqObj, _ := http.NewRequest("POST", gemmaURL+"/api/chat", bytes.NewBuffer(body))
		reqObj.Header.Set("Content-Type", "application/json")
		reqObj.Header.Set("Accept", "application/json")
		
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(reqObj)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		
		respBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[GEMMA] Raw Response: %s", string(respBytes))
		
		var result map[string]interface{}
		json.Unmarshal(respBytes, &result)
		return result, nil
	}


	result, err := performRequest(modelName)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemma offline", "details": err.Error()})
		return
	}

	// Fallback logic
	if errMsg, ok := result["error"].(string); ok && strings.Contains(errMsg, "not found") && modelName == "gemma2:2b" {
		log.Printf("[GEMMA] Model %s not found, falling back to gemma2", modelName)
		modelName = "gemma2"
		result, err = performRequest(modelName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Gemma offline during fallback"})
			return
		}
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
	s.Service.CreateFabricTicket(c.Request.Context(), 4, "gemma-ai", ticketContent)

	// Estimate tokens (chars / 4 as a heuristic)
	tokenEstimate := float64(len(req.Message) + len(respText)) / 4.0
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

	lmURL := "http://host.docker.internal:1234"

	ollamaReq := map[string]interface{}{
		"model": "gemma-2-9b-it",
		"messages": []map[string]interface{}{
			{"role": "user", "content": req.Message},
		},
		"stream": false,
	}
	body, _ := json.Marshal(ollamaReq)
	
	reqObj, _ := http.NewRequest("POST", lmURL+"/v1/chat/completions", bytes.NewBuffer(body))
	reqObj.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(reqObj)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "LM Studio offline", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	var respText string
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				respText, _ = msg["content"].(string)
			}
		}
	}

	// Estimate tokens
	tokenEstimate := float64(len(req.Message) + len(respText)) / 4.0
	_ = s.Service.IncrementMetric(c.Request.Context(), "tokens_used", tokenEstimate)

	c.JSON(http.StatusOK, gin.H{
		"response": respText,
		"tokens":   tokenEstimate,
	})
}

func (s *Server) handleLMStudioHealth(c *gin.Context) {
	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get("http://host.docker.internal:1234/v1/models")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "OFFLINE"})
		return
	}
	defer resp.Body.Close()
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

	resp, engine, err := s.AI.QuerySwarm(c.Request.Context(), req.Message)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Swarm AI nodes are offline"})
		return
	}

	// Estimate tokens
	tokenEstimate := float64(len(req.Message) + len(resp)) / 4.0
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
	expectedKey := os.Getenv("GEMINI_API_KEY")

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
		if s.Auth == nil { status = "AUTH_DEGRADED" }
		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"node": "pqr-sovereign-001",
			"uptime": time.Now().Format(time.RFC3339),
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
	// Simple telemetry mirroring S25
	c.JSON(http.StatusOK, gin.H{
		"node_id":   "ΩX9R2#",
		"status":    "SINGULARITY",
		"vitality":  98.4,
		"up_time":   "12:44:12",
		"logic":     "AELLOK-V10",
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
