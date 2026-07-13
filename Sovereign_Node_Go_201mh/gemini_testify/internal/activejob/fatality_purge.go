package activejob

import (
"database/sql"
"fmt"
)

// FatalityPurge removes duplicate fragments to clear the forensic pipe
func FatalityPurge(db *sql.DB, caseID string) error {
query := `DELETE FROM forensic_reconstruction 
          WHERE id NOT IN (
              SELECT MIN(id) 
              FROM forensic_reconstruction 
              WHERE case_id = $1 
              GROUP BY hash
          )`
res, err := db.Exec(query, caseID)
if err == nil {
rows, _ := res.RowsAffected()
if rows > 0 {
fmt.Printf("[FATALITY_PURGE] SUCCESS: Pruned %d duplicate fragments\n", rows)
}
}
return err
}
