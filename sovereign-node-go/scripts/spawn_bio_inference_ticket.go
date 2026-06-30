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
	creator := "SNAPDRAGON-NPU-NODE-01"
	assigned := "SNAPDRAGON-NPU-NODE-01"
	status := "NEW"
	priority := 30 // High
	queue := "AudioDiagnosis"
	title := "BIO_INFERENCE for Vickie Dean"

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
		"action": "BIO_INFERENCE",
		"z": []float64{0.2, -0.1, 0.5, 0.7, -0.3},
		"nutrients": map[string]float64{
			"B12": 0.2,
			"Iron": 0.5,
			"Mg": 0.4,
		},
	}
	intentBlob, _ := json.Marshal(intent)
	rawContent := "Perform bio-acoustic four-stage inference on extracted z band feature vector for Vickie Dean."

	// 2. Insert into ticket_content
	_, err = db.Exec(`
		INSERT INTO ticket_content (ticket_id, intent_blob, raw_content)
		VALUES ($1, $2, $3)
	`, ticketID, intentBlob, []byte(rawContent))
	if err != nil {
		log.Fatal("Failed to insert ticket content: ", err)
	}

	fmt.Printf("Successfully injected bio-inference test ticket %s\n", ticketID)
}
