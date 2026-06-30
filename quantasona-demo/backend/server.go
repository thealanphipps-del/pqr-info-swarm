package backend

import (
	"log"
	"net/http"
	"sovereign-node-go/pkg/rtgo"
)

func contextHandler(w http.ResponseWriter, r *http.Request) {
	// Mock a ticket since we just need to return formatted context for R3
	ticket := rtgo.TicketSummary{
		ID:         "123e4567-e89b-12d3-a456-426614174000",
		Title:      "Integrate LLM",
		Status:     rtgo.StatusPending,
		RawContent: "Need to integrate LLM to read existing code and tickets. The quantum resonance chamber is operational.",
	}
	
	formattedContext := FormatTicketForLLM(ticket)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(formattedContext))
}

func StartServer() {
	http.HandleFunc("/api/context/raw", contextHandler)
	log.Println("Starting server on :3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
