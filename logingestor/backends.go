package logingestor

import "database/sql"

type backends interface {
	GetChan() chan []byte
	GetDB() *sql.DB
}
