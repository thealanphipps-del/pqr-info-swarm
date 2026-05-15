package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/thealanphipps-del/pqr"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
	"github.com/thealanphipps-del/pqr/internal/service"
)

const (
	Version = "v1.02"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	}

	// 1. Initialize Infrastructure
	repo, err := db.NewCockroachRepository(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// 2. Initialize Service Layer
	swarmService := service.NewSwarmService(repo, repo)
	healingService := service.NewHealingService(repo, swarmService)

	// Start Background Healing Worker
	healingService.StartBackgroundWorker(context.Background())

	// 3. Initialize Schema with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := repo.InitSchema(ctx); err != nil {
		cancel()
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	
	// Seed Governance Agents and Initial Tickets
	if err := repo.SeedGovernanceAgents(ctx); err != nil {
		log.Printf("Warning: Failed to seed agents: %v", err)
	}
	if err := repo.SeedTickets(ctx); err != nil {
		log.Printf("Warning: Failed to seed tickets: %v", err)
	}
	cancel()

	log.Println("✓ Database schema initialized")
	log.Println("✓ PQR Ticketing Fabric ready")

	// 4. Start Server
	server := pqr.NewServer(swarmService, healingService)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8196"
	}
	if len(port) > 0 && port[0] != ':' {
		port = ":" + port
	}
	
	log.Printf("Starting PQR REST 2.0 API Server on %s...", port)
	
	if err := server.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
