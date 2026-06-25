package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/thealanphipps-del/pqr/internal/execution"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/auth"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
	"github.com/thealanphipps-del/pqr/internal/service"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "goback" {
		handleGoback(os.Args[2:])
		return
	}

	log.Println("[swend] starting Swarm Execution Node...")

	// 1. Load Vault token
	token := os.Getenv("PQR_VAULT_TOKEN")
	if token == "" {
		// Fallback to the known swarm token if not specified in env for dev
		token = os.Getenv("VAULT_TOKEN")
		if token == "" {
			token = "pqr-vault-token" // Default fallback based on sweep_secrets.ps1
		}
	}

	vaultClient, err := auth.NewVaultSecretManager()
	if err != nil {
		log.Fatalf("vault init failed: %v", err)
	}

	if err := vaultClient.VerifyIdentity(context.Background()); err != nil {
		log.Fatalf("vault identity verification failed: %v", err)
	}

	// 2. Build capability manifest
	osType := "linux"
	if runtime.GOOS == "windows" {
		osType = "windows_wsl"
	}

	manifest := service.CapabilityManifest{
		AgentType:    "execution_node",
		OS:           osType,
		Capabilities: []string{"windows_execution", "wsl_execution", "gcp_ops", "emulator_control", "error_solution_learning"},
	}

	// 3. Register with swarm
	swarmClient := service.NewSwarmClient()
	agentInfo, err := swarmClient.Register(manifest)
	if err != nil {
		log.Fatalf("swarm registration failed: %v", err)
	}

	log.Printf("[swend] registered as agent %s (shortcode %s)",
		agentInfo.AgentID, agentInfo.Shortcode)

	// 4. Initialize execution engine
	engine := execution.NewExecutionEngine()

	// 5. Initialize CockroachDB memory
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	}
	memoryDB, err := db.NewCockroachRepository(connStr)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	if err := memoryDB.InitSchema(context.Background()); err != nil {
		log.Fatalf("db schema init failed: %v", err)
	}

	engine.Memory = memoryDB

	// 6. Start execution stream
	go service.StartExecutionStream(swarmClient, engine, memoryDB)

	// 7. Optional command palette
	if len(os.Args) > 1 && os.Args[1] == "menu" {
		execution.StartPalette(engine, memoryDB, func() {
			fmt.Print("Enter timestamp to rewind to (or 'last'): ")
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			choice := strings.TrimSpace(line)
			
			gb := service.NewGobackService(memoryDB)
			if choice == "last" {
				if err := gb.Last(); err != nil {
					fmt.Printf("Undo last failed: %v\n", err)
				} else {
					fmt.Println("Successfully undid last action.")
				}
			} else if choice != "" {
				if err := gb.System(choice); err != nil {
					fmt.Printf("System rewind failed: %v\n", err)
				} else {
					fmt.Printf("Successfully rewound system to %s\n", choice)
				}
			}
		})
	} else {
		log.Println("[swend] execution node running in background. Waiting for swarm commands...")
	}

	select {} // block forever
}

func handleGoback(args []string) {
	flags := flag.NewFlagSet("goback", flag.ExitOnError)
	ts := flags.String("to", "", "timestamp to rewind to")
	last := flags.Bool("last", false, "undo last change")
	chain := flags.String("chain", "", "undo change and all dependent changes")
	fixes := flags.String("fixes", "", "rewind only fix memory to timestamp")

	flags.Parse(args)

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgresql://root@localhost:26257/antigravity?sslmode=disable"
	}
	memoryDB, err := db.NewCockroachRepository(connStr)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	gb := service.NewGobackService(memoryDB)

	switch {
	case *last:
		if err := gb.Last(); err != nil {
			log.Fatalf("undo last failed: %v", err)
		}
		fmt.Println("Successfully undid last action.")
	case *ts != "":
		if err := gb.System(*ts); err != nil {
			log.Fatalf("system rewind failed: %v", err)
		}
		fmt.Println("Successfully rewound system to", *ts)
	case *chain != "":
		if err := gb.Chain(*chain); err != nil {
			log.Fatalf("chain undo failed: %v", err)
		}
		fmt.Println("Successfully undid chain from", *chain)
	case *fixes != "":
		if err := gb.Fixes(*fixes); err != nil {
			log.Fatalf("fixes rewind failed: %v", err)
		}
		fmt.Println("Successfully rewound fixes to", *fixes)
	default:
		fmt.Println("Usage: swend goback --to <timestamp> | --last | --chain <timestamp> | --fixes <timestamp>")
	}
}
