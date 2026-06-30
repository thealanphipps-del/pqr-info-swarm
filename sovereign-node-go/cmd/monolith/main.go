package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sovereign-node-go/pkg/github"
	"sovereign-node-go/pkg/mesh"
	"sovereign-node-go/pkg/tickets"
)

var (
	strikeMgr *mesh.StrikeManager
)

func main() {
	fmt.Println("Sovereign Super Monolith v7.2: Initializing Semantic Strike Pipeline...")

	// 1. Initialize Infrastructure
	client := github.NewClient("")
	ticketMgr, err := tickets.NewManager(client, "owner", "repo")
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize ticket manager: %v", err)
	}

	gov := mesh.NewGovernance(ticketMgr)
	strikeMgr = mesh.NewStrikeManager(ticketMgr, gov)

	// 2. HTTP Routes
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	http.HandleFunc("/strike/propose", handlePropose)
	http.HandleFunc("/strike/vote", handleVote)
	http.HandleFunc("/strike/status", handleStatus)

	addr := "0.0.0.0:11111"
	fmt.Printf("[SUCCESS] Monolith HUD active on http://localhost:11111\n")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handlePropose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AgentID string                `json:"agent_id"`
		Layer   int                   `json:"layer"`
		Cluster mesh.ClusterID        `json:"cluster"`
		Content tickets.FabricContent `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strike, err := strikeMgr.ProposeStrike(context.Background(), req.AgentID, req.Layer, req.Cluster, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(strike)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement voting logic
	fmt.Fprintf(w, "Voting implementation pending forensic audit.\n")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement status check
	fmt.Fprintf(w, "Status check implementation pending forensic audit.\n")
}
