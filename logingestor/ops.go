package logingestor

import (
	"database/sql"
	"encoding/json"
	"log"
	"logingestor/logingestor/types"
)

func writeLog(message []byte, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	logEntry := types.LogEntry{}
	if err := json.Unmarshal(message, &logEntry); err != nil {
		return err
	}

	const maxRetries = 3
	for retries := 0; retries <= maxRetries; retries++ {
		if err = AddLog(tx, logEntry); err == nil {
			break
		}
		if retries == maxRetries {
			return err
		}

		log.Printf("Retrying to add log entry: attempt %d", retries+1)
	}

	return nil
}
