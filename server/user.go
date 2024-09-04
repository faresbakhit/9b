package server

import (
	"net/http"

	"github.com/faresbakhit/9b/views"
)

func (s *Server) UserPath(name string) string {
	// TODO Escape `name`?
	return "/u/" + name
}

func (s *Server) UserGETPattern() string {
	return "GET /u/{name}"
}

func (s *Server) UserGETHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	user, err := s.store.UserFromUsername(name)
	if err != nil {
		views.NotFound(w, "No such user.")
		return
	}
	userPageData := views.UserData{Username: user.Username, CreatedAt: user.CreatedAt}
	if self := s.getUser(r); self != nil && user.Id == self.Id {
		userPageData.IsSelf = true
	}
	views.User(w, userPageData, http.StatusOK)
}

func (s *Server) UserPOSTPattern() string {
	return "POST /u/{name}"
}

func (s *Server) UserPOSTHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		views.Unauthorized(w)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	err := s.store.UserUpdatePassword(user, []byte(currentPassword), []byte(newPassword))
	if err != nil {
		userPageData := views.UserData{
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			IsSelf:    true,
			Err:       "Changing password failed; Invalid login credentials."}
		views.User(w, userPageData, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) ProtectedPage(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		views.Unauthorized(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
