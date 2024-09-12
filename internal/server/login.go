package server

import (
	"log"
	"net/http"

	"github.com/faresbakhit/9b/internal/config"
	"github.com/faresbakhit/9b/internal/store"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionToken, err := s.store.UserUpdateSessionToken(username, []byte(password))

	if err != nil {
		http.Error(w, "Invalid username or password.", http.StatusBadRequest)
		return
	}

	setSessionTokenCookie(w, sessionToken)

	w.WriteHeader(http.StatusOK)
}

func getSessionTokenCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(config.SESSION_TOKEN_COOKIE_NAME)
}

func setSessionTokenCookie(w http.ResponseWriter, sessionToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.SESSION_TOKEN_COOKIE_NAME,
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
		log.Print(err)
		return nil
	}
	user, err := s.store.UserFromSessionToken(sessionToken.Value)
	if err != nil {
		log.Print(err)
		return nil
	}
	return user
}
