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

	var posts []*store.UserPost
	for p, err := range s.store.UserPostListToday(10, 0) {
		if err != nil {
			log.Printf("HomeHandler:UserPostListToday:%q", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		posts = append(posts, p)
	}

	user := s.getUser(r)

	views.Home(w, &views.HomeData{User: user, Posts: posts}, http.StatusOK)
}
