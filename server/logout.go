package server

import "net/http"

func (s *Server) LogoutPath() string {
	return "/logout"
}

func (s *Server) LogoutPOSTPattern() string {
	return "POST " + s.LogoutPath()
}

func (s *Server) LogoutPOSTHandler(w http.ResponseWriter, r *http.Request) {
	if user := s.getUser(r); user != nil && user.SessionToken.Valid {
		s.store.UserDeleteSessionToken(user.SessionToken.String)
	}
	w.WriteHeader(http.StatusOK)
}
