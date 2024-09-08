package server

import (
	"github.com/faresbakhit/9b/internal/config"
	"github.com/faresbakhit/9b/internal/store"
)

type Server struct {
	store *store.Store
}

func New() (Server, error) {
	store, err := store.New(config.SQLITE_DATA_SOURCE)
	if err != nil {
		return Server{}, err
	}
	return Server{store}, nil
}

func (s *Server) Close() {
	s.store.Close()
}
