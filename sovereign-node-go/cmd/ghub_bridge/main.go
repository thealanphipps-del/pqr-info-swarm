package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/bridge", func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		fmt.Printf("[GHUB] Command received: %s\n", cmd)
		w.Write([]byte("Command routed to Sovereign Mesh"))
	})

	fmt.Println("[GHUB] Bridge Listener active on Port 9191")
	http.ListenAndServe(":9191", nil)
}
