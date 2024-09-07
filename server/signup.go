package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/faresbakhit/9b/store"
)

func (s *Server) SignupPath() string {
	return "/signup"
}

func (s *Server) SignupPOSTPattern() string {
	return "POST " + s.SignupPath()
}

func (s *Server) SignupPOST(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.store.UserNew(username, []byte(password))

	if err != nil {
		if errors.Is(err, store.UserErrUsername) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Username not within constraints.")
		} else {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Username is taken.")
		}
		return
	}

	setSessionTokenCookie(w, user.SessionToken.String)

	w.WriteHeader(http.StatusOK)
}
