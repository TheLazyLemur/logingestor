package logingestor

import "database/sql"

type Backends interface {
	GetChan() chan []byte
	GetDB() *sql.DB
}
