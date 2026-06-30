package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"amln-sen/internal/pqr"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rtgo-client [recall/create] [ticket_id/title] [optional: baseURL]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	target := os.Args[2]

	baseURL := "http://localhost:8085"
	if len(os.Args) >= 4 {
		baseURL = os.Args[3]
	}

	session := pqr.NewSession(baseURL, "SYSTEM")

	ctx := context.Background()

	switch cmd {
	case "recall":
		memory, err := session.RecallMemory(ctx, target)
		if err != nil {
			fmt.Printf("ERROR: failed to recall memory: %v\n", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(memory, "", "  ")
		fmt.Println(string(out))

	case "create":
		data := map[string]interface{}{
			"description": "Created via Go RT REST 2.0 Client CLI",
			"timestamp":   "dynamic",
		}
		ticketID, err := session.CreateMemory(ctx, target, data)
		if err != nil {
			fmt.Printf("ERROR: failed to create memory ticket: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SUCCESS: ticket %s created\n", ticketID)

	default:
		fmt.Printf("ERROR: unknown command %s\n", cmd)
		os.Exit(1)
	}
}
