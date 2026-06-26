package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting Swarm Proxy Daemon...")

	broadcaster := NewBroadcaster()

	// Start Broadcaster WS server
	go func() {
		http.HandleFunc("/ws", broadcaster.HandleWS)
		fmt.Println("[Broadcaster] Listening on :8081/ws")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// Discover Yoga Laptop
	fmt.Println("Discovering Yoga Laptop (Ollama) on port 11434...")
	yogaIP, err := findOllamaServer("11434")
	var ollamaTarget string
	if err != nil {
		fmt.Printf("Warning: Could not discover Yoga laptop: %v. Defaulting to localhost:11434\n", err)
		ollamaTarget = "http://127.0.0.1:11434"
	} else {
		fmt.Printf("Found Yoga Laptop at %s\n", yogaIP)
		ollamaTarget = fmt.Sprintf("http://%s:11434", yogaIP)
	}

	// Setup Gemma Proxy (Listen 1234 -> Forward 1233)
	setupReverseProxy("127.0.0.1:1234", "http://127.0.0.1:1233", "Gemma", broadcaster)

	// Setup Ollama Proxy (Listen 11434 -> Forward Yoga Laptop)
	setupReverseProxy("127.0.0.1:11434", ollamaTarget, "Ollama", broadcaster)

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down proxy daemon...")
}
