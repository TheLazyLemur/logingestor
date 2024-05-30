package logingestor

import (
	"context"
	"database/sql"
	"log/slog"
)

func ConsumeLogsForever(ctx context.Context, publisherChannel <-chan []byte, db *sql.DB) {
	for {
		message, ok := <-publisherChannel
		if !ok {
			slog.Error("Failed to get log from channel")
		}

		if err := writeLog(ctx, message, db); err != nil {
			slog.Error("Could not upload log, sending to failed queue", "err", err.Error())
		}
	}
}
