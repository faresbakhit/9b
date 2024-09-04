package server

import (
	"errors"
	"net/http"

	"github.com/faresbakhit/9b/store"
	"github.com/faresbakhit/9b/views"
)

func (s *Server) SignupPath() string {
	return "/signup"
}

func (s *Server) SignupGETPattern() string {
	return "GET " + s.SignupPath()
}

func (s *Server) SignupGET(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.URL.Query().Get("goto")

	if user := s.getUser(r); user != nil {
		if gotoValue == "" {
			gotoValue = s.UserPath(user.Username)
		}
		http.Redirect(w, r, gotoValue, http.StatusSeeOther)
		return
	}

	views.Signup(w, views.SignupData{Goto: gotoValue}, http.StatusOK)
}

func (s *Server) SignupPOSTPattern() string {
	return "POST " + s.SignupPath()
}

func (s *Server) SignupPOST(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.FormValue("goto")
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.store.UserNew(username, []byte(password))

	if err != nil {
		data := views.SignupData{Goto: gotoValue}
		if errors.Is(err, store.UserErrUsername) {
			data.Err = "Username not within constraints."
			views.Signup(w, data, http.StatusUnauthorized)
			return
		}
		data.Err = "Username is taken."
		views.Signup(w, data, http.StatusConflict)
		return
	}

	setSessionTokenCookie(w, user.SessionToken.String)

	if gotoValue == "" {
		gotoValue = s.UserPath(username)
	}

	http.Redirect(w, r, gotoValue, http.StatusSeeOther)
}
