package store

import (
	"database/sql"
	_ "embed"
	"github.com/faresbakhit/9b/config"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

//go:embed schema.sql
var sqlSchema string

func NewStore() (Store, error) {
	db, err := sql.Open("sqlite3", config.SQLITE_SOURCE_NAME)
	if err != nil {
		return Store{}, err
	}
	_, err = db.Exec(sqlSchema)
	return Store{db}, nil
}
