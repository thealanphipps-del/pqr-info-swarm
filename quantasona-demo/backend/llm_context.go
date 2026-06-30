package backend

import (
	"fmt"
	"sovereign-node-go/pkg/rtgo"
)

// FormatTicketForLLM takes a TicketSummary and returns a string formatted for an LLM prompt.
func FormatTicketForLLM(ticket rtgo.TicketSummary) string {
	return fmt.Sprintf("Ticket ID: %s\nTitle: %s\nStatus: %s\nContent: %s", 
		ticket.ID, ticket.Title, ticket.Status, ticket.RawContent)
}
