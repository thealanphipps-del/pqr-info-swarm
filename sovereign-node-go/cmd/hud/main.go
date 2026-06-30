package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sovereign-node-go/pkg/agent"
	"sovereign-node-go/pkg/llm"
	"sovereign-node-go/pkg/rtgo"
	"time"
)

var rtgoManager *rtgo.Manager

func main() {
	// Initialize Local LLM Gateway
	llmGateway := llm.NewLocalGateway(llm.DefaultConfig)

	connStr := os.Getenv("DATABASE_URL")
	var err error
	rtgoManager, err = rtgo.NewManager(connStr, llmGateway)
	if err != nil {
		log.Fatalf("Failed to initialize rtgo manager: %v", err)
	}

	// Start Agent Background Worker to process tickets
	worker := agent.NewTicketWorker(rtgoManager, llmGateway, nil) // no specific agents assigned yet
	go worker.Start(context.Background())

	rtgoServer := rtgo.NewServer(rtgoManager)

	// Static files
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/app.js", serveStatic)
	http.HandleFunc("/style.css", serveStatic)

	// Real-time stream
	http.HandleFunc("/api/tickets/stream", handleSSE)

	// Mount RTGO REST 2.0 API (using gin router as a handler)
	http.Handle("/REST/2.0/", rtgoServer.Router)

	fmt.Println("Sovereign HUD active on http://localhost:8082")
	fmt.Println("REST 2.0 API available at http://localhost:8082/REST/2.0/")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web"+r.URL.Path)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			tix, err := rtgoManager.FetchTickets(ctx)
			if err != nil {
				log.Printf("[SSE ERROR] %v", err)
				fmt.Fprintf(w, "event: error\ndata: %v\n\n", err)
			} else {
				data, _ := json.Marshal(tix)
				fmt.Fprintf(w, "data: %s\n\n", data)
			}
			flusher.Flush()
			time.Sleep(2 * time.Second)
		}
	}
}
