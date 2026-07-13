package activejob

import (
"bufio"
"crypto/sha256"
"database/sql"
"fmt"
"gemini_testify/internal/model"
"os"
"time"
)

func ParseAndCommitHistory(db *sql.DB, filePath string, caseID string) {
file, err := os.Open(filePath)
if err != nil {
return
}
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
line := scanner.Text()
hash := fmt.Sprintf("%x", sha256.Sum256([]byte(line)))

ticket := model.ForensicTicket{
CaseID:    caseID,
Source:    "wiki_mud_log",
Timestamp: time.Now(),
Payload:   line,
Hash:      hash,
}

// Executing Commit via RTdb Bridge
err := CommitTicket(db, ticket)
if err != nil {
fmt.Printf("[CHILD3] COMMIT_ERROR: %s | HASH: %s\n", err, hash[:8])
} else {
fmt.Printf("[CHILD3] COMMIT_SUCCESS: %s\n", hash[:8])
}
}
}
