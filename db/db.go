package db

import (
	"database/sql"
	"log"
	"log/slog"
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
		slog.Info("WAL mode is enabled.")
	} else {
		slog.Warn("WAL mode is not enabled.")
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
