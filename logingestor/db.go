package logingestor

import (
	"context"
	"database/sql"
	"sync"

	"logingestor/logingestor/types"

	_ "modernc.org/sqlite"
)

func AddLog(ctx context.Context, tx *sql.Tx, logEntry types.LogEntry) error {
	insertEntryQuery := `
	INSERT INTO log_entries (level, message, resource_id, timestamp, trace_id, span_id, "commit") 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := tx.Exec(insertEntryQuery, logEntry.Level, logEntry.Message, logEntry.ResourceID, logEntry.Timestamp, logEntry.TraceID, logEntry.SpanID, logEntry.Commit)
	if err != nil {
		return err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	insertMetadataQuery := `INSERT INTO log_metadata (log_entry_id, key, value) VALUES (?, ?, ?)`

	for key, value := range logEntry.Metadata {
		if _, err = tx.Exec(insertMetadataQuery, newID, key, value); err != nil {
			return err
		}
	}

	return nil
}

func getLogEntries(ctx context.Context, b backends, message string, traceID string, level string) ([]types.LogEntry, error) {
	selectLogEntryQuery := `
	SELECT * FROM log_entries
	`

	clauses := map[string]string{}
	if message != "" {
		clauses["message"] = "'" + message + "'"
	}

	if traceID != "" {
		clauses["trace_id"] = "'" + traceID + "'"
	}

	if level != "" {
		clauses["level"] = "'" + level + "'"
	}

	count := 0
	for key, value := range clauses {
		if count < 1 {
			selectLogEntryQuery = selectLogEntryQuery + " WHERE " + key + " = " + value
		} else {
			selectLogEntryQuery = selectLogEntryQuery + " AND " + key + " = " + value
		}
		count++
	}

	selectLogEntryQuery = selectLogEntryQuery + " ORDER BY timestamp LIMIT 2000;"

	var entries []types.LogEntry
	rows, err := b.GetDB().QueryContext(ctx, selectLogEntryQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var logID int64
		var log types.LogEntry
		err := rows.Scan(&logID, &log.Level, &log.Message, &log.ResourceID, &log.Timestamp, &log.TraceID, &log.SpanID, &log.Commit)
		if err != nil {
			return nil, err
		}

		entries = append(entries, log)
	}

	return entries, nil
}

func getLogMetadata(ctx context.Context, db *sql.DB) ([]types.LogMetadata, error) {
	selectLogEntryQuery := `
	SELECT * FROM log_metadata LIMIT 200000;
	`

	var entries []types.LogMetadata
	rows, err := db.QueryContext(ctx, selectLogEntryQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log types.LogMetadata
		err := rows.Scan(&log.LogEntryID, &log.Key, &log.Value)
		if err != nil {
			return nil, err
		}

		entries = append(entries, log)
	}

	return entries, nil
}

func GetLogs(ctx context.Context, b backends, message string, traceID string, level string) ([]types.LogEntry, error) {
	errors := make(chan error, 1)
	wg := sync.WaitGroup{}

	wg.Add(1)
	var logEntries []types.LogEntry
	go func() {
		defer wg.Done()
		entries, err := getLogEntries(ctx, b, message, traceID, level)
		if err != nil {
			errors <- err
		}

		logEntries = entries
	}()

	wg.Add(1)
	var logMetadataEntries []types.LogMetadata
	go func() {
		defer wg.Done()
		entries, err := getLogMetadata(ctx, b.GetDB())
		if err != nil {
			errors <- err
		}

		logMetadataEntries = entries
		_ = logMetadataEntries
	}()

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return logEntries, nil
}
