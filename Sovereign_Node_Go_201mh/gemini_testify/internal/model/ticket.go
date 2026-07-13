package model
import "time"
type ForensicTicket struct {
	ID        int       `db:"id"`
	CaseID    string    `db:"case_id"`
	Source    string    `db:"source"`
	Timestamp time.Time `db:"timestamp"`
	Payload   string    `db:"payload"`
	Hash      string    `db:"hash"`
}
