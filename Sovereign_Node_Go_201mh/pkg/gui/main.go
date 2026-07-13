package main
import (
	"fmt"
	"net/http"
)
func main() {
	fs := http.FileServer(http.Dir("/data/data/com.termux/files/home/downloads/Takeout"))
	http.Handle("/", fs)
	fmt.Println("GSH-MESH: Takeout Browser Active on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
