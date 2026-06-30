package pqr

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleAgentUsage demonstrates how agents use the ticketing system for memory
func ExampleAgentUsage() {
	// Initialize client
	client := NewClient("http://localhost:8080")
	agentID := "processing-agent-001"
	ctx := context.Background()

	// 1. Check service health
	healthy, err := client.Health(ctx)
	if !healthy {
		log.Fatal("Service not healthy:", err)
	}
	fmt.Println("✓ Service healthy")

	// 2. Create a task ticket (working memory)
	ticketID, err := client.CreateTicket(ctx,
		"Process Data Batch",
		"data-processing",
		"Initial batch of 1000 records to process",
		agentID,
		map[string]interface{}{
			"batch_id": "batch-2024-001",
			"items":    1000,
			"source":   "api-feed",
		},
	)
	if err != nil {
		log.Fatal("Failed to create ticket:", err)
	}
	fmt.Printf("✓ Created ticket: %s\n", ticketID)

	// 3. Store working state as agent memory
	err = client.StoreMemory(ctx, agentID, ticketID, "context",
		map[string]interface{}{
			"status":           "processing",
			"items_processed":  0,
			"items_total":      1000,
			"current_position": 0,
			"errors":           0,
		},
		1.0, // Full relevance for current task
	)
	if err != nil {
		log.Fatal("Failed to store memory:", err)
	}
	fmt.Println("✓ Stored initial memory")

	// 4. Simulate processing and update memory
	for i := 1; i <= 5; i++ {
		// Do some work...
		time.Sleep(100 * time.Millisecond)

		// Update memory with progress
		err = client.StoreMemory(ctx, agentID, ticketID, "context",
			map[string]interface{}{
				"status":           "processing",
				"items_processed":  i * 200,
				"items_total":      1000,
				"current_position": i * 200,
				"errors":           0,
			},
			1.0,
		)
		if err != nil {
			log.Fatal("Failed to update memory:", err)
		}
		fmt.Printf("✓ Progress: %d/1000 items\n", i*200)
	}

	// 5. Store learned knowledge for future use
	err = client.StoreMemory(ctx, agentID, ticketID, "knowledge",
		map[string]interface{}{
			"patterns": []string{
				"null_value_handling",
				"date_format_parsing",
				"unicode_normalization",
			},
			"success_rate": 0.98,
			"learned_at":   time.Now(),
		},
		0.85, // Lower relevance for knowledge vs current task
	)
	if err != nil {
		log.Fatal("Failed to store knowledge:", err)
	}
	fmt.Println("✓ Stored learned patterns")

	// 6. Complete the task
	err = client.UpdateTicket(ctx, ticketID, "COMPLETED", "Batch processing complete")
	if err != nil {
		log.Fatal("Failed to update ticket:", err)
	}
	fmt.Println("✓ Marked task complete")

	// 7. Retrieve final memory state for review
	finalState, err := client.GetMemory(ctx, agentID, ticketID, "context")
	if err != nil {
		log.Fatal("Failed to retrieve memory:", err)
	}
	fmt.Printf("✓ Final state: %v\n", finalState)

	// 8. Review audit trail
	audit, err := client.GetAuditTrail(ctx, ticketID)
	if err != nil {
		log.Fatal("Failed to get audit:", err)
	}
	fmt.Printf("✓ Audit trail: %d entries\n", len(audit))

	// 9. Get all agent context (useful for next task)
	allMemories, err := client.GetContext(ctx, agentID)
	if err != nil {
		log.Fatal("Failed to get context:", err)
	}
	fmt.Printf("✓ Agent context: %d tickets available\n", len(allMemories))

	// 10. Demonstrate state resumption
	// Create an incomplete ticket to simulate an interrupted task
	crashTicketID, _ := client.CreateTicket(ctx, "Crash Test", "testing", "Task that will be interrupted", agentID, nil)
	client.StoreMemory(ctx, agentID, crashTicketID, "context", map[string]interface{}{
		"step": 2,
		"last_processed_id": 42,
	}, 1.0)
	
	// Simulated crash and restart with a new session
	newSession := NewAgentSession("http://localhost:8080", agentID)
	resumedTicket, state, err := newSession.Resume(ctx)
	if err == nil {
		fmt.Printf("✓ Resumed agent state! Ticket: %s\n", resumedTicket)
		if step, ok := state["step"].(float64); ok {
			fmt.Printf("  Resuming from step: %v\n", step)
		}
		client.UpdateTicket(ctx, resumedTicket, "COMPLETED", "Resumed and finished")
	}

	fmt.Println("\n✓ Example complete - Agent memory system working!")
}

// ExampleMultiAgentCoordination shows agents coordinating through linked tickets
func ExampleMultiAgentCoordination() {
	client := NewClient("http://localhost:8080")
	ctx := context.Background()

	// Agent A creates analysis ticket
	analysisID, err := client.CreateTicket(ctx,
		"Data Analysis",
		"analytics",
		"Analyze processed data",
		"analysis-agent-001",
		map[string]interface{}{"type": "statistical"},
	)
	if err != nil {
		log.Fatal("Failed to create analysis ticket:", err)
	}
	fmt.Printf("Analysis Agent created: %s\n", analysisID)

	// Agent A stores its analysis plan
	err = client.StoreMemory(ctx, "analysis-agent-001", analysisID, "context",
		map[string]interface{}{
			"methods":  []string{"mean", "median", "stddev"},
			"datasets": 1,
			"results":  0,
		},
		1.0,
	)
	if err != nil {
		log.Fatal("Failed to store memory:", err)
	}

	// Agent B creates report ticket
	reportID, err := client.CreateTicket(ctx,
		"Generate Report",
		"reporting",
		"Create final report from analysis",
		"reporting-agent-001",
		map[string]interface{}{"type": "business_report"},
	)
	if err != nil {
		log.Fatal("Failed to create report ticket:", err)
	}
	fmt.Printf("Reporting Agent created: %s\n", reportID)

	// Agent B links its work to Agent A's analysis
	err = client.LinkTickets(ctx, analysisID, reportID, "CONSEQUENCE", "reporting-agent-001")
	if err != nil {
		log.Fatal("Failed to link tickets:", err)
	}
	fmt.Println("✓ Linked analysis → report")

	// Agent A can see Agent B's dependent work
	ticket, err := client.GetTicket(ctx, analysisID)
	if err != nil {
		log.Fatal("Failed to get ticket:", err)
	}
	fmt.Printf("Analysis ticket status: %s\n", ticket["status"])

	// Both agents can track the relationship
	audit, err := client.GetAuditTrail(ctx, analysisID)
	if err != nil {
		log.Fatal("Failed to get audit trail:", err)
	}
	fmt.Printf("Analysis has %d audit entries\n", len(audit))

	fmt.Println("✓ Multi-agent coordination working!")
}
