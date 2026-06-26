package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/thealanphipps-del/pqr"
	"github.com/thealanphipps-del/pqr/internal/execution"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/auth"
	"github.com/thealanphipps-del/pqr/internal/infrastructure/db"
	"github.com/thealanphipps-del/pqr/internal/service"
	"github.com/thealanphipps-del/pqr/internal/tui"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "goback" {
			handleGoback(os.Args[2:])
			return
		}
		if os.Args[1] == "install" || os.Args[1] == "genesis" {
			handleInstall()
			return
		}
		if os.Args[1] == "monitor" {
			tui.StartTUI()
			return
		}
		if os.Args[1] == "rt" {
			handleRT(os.Args[2:])
			return
		}
		if os.Args[1] == "swarm" {
			handleSwarm(os.Args[2:])
			return
		}
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
		log.Printf("[WARNING] vault identity verification failed (proceeding without vault): %v", err)
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

	swarmClient := service.NewSwarmClient()
	agentInfo, err := swarmClient.Register(manifest)
	if err != nil {
		log.Printf("[WARNING] swarm registration failed (proceeding in local mode): %v", err)
	} else {
		log.Printf("[swend] registered as agent %s (shortcode %s)",
			agentInfo.AgentID, agentInfo.Shortcode)
	}

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

	// 5.5 Start Dashboard Service
	gobackSvc := service.NewGobackService(memoryDB)
	dashboard := service.NewDashboardService(memoryDB, gobackSvc)
	dashboard.Start()
	log.Println("[swend] Observability Dashboard started on http://127.0.0.1:7777")

	// 5.7 Start Knowledge Ingestion
	ingester := service.NewKnowledgeIngester(memoryDB)
	ingester.Start(context.Background())
	log.Println("[swend] Knowledge Base Ingestion started")

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

func handleInstall() {
	log.Println("[swend] initializing Genesis Restore Point...")

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

	gb := service.NewGobackService(memoryDB)

	hasWSL := false
	if runtime.GOOS == "windows" {
		hasWSL = true 
	}

	engineSnap := &service.EngineSnapshot{
		OSName:           runtime.GOOS,
		HasWSL:           hasWSL,
		DefaultTimeout:   60,
		EmulatorDefaults: map[string]string{"gpu_mode": "auto", "cold_boot": "false"},
		VaultAddress:     os.Getenv("VAULT_ADDR"),
		VaultRole:        "execution_node",
	}

	configSnap := &service.ConfigSnapshot{
		NodeID:       "swen-node-1",
		SwarmAddress: os.Getenv("PQR_SWARM_ADDR"),
		CockroachDSN: connStr,
		Env: map[string]string{
			"VAULT_ADDR":     os.Getenv("VAULT_ADDR"),
			"VAULT_TOKEN":    os.Getenv("VAULT_TOKEN"),
			"SWARM_ENDPOINT": os.Getenv("PQR_SWARM_ADDR"),
			"SWEN_NODE_ID":   "swen-node-1",
			"ANDROID_HOME":   os.Getenv("ANDROID_HOME"),
			"PATH":           os.Getenv("PATH"),
		},
	}

	if err := gb.Genesis(engineSnap, configSnap, "proto/swarm.proto"); err != nil {
		log.Fatalf("genesis failed: %v", err)
	}

	log.Println("[swend] installation complete. Genesis snapshot created.")
}

func handleRT(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: swend rt <list|show|take|assign|comment|priority>")
		return
	}

	apiURL := os.Getenv("PQR_API_URL")
	if apiURL == "" {
		apiURL = "http://127.0.0.1:8196"
	}
	client := pqr.NewClient(apiURL)
	ctx := context.Background()

	switch args[0] {
	case "list":
		tickets, err := client.ListTickets(ctx)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("\n=== PQR FABRIC TICKETS ===")
		fmt.Printf("%-36s | %-12s | %-12s | %-10s | %s\n", "TICKET ID", "CREATOR", "ASSIGNEE", "STATUS", "CREATED AT")
		fmt.Println(strings.Repeat("-", 90))
		for _, t := range tickets {
			id, _ := t["id"].(string)
			creator, _ := t["creator"].(string)
			assignee, _ := t["assigned_to"].(string)
			status, _ := t["status"].(string)
			created, _ := t["created_at"].(string)
			fmt.Printf("%-36s | %-12s | %-12s | %-10s | %s\n", id, creator, assignee, status, created)
		}
	case "show":
		if len(args) < 2 {
			log.Fatal("Usage: swend rt show <id>")
		}
		t, err := client.GetTicket(ctx, args[1])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("\n=== TICKET DETAILS ===")
		fmt.Printf("ID:       %s\n", t["id"])
		fmt.Printf("Layer:    %v\n", t["layer"])
		fmt.Printf("Creator:  %s\n", t["creator"])
		fmt.Printf("Assignee: %s\n", t["assigned_to"])
		fmt.Printf("Status:   %s\n", t["status"])
		fmt.Printf("Created:  %s\n", t["created_at"])
		fmt.Println("\n--- CONTENT ---")
		fmt.Println(t["content"])
		
		if intent, ok := t["intent"].(map[string]interface{}); ok {
			fmt.Println("\n--- METADATA (INTENT BLOB) ---")
			for k, v := range intent {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	case "take":
		if len(args) < 2 {
			log.Fatal("Usage: swend rt take <id>")
		}
		err := client.UpdateTicketExtended(ctx, args[1], "", "", "operator", "operator", "", "")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("Ticket successfully taken by operator.")
	case "assign":
		if len(args) < 3 {
			log.Fatal("Usage: swend rt assign <id> <agent-id>")
		}
		err := client.UpdateTicketExtended(ctx, args[1], "", "", "operator", args[2], "", "")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Ticket successfully assigned to agent %s.\n", args[2])
	case "comment":
		if len(args) < 3 {
			log.Fatal("Usage: swend rt comment <id> \"comment text\"")
		}
		err := client.CommentTicket(ctx, args[1], "operator", args[2])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("Comment successfully posted.")
	case "priority":
		if len(args) < 3 {
			log.Fatal("Usage: swend rt priority <id> <low|medium|high>")
		}
		err := client.UpdateTicketExtended(ctx, args[1], "", "", "operator", "", args[2], "")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Printf("Ticket priority updated to %s.\n", args[2])
	default:
		fmt.Println("Unknown command. Usage: swend rt <list|show|take|assign|comment|priority>")
	}
}

func handleSwarm(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: swend swarm <status|stop|kill>")
		return
	}

	switch args[0] {
	case "status":
		fmt.Println("Checking legacy swarm status...")
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-Command", "Get-Process -Name python,gemma,inference -ErrorAction SilentlyContinue")
		} else {
			cmd = exec.Command("sh", "-c", "ps aux | grep -E 'gemma|inference|swarm' | grep -v grep")
		}
		out, err := cmd.CombinedOutput()
		if err != nil || len(out) == 0 {
			fmt.Println("No active legacy processes detected.")
			return
		}
		fmt.Println(string(out))
	case "stop", "kill":
		fmt.Println("Executing safe kill-switch for legacy swarm processes...")
		tui.ExecKillSwitch()
		fmt.Println("✓ Legacy swarm processes terminated successfully.")
	default:
		fmt.Println("Unknown command. Usage: swend swarm <status|stop|kill>")
	}
}
