package logingestor

import (
	"io"
	"log/slog"
	"logingestor/logingestor/views"
	"net/http"
)

func LogIngestionHandler(b backends) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}
		defer r.Body.Close()

		submitLog(b, body)
		w.WriteHeader(http.StatusAccepted)
	}
}

func LogViewHandler(b backends) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := views.Home().Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			return
		}
	}
}

func ListLogHandler(b backends) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		message := r.FormValue("message")
		traceID := r.FormValue("traceID")
		level := r.FormValue("level")

		entries, err := GetLogs(r.Context(), b, message, traceID, level)
		if err != nil {
			panic(err)
		}

		views.Logs(entries).Render(r.Context(), w)
	}
}
