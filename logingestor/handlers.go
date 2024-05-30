package logingestor

import (
	"database/sql"
	"io"
	"log/slog"
	"logingestor/logingestor/types"
	"logingestor/logingestor/views"
	"net/http"
)

func LogIngestionHandler(ch chan []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}
		defer r.Body.Close()

		ch <- body
		w.WriteHeader(http.StatusAccepted)
	}
}

func LogViewHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Home().Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}
	}
}

func ListLogHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := GetLogs(db)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			views.LogEntry(types.LogEntry{
				Level:   entry.Level,
				Message: entry.Message,
			}).Render(r.Context(), w)
		}
	}
}
