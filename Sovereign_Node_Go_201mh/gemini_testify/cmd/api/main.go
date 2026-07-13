package main
import (
    "database/sql"
    "fmt"
    "gemini_testify/internal/activejob"
    "gemini_testify/internal/controller"
    _ "github.com/lib/pq"
    "time"
)
func main() {
    fmt.Println("[MASTER] Sovereign Vitality Loop Ignited")
    connStr := "postgres://postgres@204.168.138.60:5432/sovereign_forensics?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil { fmt.Printf("[FATAL] DB_BRIDGE_FAILURE: %v\n", err); return }
    
    monitor := &controller.VitalityMonitor{}
    go monitor.Start(func() int {
        var count int
        db.QueryRow("SELECT count(*) FROM forensic_reconstruction WHERE case_id='0176811576'").Scan(&count)
        return count
    })
    
    for {
        activejob.FatalityPurge(db, "0176811576")
        time.Sleep(5 * time.Minute)
    }
}
