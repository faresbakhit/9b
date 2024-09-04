package store

import (
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(sqliteDataSource string) (*Store, error) {
	log.Printf("loading database from %q", sqliteDataSource)
	db, err := sql.Open("sqlite3", sqliteDataSource)
	if err != nil {
		return nil, err
	}
	store := Store{db}
	store.migrateDB()
	return &store, nil
}

//go:embed schema/v1.sql
var schemaV1 string

func (s Store) migrateDB() error {
	var schemaVersion int
	schemaVersionRow := s.db.QueryRow("PRAGMA user_version", nil)
	if err := schemaVersionRow.Scan(&schemaVersion); err != nil {
		return err
	}
	log.Printf("database schema version: %d", schemaVersion)
	switch schemaVersion {
	case 0:
		if _, err := s.db.Exec(schemaV1); err != nil {
			return err
		}
		if _, err := s.db.Exec("PRAGMA user_version = 1"); err != nil {
			return err
		}
		log.Printf("database schema version: 1")
		// fallthrough. case 1: fallthrough...
	}
	return nil
}
