package logingestor

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"logingestor/logingestor/types"

	_ "modernc.org/sqlite"
)

func SetupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:example.db?mode=rwc")
	if err != nil {
		return nil, err
	}

	sqliteSettings := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA automatic_index = ON;",
	}

	for _, s := range sqliteSettings {
		_, err = db.Exec(s)
		if err != nil {
			return nil, err
		}
	}
	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
	if err != nil {
		log.Fatal(err)
	}
	if journalMode == "wal" {
		fmt.Println("WAL mode is enabled.")
	} else {
		fmt.Println("WAL mode is not enabled.")
	}

	logTable := `
	CREATE TABLE IF NOT EXISTS log_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		level TEXT,
		message TEXT,
		resource_id TEXT,
		timestamp TEXT,
		trace_id TEXT,
		span_id TEXT,
		"commit" TEXT
	);

	CREATE TABLE IF NOT EXISTS log_metadata (
		log_entry_id INTEGER,
		key TEXT,
		value TEXT,
		FOREIGN KEY (log_entry_id) REFERENCES log_entries(id)
	);

	CREATE INDEX IF NOT EXISTS idx_log_entries_resource_id ON log_entries(resource_id);
	CREATE INDEX IF NOT EXISTS idx_log_entries_timestamp ON log_entries(timestamp);
	CREATE INDEX IF NOT EXISTS idx_log_entries_trace_id ON log_entries(trace_id);
	CREATE INDEX IF NOT EXISTS idx_log_entries_level ON log_entries(level);
	CREATE INDEX IF NOT EXISTS idx_log_metadata_log_entry_id ON log_metadata(log_entry_id);
	CREATE INDEX IF NOT EXISTS idx_log_metadata_key ON log_metadata(key);
	CREATE INDEX IF NOT EXISTS idx_log_metadata_value ON log_metadata(value);
	`
	_, err = db.Exec(logTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddLog(tx *sql.Tx, logEntry types.LogEntry) error {
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

func getLogEntries(db *sql.DB) ([]types.LogEntry, error) {
	selectLogEntryQuery := `
	SELECT * FROM log_entries;
	`

	var entries []types.LogEntry
	rows, err := db.Query(selectLogEntryQuery)
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

	fmt.Println("entries: ", len(entries))
	return entries, nil
}

func getLogMetadata(db *sql.DB) ([]types.LogMetadata, error) {
	selectLogEntryQuery := `
	SELECT * FROM log_metadata;
	`

	var entries []types.LogMetadata
	rows, err := db.Query(selectLogEntryQuery)
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

	fmt.Println("entries: ", len(entries))
	return entries, nil
}

func GetLogs(db *sql.DB) ([]types.LogEntry, error) {
	errors := make(chan error, 1)
	wg := sync.WaitGroup{}

	wg.Add(1)
	var logEntries []types.LogEntry
	go func() {
		defer wg.Done()
		entries, err := getLogEntries(db)
		if err != nil {
			errors <- err
		}

		logEntries = entries
	}()

	wg.Add(1)
	var logMetadataEntries []types.LogMetadata
	go func() {
		defer wg.Done()
		entries, err := getLogMetadata(db)
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
		} else {
			fmt.Println("worker completed successfully")
		}
	}

	return logEntries, nil
}
