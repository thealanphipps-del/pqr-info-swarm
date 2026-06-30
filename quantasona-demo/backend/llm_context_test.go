package backend

import (
	"testing"
	"sovereign-node-go/pkg/rtgo"
)

func TestFormatTicketForLLM(t *testing.T) {
	ticket := rtgo.TicketSummary{
		ID:         "123e4567-e89b-12d3-a456-426614174000",
		Title:      "Integrate LLM",
		Status:     rtgo.StatusPending,
		RawContent: "Need to integrate LLM to read existing code and tickets.",
	}

	result := FormatTicketForLLM(ticket)
	expected := "Ticket ID: 123e4567-e89b-12d3-a456-426614174000\nTitle: Integrate LLM\nStatus: PENDING\nContent: Need to integrate LLM to read existing code and tickets."

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
