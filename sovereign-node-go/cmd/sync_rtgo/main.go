package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("[SYNC] Initializing RTGO -> CockroachDB Sync...")

	// Connect to local CockroachDB
	connStr := "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to local CockroachDB: %v", err)
	}
	defer db.Close()

	fmt.Println("[SYNC] Fetching tickets from 39.MH via secure mesh SSH...")

	// Execute SSH to fetch data as CSV to bypass port 5432 restrictions
	// Using the known 39.mh host alias or IP.
	cmd := exec.Command("ssh", "root@204.168.138.60", "sudo -u postgres psql -d rtgo_ticketing_system -c \"COPY tickets TO STDOUT WITH CSV HEADER\"")
	
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("[FATAL] SSH command failed. Ensure your local SSH keys are configured for 39.mh: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("[FATAL] Failed to parse CSV data: %v", err)
	}

	if len(records) <= 1 {
		fmt.Println("[SYNC] No tickets found or only header returned.")
		return
	}

	fmt.Printf("[SYNC] Fetched %d records. Upserting to local CockroachDB...\n", len(records)-1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	successCount := 0
	for i, row := range records {
		if i == 0 {
			continue // Skip header
		}

		// Schema: ticket_id, ticket_type, content, path, status, created_at, updated_at
		if len(row) < 7 {
			continue
		}

		_, err := db.ExecContext(ctx, `
			UPSERT INTO rtgo (ticket_id, ticket_type, content, path, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, row[0], row[1], row[2], row[3], row[4], row[5], row[6])

		if err != nil {
			fmt.Printf("[ERROR] Failed to upsert ticket %s: %v\n", row[0], err)
		} else {
			successCount++
		}
	}

	fmt.Printf("[SUCCESS] Successfully synchronized %d tickets into local CockroachDB.\n", successCount)
}
