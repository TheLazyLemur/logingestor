package logingestor

import (
	"database/sql"
	"log/slog"
)

func ConsumeLogsForever(publisherChannel <-chan []byte, db *sql.DB) {
	for {
		message, ok := <-publisherChannel
		if !ok {
			slog.Error("Failed to get log from channel")
		}

		if err := writeLog(message, db); err != nil {
			slog.Error("Could not upload log, sending to failed queue", "err", err.Error())
		}
	}
}
