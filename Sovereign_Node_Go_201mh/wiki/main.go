package main
import (
	"fmt"
	"net/http"
	"os"
)
func main() {
	port := os.Getenv("TARGET_PORT")
	if port == "" { port = "9111" }
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		c, _ := os.ReadFile("content.txt")
		fmt.Fprintf(w, "[SOVEREIGN_WIKI_V4.3]\nPORT: %s\n\n%s\n", port, string(c))
	})
	fmt.Printf("[IGNITION] Wiki Live on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
