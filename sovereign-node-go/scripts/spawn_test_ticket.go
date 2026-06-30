package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ticketID := uuid.New()
	layerID := 2
	creator := "GATE-01"
	assigned := "ORACLE"
	status := "NEW"
	priority := 30 // High
	queue := "Database"
	title := "Trigger test Database Sync"

	// 1. Insert into tickets
	_, err = db.Exec(`
		INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status, priority, queue, assigned_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, ticketID, layerID, creator, status, priority, queue, assigned)
	if err != nil {
		log.Fatal("Failed to insert ticket: ", err)
	}

	intent := map[string]interface{}{
		"title": title,
		"action": "SYNC_TEST",
	}
	intentBlob, _ := json.Marshal(intent)
	rawContent := "Perform validation on mesh node syncing. Note: If connection is broken, output 'ERROR: Database Connection Failure. Initiating self-healing sequence.' to request recovery."

	// 2. Insert into ticket_content
	_, err = db.Exec(`
		INSERT INTO ticket_content (ticket_id, intent_blob, raw_content)
		VALUES ($1, $2, $3)
	`, ticketID, intentBlob, []byte(rawContent))
	if err != nil {
		log.Fatal("Failed to insert ticket content: ", err)
	}

	fmt.Printf("Inserted test ticket %s\n", ticketID)
}
