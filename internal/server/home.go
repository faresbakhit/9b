package server

import (
	"log"
	"net/http"

	"github.com/faresbakhit/9b/internal/store"
	"github.com/faresbakhit/9b/internal/views"
)

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Get("search")
	r.URL.Query().Get("type")
	r.URL.Query().Get("sortby")

	var userId int

	user := s.getUser(r)
	if user == nil {
		userId = 0
	} else {
		userId = user.Id
	}

	var posts []*store.UserPostListResult
	for p, err := range s.store.UserPostList(userId, 10, 0) {
		if err != nil {
			log.Printf("HomeHandler:UserPostListToday:%q", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	data := &views.HomeData{LoggedIn: user != nil, Posts: posts}
	views.Home(w, data, http.StatusOK)
}
