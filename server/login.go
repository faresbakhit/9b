package server

import (
	"net/http"
	"net/url"

	"github.com/faresbakhit/9b/store"
	"github.com/faresbakhit/9b/views"
)

func (s *Server) LoginPath() string {
	return "/login"
}

func (s *Server) LoginPathWithGoto(gotoPath string) string {
	return s.LoginPath() + "?goto=" + url.QueryEscape(gotoPath)
}

func (s *Server) LoginGETPattern() string {
	return "GET " + s.LoginPath()
}

func (s *Server) LoginGETHandler(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.URL.Query().Get("goto")

	if user := s.getUser(r); user != nil {
		if gotoValue == "" {
			gotoValue = s.UserPath(user.Username)
		}
		http.Redirect(w, r, gotoValue, http.StatusSeeOther)
		return
	}

	views.Login(w, views.LoginData{Goto: gotoValue}, http.StatusOK)
}

func (s *Server) LoginPOSTPattern() string {
	return "POST " + s.LoginPath()
}

func (s *Server) LoginPOSTHandler(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.FormValue("goto")
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionToken, err := s.store.UserUpdateSessionToken(username, []byte(password))

	if err != nil {
		data := views.LoginData{Goto: gotoValue, Err: "Login failed; Invalid username or password."}
		views.Login(w, data, http.StatusConflict)
		return
	}

	setSessionTokenCookie(w, sessionToken)

	if gotoValue == "" {
		gotoValue = s.UserPath(username)
	}

	http.Redirect(w, r, gotoValue, http.StatusSeeOther)
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
