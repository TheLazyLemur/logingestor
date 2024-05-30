package main

import (
	"database/sql"
	"log"
	"net/http"

	"logingestor/logingestor"

	_ "modernc.org/sqlite"
)

type Backends interface {
	GetChan() *chan []byte
	GetDB() *sql.DB
}

type backendsImpl struct {
	logIngestionChannel chan []byte
	db                  *sql.DB
}

func (b *backendsImpl) GetChan() chan []byte {
	if b.logIngestionChannel == nil {
		b.logIngestionChannel = make(chan []byte, 100)
	}

	return b.logIngestionChannel
}

func (b *backendsImpl) GetDB() *sql.DB {
	if b.db == nil {
		db, err := logingestor.SetupDB()
		if err != nil {
			panic("Failed to setup DB, cannot recover. err: " + err.Error())
		}

		b.db = db
	}

	return b.db
}

func main() {
	b := &backendsImpl{}
	defer b.GetDB().Close()

	defer close(b.GetChan())

	r := http.NewServeMux()

	r.HandleFunc("POST /", logingestor.LogIngestionHandler(b.GetChan()))
	r.HandleFunc("GET /view", logingestor.LogViewHandler())
	r.HandleFunc("GET /logs", logingestor.ListLogHandler(b.GetDB()))

	go logingestor.ConsumeLogsForever(b.GetChan(), b.GetDB())

	log.Fatal(http.ListenAndServe(":3000", r))
}
