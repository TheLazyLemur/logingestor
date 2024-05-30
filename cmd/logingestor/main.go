package main

import (
	"context"
	"log"
	"net/http"

	"logingestor/logingestor"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

func traceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "trace", uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	ctx := context.Background()

	b := &backendsImpl{}
	defer b.GetDB().Close()
	defer close(b.GetChan())

	r := chi.NewRouter()
	r.Use(traceIDMiddleware)

	r.Post("/", logingestor.LogIngestionHandler(b))
	r.Get("/view", logingestor.LogViewHandler(b))
	r.Get("/logs", logingestor.ListLogHandler(b))

	go logingestor.ConsumeLogsForever(ctx, b.GetChan(), b.GetDB())

	log.Fatal(http.ListenAndServe(":3000", r))
}
