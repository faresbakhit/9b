package server

import "net/http"

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if user := s.getUser(r); user != nil && user.SessionToken.Valid {
		s.store.UserDeleteSessionToken(user.SessionToken.String)
	}
	w.WriteHeader(http.StatusOK)
}
