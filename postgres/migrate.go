package postgres

import (
	"log"

	"github.com/pressly/goose"
)

func migrate(store *Storage) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(store.db.DB, "migrations"); err != nil {
		log.Fatal(err)
	}
}
