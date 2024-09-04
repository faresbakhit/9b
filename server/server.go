package server

import (
	"github.com/faresbakhit/9b/config"
	"github.com/faresbakhit/9b/store"
)

type Server struct {
	store *store.Store
}

func NewServer() (Server, error) {
	store, err := store.NewStore(config.SQLITE_DATA_SOURCE)
	if err != nil {
		return Server{}, err
	}
	return Server{store}, nil
}
