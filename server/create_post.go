package server

import (
	"log"
	"net/http"

	"github.com/faresbakhit/9b/store"
	"github.com/faresbakhit/9b/views"
)

func (s *Server) CreatePostPath() string {
	return "/post"
}

func (s *Server) CreatePostGETPattern() string {
	return "GET " + s.CreatePostPath()
}

func (s *Server) CreatePostGETHandler(w http.ResponseWriter, r *http.Request) {
	if s.getUser(r) == nil {
		http.Redirect(w, r, s.LoginPathWithGoto(s.CreatePostPath()), http.StatusSeeOther)
		return
	}
	views.CreatePost(w, views.CreatePostData{}, http.StatusOK)
}

func (s *Server) CreatePostPOSTPattern() string {
	return "POST " + s.CreatePostPath()
}

func (s *Server) CreatePostPOSTHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		views.Unauthorized(w)
		return
	}
	postTitle := r.FormValue("title")
	postUrl := r.FormValue("url")
	postBody := r.FormValue("body")
	if postTitle == "" {
		data := views.CreatePostData{Err: "Invalid submission; Title is empty."}
		views.CreatePost(w, data, http.StatusUnprocessableEntity)
		return
	}
	if postUrl == "" && postBody == "" {
		data := views.CreatePostData{
			Err: "Invalid Submission; Must have one or both of link and body."}
		views.CreatePost(w, data, http.StatusUnprocessableEntity)
		return
	}
	if err := s.store.UserPostNew(&store.UserPost{
		UserId: user.Id,
		Title:  postTitle,
		Url:    postUrl,
		Body:   postBody,
	}); err != nil {
		log.Printf("new post internal error: %q", err)
		data := views.CreatePostData{Err: "Internal server error."}
		views.CreatePost(w, data, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, s.UserPath(user.Username), http.StatusSeeOther)
}
