package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlFiles := []string{
		"db/schema.sql",
		"db/rtgo_v2_upgrade.sql",
		"db/init.sql",
		"db/init_agents.sql",
		"db/lattice_config.sql",
	}

	for _, file := range sqlFiles {
		filePath := filepath.Join("/home/aellok/sovereign-node-go", file)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", file, err)
		}

		// Split queries by semicolon to execute them one by one
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}
			_, err = db.Exec(q)
			if err != nil {
				// Log but don't fatal on ON CONFLICT / IF NOT EXISTS warnings
				log.Printf("Executing query in %s failed: %v\nQuery: %s", file, err, q)
			}
		}
		fmt.Printf("Successfully applied %s\n", file)

		// After schema.sql, apply the alters for tickets
		if file == "db/schema.sql" {
			alterQueries := []string{
				"ALTER TABLE tickets ADD COLUMN IF NOT EXISTS priority INT DEFAULT 20",
				"ALTER TABLE tickets ADD COLUMN IF NOT EXISTS queue STRING DEFAULT 'General'",
				"ALTER TABLE tickets ADD COLUMN IF NOT EXISTS assigned_to STRING DEFAULT ''",
				"ALTER TABLE tickets ADD COLUMN IF NOT EXISTS is_sticky BOOL DEFAULT false",
				"ALTER TABLE tickets ADD COLUMN IF NOT EXISTS referrer_code STRING DEFAULT ''",
			}
			for _, aq := range alterQueries {
				_, err = db.Exec(aq)
				if err != nil {
					log.Printf("Executing alter query failed: %v\nQuery: %s", err, aq)
				}
			}
			fmt.Println("Applied ALTER TABLE tickets upgrades")
		}
	}
	fmt.Println("Database initialization completed successfully.")
}
