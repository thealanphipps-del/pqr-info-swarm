package pqr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCockpitRESTAndClient(t *testing.T) {
	// 1. Setup mock Server that matches our new endpoints
	r := gin.Default()
	r.GET("/REST/2.0/tickets", func(c *gin.Context) {
		c.JSON(http.StatusOK, []map[string]interface{}{
			{"id": "ticket-100", "status": "OPEN", "assigned_to": "agent-001", "created_at": "2026-06-25T20:00:00Z"},
		})
	})
	r.GET("/REST/2.0/ticket/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"id":          c.Param("id"),
			"layer":       2,
			"creator":     "agent-001",
			"assigned_to": "agent-001",
			"status":      "OPEN",
			"created_at":  "2026-06-25T20:00:00Z",
			"content":     "Hello world",
		})
	})
	r.PUT("/REST/2.0/ticket/:id", func(c *gin.Context) {
		var req struct {
			Status     string `json:"Status"`
			Title      string `json:"Title"`
			Creator    string `json:"Creator"`
			AssignedTo string `json:"AssignedTo"`
			Priority   string `json:"Priority"`
			Queue      string `json:"Queue"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated", "assigned_to": req.AssignedTo, "priority": req.Priority})
	})
	r.POST("/REST/2.0/ticket/:id/comment", func(c *gin.Context) {
		var req struct {
			AgentID string `json:"AgentID"`
			Comment string `json:"Comment"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "comment added", "comment": req.Comment})
	})

	srv := httptest.NewServer(r)
	defer srv.Close()

	client := NewClient(srv.URL)
	ctx := context.Background()

	// 2. Test ListTickets
	tickets, err := client.ListTickets(ctx)
	if err != nil {
		t.Fatalf("failed to list tickets: %v", err)
	}
	if len(tickets) != 1 || tickets[0]["id"] != "ticket-100" {
		t.Errorf("unexpected tickets list: %v", tickets)
	}

	// 3. Test UpdateTicketExtended
	err = client.UpdateTicketExtended(ctx, "ticket-100", "OPEN", "New Title", "operator", "agent-002", "high", "Swarm::Operator")
	if err != nil {
		t.Fatalf("failed to update ticket extended: %v", err)
	}

	// 4. Test CommentTicket
	err = client.CommentTicket(ctx, "ticket-100", "operator", "test comment")
	if err != nil {
		t.Fatalf("failed to comment ticket: %v", err)
	}
}
