package server

import (
	"net/http"

	"github.com/faresbakhit/9b/store"
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
	var posts []*store.UserPost
	for post := range s.store.UserPostListFromUser(user.Id, 10, 0) {
		posts = append(posts, post)
	}
	self := s.getUser(r)
	userPageData := views.UserData{
		Posts:     posts,
		User:      self,
		UserOther: user}
	if self != nil && user.Id == self.Id {
		userPageData.IsSelf = true
	}
	views.User(w, userPageData, http.StatusOK)
}
