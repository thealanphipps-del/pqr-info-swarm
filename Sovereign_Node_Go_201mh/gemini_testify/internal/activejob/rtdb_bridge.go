package activejob
import (
	"database/sql"
	_ "github.com/lib/pq"
	"gemini_testify/internal/model"
)
func CommitTicket(db *sql.DB, ticket model.ForensicTicket) error {
	query := `INSERT INTO forensic_reconstruction (case_id, source, timestamp, payload, hash) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, ticket.CaseID, ticket.Source, ticket.Timestamp, ticket.Payload, ticket.Hash)
	return err
}
