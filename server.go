package pqr

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thealanphipps-del/pqr/internal/domain"
	"github.com/thealanphipps-del/pqr/internal/service"
)

type Server struct {
	Service *service.SwarmService
	Healing *service.HealingService
	Router  *gin.Engine
}

func NewServer(svc *service.SwarmService, healing *service.HealingService) *Server {
	r := gin.Default()
	s := &Server{
		Service: svc,
		Healing: healing,
		Router:  r,
	}

	// Static UI serving
	r.StaticFile("/", "./web/index.html")
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
		api.POST("/ticket/:parentID/link/:childID", s.handleLinkTickets)
		
		// Database health check
		api.GET("/health", s.handleHealth)
		
		// Self-healing
		api.POST("/healing/ticket", s.handleCreateHealingTicket)
		api.POST("/healing/iterate/:id", s.handleProcessHealingIteration)
		api.POST("/healing/failure", s.handleRecordHealingFailure)
		api.POST("/healing/resolve", s.handleResolveHealingTicket)
		
		// Initialize schema
		api.POST("/init", s.handleInitSchema)

		// Documentation
		api.GET("/docs/:name", s.handleGetDoc)
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
		Status string `json:"Status"`
		Title  string `json:"Title"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Service.UpdateTicket(c.Request.Context(), ticketID, req.Status, req.Title)
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
	// Simple health check
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "PQR-ticketing"})
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
