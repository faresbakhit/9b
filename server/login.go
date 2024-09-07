package server

import (
	"fmt"
	"net/http"

	"github.com/faresbakhit/9b/store"
)

func (s *Server) LoginPath() string {
	return "/login"
}

func (s *Server) LoginPOSTPattern() string {
	return "POST " + s.LoginPath()
}

func (s *Server) LoginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionToken, err := s.store.UserUpdateSessionToken(username, []byte(password))

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Invalid username or password.")
		return
	}

	setSessionTokenCookie(w, sessionToken)

	w.WriteHeader(http.StatusOK)
}

func sessionTokenCookieName() string {
	return "session_token"
}

func getSessionTokenCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(sessionTokenCookieName())
}

func setSessionTokenCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionTokenCookieName(),
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) getUser(r *http.Request) *store.User {
	sessionToken, err := getSessionTokenCookie(r)
	if err != nil {
		return nil
	}
	user, err := s.store.UserFromSessionToken(sessionToken.Value)
	if err != nil {
		return nil
	}
	return user
}
