package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/faresbakhit/9b/render"
	"github.com/faresbakhit/9b/store"
)

type Server struct {
	store store.Store
}

func NewServer() (Server, error) {
	store, err := store.NewStore()
	if err != nil {
		return Server{}, err
	}
	return Server{store}, nil
}

func (s *Server) SignupPage(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.URL.Query().Get("goto")

	if user := s.getUser(r); user != nil {
		http.Redirect(w, r, gotoValue, http.StatusSeeOther)
		return
	}

	render.SignupPage(w, render.SignupPageData{Goto: gotoValue}, http.StatusOK)
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.FormValue("goto")
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.store.CreateUser(username, []byte(password))

	if err != nil {
		data := render.SignupPageData{Goto: gotoValue}
		if errors.Is(err, store.ErrInvalidUsername) {
			data.Err = "Username not within constraints."
			render.SignupPage(w, data, http.StatusUnauthorized)
			return
		}
		data.Err = fmt.Sprintf("Username '%s' is taken.", username)
		render.SignupPage(w, data, http.StatusConflict)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.SessionToken.String,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if gotoValue == "" {
		gotoValue = "/u/" + username
	}

	http.Redirect(w, r, gotoValue, http.StatusSeeOther)
}

func (s *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.URL.Query().Get("goto")

	if user := s.getUser(r); user != nil {
		http.Redirect(w, r, gotoValue, http.StatusSeeOther)
		return
	}

	render.LoginPage(w, render.LoginPageData{Goto: gotoValue}, http.StatusOK)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	gotoValue := r.FormValue("goto")
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionToken, err := s.store.UpdateUserSessionToken(username, []byte(password))

	if err != nil {
		data := render.LoginPageData{Goto: gotoValue, Err: "Login failed; Invalid username or password."}
		render.LoginPage(w, data, http.StatusConflict)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if gotoValue == "" {
		gotoValue = "/u/" + username
	}

	http.Redirect(w, r, gotoValue, http.StatusSeeOther)
}

func (s *Server) UserPage(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	user, err := s.store.GetUserByName(name)
	if err != nil {
		render.NotFound(w, "No such user.")
		return
	}
	userPageData := render.UserPageData{Name: user.Name, JoinedAt: user.CreatedAt}
	if self := s.getUser(r); self != nil && user.Id == self.Id {
		userPageData.IsSelf = true
	}
	render.UserPage(w, userPageData, http.StatusOK)
}

func (s *Server) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		render.Unauthorized(w)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	err := s.store.UpdateUserPassword(user, []byte(currentPassword), []byte(newPassword))
	if err != nil {
		userPageData := render.UserPageData{
			Name:     user.Name,
			JoinedAt: user.CreatedAt,
			IsSelf:   true,
			Err:      "Changing password failed; Invalid login credentials."}
		render.UserPage(w, userPageData, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		render.Unauthorized(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getUser(r *http.Request) *store.User {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	user, err := s.store.GetUserBySessionToken(sessionToken.Value)
	if err != nil {
		return nil
	}
	return user
}
