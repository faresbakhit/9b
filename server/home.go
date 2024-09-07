package server

import (
	"net/http"

	"github.com/faresbakhit/9b/views"
)

func (s *Server) HomePath() string {
	return "/home"
}

func (s *Server) HomeGETPattern() string {
	return "GET " + s.HomePath()
}

func (s *Server) HomeGET(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	views.Home(w, &views.HomeData{User: user}, http.StatusOK)
}
