package main

import (
	"database/sql"
	"logingestor/db"
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
		db, err := db.SetupDB()
		if err != nil {
			panic("Failed to setup DB, cannot recover. err: " + err.Error())
		}

		b.db = db
	}

	return b.db
}
