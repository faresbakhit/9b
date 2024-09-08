package server

import (
	"errors"
	"net/http"

	"github.com/faresbakhit/9b/internal/store"
)

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.store.UserNew(username, []byte(password))

	if err != nil {
		var error string
		if errors.Is(err, store.UserErrUsername) {
			error = "Username not within constraints."
		} else {
			error = "Username is taken."
		}
		http.Error(w, error, http.StatusBadRequest)
		return
	}

	setSessionTokenCookie(w, user.SessionToken.String)

	w.WriteHeader(http.StatusOK)
}
