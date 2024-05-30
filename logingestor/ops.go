package logingestor

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"logingestor/logingestor/types"
)

func submitLog(b backends, message []byte) {
	b.GetChan() <- message
}

func writeLog(ctx context.Context, message []byte, db *sql.DB) (err error) {
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
		if err = AddLog(ctx, tx, logEntry); err == nil {
			break
		}
		if retries == maxRetries {
			return err
		}

		log.Printf("Retrying to add log entry: attempt %d", retries+1)
	}

	return nil
}
