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

func New(sqliteDataSource string) (*Store, error) {
	log.Printf("loading database from %q", sqliteDataSource)
	db, err := sql.Open("sqlite3", sqliteDataSource)
	if err != nil {
		return nil, err
	}

	store := Store{db}
	store.migrateDB()

	var foreignKeys int
	db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	log.Printf("PRAGMA foreign_keys = %d", foreignKeys)

	return &store, nil
}

func (s *Store) Close() {
	s.db.Close()
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
		log.Printf("database schema version: 1")
		// fallthrough. case 1: fallthrough...
	}
	return nil
}
