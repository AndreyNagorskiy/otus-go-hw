package sqlstorage

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/pressly/goose"
	"log"
)

func Migrate(dsn string, retryWithoutSSL bool) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open connection to database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose: failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "internal/storage/migrations"); err != nil {
		if errors.Is(err, pq.ErrSSLNotSupported) && !retryWithoutSSL {
			Migrate(dsn+"?sslmode=disable", true)
		} else {
			log.Fatalf("goose up error: %v", err)
		}
	}
}
