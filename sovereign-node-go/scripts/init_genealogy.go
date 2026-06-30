package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

	file, err := os.Open("rt_modules.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nodeMap := make(map[string]uuid.UUID)
	genesisID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	// 1. Create root 'RT' ticket if it doesn't exist
	// In this implementation, we'll just create all from scratch or ignore if exists
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		ticketID := uuid.New()
		_, err = db.Exec(`
			INSERT INTO tickets (ticket_id, layer_id, creator_agent_id, status, priority, queue)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, ticketID, 1, "SYSTEM-GENEALOGY", "IMMUTABLE", 0, "Genealogy")
		if err != nil {
			log.Printf("Failed to insert ticket %s: %v", line, err)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO ticket_content (ticket_id, intent_blob, raw_content)
			VALUES ($1, $2, $3)
		`, ticketID, fmt.Sprintf(`{"title": "%s", "type": "GENEALOGY_PLACEHOLDER"}`, line), []byte("Genealogical placeholder for "+line))
		if err != nil {
			log.Printf("Failed to insert content %s: %v", line, err)
			continue
		}

		nodeMap[line] = ticketID
	}

	// 2. Establish relationships based on naming hierarchy
	for childName, childID := range nodeMap {
		parentName := ""
		if idx := strings.LastIndex(childName, "::"); idx != -1 {
			parentName = childName[:idx]
		}

		parentID := genesisID
		if parentName != "" {
			if id, ok := nodeMap[parentName]; ok {
				parentID = id
			}
		}

		_, err = db.Exec(`
			INSERT INTO ticket_relationships (parent_id, child_id, relationship_type)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, parentID, childID, "GENESIS")
		if err != nil {
			log.Printf("Failed to link %s -> %s: %v", parentName, childName, err)
		}
	}

	fmt.Println("Genealogy initialization complete.")
}
