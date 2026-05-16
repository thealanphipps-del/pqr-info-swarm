package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/thealanphipps-del/pqr"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/auth"
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
	aiService := service.NewAIService()
	swarmService := service.NewSwarmService(repo, repo)
	healingService := service.NewHealingService(repo, swarmService, aiService)

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

	// 4. Initialize Auth Service (SAML IdP)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://pqr.info"
	}
	
	certPath := "certs/origin_ca.pem"
	keyPath := "certs/origin_ca.key"
	
	var privKey *rsa.PrivateKey
	var cert *x509.Certificate
	
	if _, err := os.Stat(certPath); err == nil {
		log.Printf("✓ Loading Cloudflare Origin CA Certificate from %s", certPath)
		privKey, cert, err = auth.LoadCertFromFiles(certPath, keyPath)
		if err != nil {
			log.Printf("⚠️ Error: Failed to load Origin CA: %v. Falling back to self-signed.", err)
			privKey, cert, _ = auth.GenerateSelfSignedCert("pqr.info")
		}
	} else {
		log.Printf("ℹ️ Origin CA not found at %s. Generating self-signed cert.", certPath)
		privKey, cert, err = auth.GenerateSelfSignedCert("pqr.info")
		if err != nil {
			log.Printf("Warning: Failed to generate SAML cert: %v", err)
		}
	}
	
	authService, err := service.NewAuthService(repo, baseURL, privKey, cert)
	if err != nil {
		log.Printf("Warning: Failed to initialize Auth Service: %v", err)
	} else {
		log.Println("✓ SAML Identity Provider ready at", baseURL+"/saml/metadata")
	}

	// 5. Start Monitoring Service (Healing Agent)
	monitor := service.NewMonitoringService(healingService, authService, "pqr.info")
	monitor.Start(context.Background())

	// 6. Start Server
	server := pqr.NewServer(swarmService, healingService, authService, aiService)
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
