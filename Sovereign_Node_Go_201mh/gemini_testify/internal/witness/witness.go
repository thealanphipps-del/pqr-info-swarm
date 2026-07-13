package witness
import (
	"crypto/sha256"
	"fmt"
	"time"
)
// Evidence represents a sealed forensic fragment
type Evidence struct {
	Timestamp int64
	TicketID  string
	Payload   string
	Hash      string
}
// SealEvidence generates an immutable hash for the forensic trail
func SealEvidence(ticketID string, data string) *Evidence {
	ts := time.Now().Unix()
	raw := fmt.Sprintf("%d|%s|%s", ts, ticketID, data)
	h := sha256.Sum256([]byte(raw))
	return &Evidence{Timestamp: ts, TicketID: ticketID, Payload: data, Hash: fmt.Sprintf("%x", h)}
}
