package rtgo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	Manager *Manager
	Router  *gin.Engine
}

func NewServer(mgr *Manager) *Server {
	r := gin.Default()
	s := &Server{
		Manager: mgr,
		Router:  r,
	}

	api := r.Group("/REST/2.0")
	{
		api.POST("/ticket", s.handleCreateTicket)
		api.GET("/ticket/:id", s.handleGetTicket)
		api.PUT("/ticket/:id", s.handleUpdateTicket)
		api.GET("/tickets", s.handleSearchTickets)
		api.GET("/tickets/tree", s.handleTicketTree)
		api.GET("/ticket/:id/forensic", s.handleExportTicket)
		api.GET("/search/semantic", s.handleSemanticSearch)
		api.GET("/health", s.handleHealth)

		// Human-in-the-Loop Replay Determination Endpoints
		api.POST("/hitl/replay/:id", s.handleHITLRequestReplay)
		api.POST("/hitl/resolve/:id", s.handleHITLResolveReplay)
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
		req.Layer = 2 // Default to Layer 2 (child of Genesis)
	}

	ctx := c.Request.Context()
	fabricContent := FabricContent{
		IntentBlob: req.Intent,
		RawContent: []byte(req.Content),
	}
	if fabricContent.IntentBlob == nil {
		fabricContent.IntentBlob = map[string]interface{}{"subject": req.Subject, "queue": req.Queue}
	}

	ticketID, err := s.Manager.CreateFabricTicketV71(ctx, req.Layer, req.AgentID, fabricContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Link to Genesis if it's a new chain
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	s.Manager.LinkTicketsV71(ctx, genesisID, ticketID, RelEvolution)

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

	ctx := c.Request.Context()
	// Using the fabric logic to get context
	// We need to know which agent is requesting to get their "active ticket" context,
	// but here we'll just fetch the specific ticket.

	var t FabricTicket
	var content FabricContent
	var intentJSON []byte

	err = s.Manager.store.db.QueryRowContext(ctx, `
		SELECT t.ticket_id, t.layer_id, t.creator_agent_id, t.status, t.priority, t.queue, t.assigned_to, t.is_sticky, t.referrer_code, t.created_at,
		       c.intent_blob, c.raw_content
		FROM tickets t
		LEFT JOIN ticket_content c ON t.ticket_id = c.ticket_id
		WHERE t.ticket_id = $1
	`, ticketID).Scan(&t.ID, &t.LayerID, &t.CreatorAgentID, &t.Status, &t.Priority, &t.Queue, &t.AssignedTo, &t.IsSticky, &t.ReferrerCode, &t.CreatedAt, &intentJSON, &content.RawContent)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	if intentJSON != nil {
		json.Unmarshal(intentJSON, &content.IntentBlob)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            t.ID.String(),
		"layer":         t.LayerID,
		"creator":       t.CreatorAgentID,
		"status":        t.Status,
		"priority":      t.Priority,
		"queue":         t.Queue,
		"assigned_to":   t.AssignedTo,
		"is_sticky":     t.IsSticky,
		"referrer_code": t.ReferrerCode,
		"created_at":    t.CreatedAt,
		"intent":        content.IntentBlob,
		"content":       string(content.RawContent),
	})
}

func (s *Server) handleUpdateTicket(c *gin.Context) {
	idStr := c.Param("id")
	var req struct {
		Status     string `json:"Status"`
		Title      string `json:"Title"`
		Priority   int    `json:"Priority"`
		AssignedTo string `json:"AssignedTo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.Manager.UpdateTicket(c.Request.Context(), idStr, req.Title, req.Status, req.Priority, req.AssignedTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (s *Server) handleSearchTickets(c *gin.Context) {
	tix, err := s.Manager.FetchTickets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tix == nil {
		tix = []TicketSummary{}
	}
	c.JSON(http.StatusOK, tix)
}

func (s *Server) handleTicketTree(c *gin.Context) {
	tree, err := s.Manager.FetchTicketTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tree == nil {
		tree = []*TreeNode{}
	}
	c.JSON(http.StatusOK, tree)
}

func (s *Server) handleExportTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}
	packet, err := s.Manager.ExportForensicPacket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=forensic_%s.md", idStr[:8]))
	c.Data(http.StatusOK, "text/markdown", []byte(packet))
}

func (s *Server) handleSemanticSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query required"})
		return
	}
	results, err := s.Manager.SemanticSearch(c.Request.Context(), query, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "Sovereign Node Go HUD",
		"version": "v8.5-RTGO",
	})
}

func (s *Server) handleHITLRequestReplay(c *gin.Context) {
	idStr := c.Param("id")
	err := s.Manager.UpdateTicket(c.Request.Context(), idStr, "", StatusHITLRequired, PriorityCritical, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sequence buffered for Human-In-The-Loop Replay Determination", "status": StatusHITLRequired})
}

func (s *Server) handleHITLResolveReplay(c *gin.Context) {
	idStr := c.Param("id")
	var req struct {
		Action string `json:"action"` // "APPROVE" or "REJECT"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStatus := StatusResolved
	if req.Action == "REJECT" {
		newStatus = StatusDeleted
	}

	err := s.Manager.UpdateTicket(c.Request.Context(), idStr, "", newStatus, PriorityNormal, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("HITL Replay Determination: %s", req.Action), "status": newStatus})
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}
