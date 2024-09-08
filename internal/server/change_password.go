package server

import (
	"fmt"
	"net/http"
)

func (s *Server) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized.")
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	err := s.store.UserUpdatePassword(user, []byte(currentPassword), []byte(newPassword))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Incorrect password.")
		return
	}

	w.WriteHeader(http.StatusOK)
}
