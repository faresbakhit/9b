package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/faresbakhit/9b/internal/config"
	"github.com/faresbakhit/9b/internal/store"
	"github.com/faresbakhit/9b/internal/views"
)

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	post := s.store.PostGet(user.Id, postId)
	data := &views.PostData{LoggedIn: user != nil, Post: post}
	views.Post(w, data, http.StatusOK)
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	if user == nil {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}
	postTitle := r.FormValue("title")
	postURL := r.FormValue("url")
	postBody := r.FormValue("body")
	if postTitle == "" {
		http.Error(w, "Empty title.", http.StatusBadRequest)
		return
	}
	if postURL == "" && postBody == "" {
		http.Error(w, "Must have one or both of link and body.", http.StatusBadRequest)
		return
	}
	if postURL != "" {
		postURLParse, err := url.Parse(postURL)
		if err != nil || postURLParse.Host == "" {
			http.Error(w, "Invalid link address.", http.StatusBadRequest)
			return
		}
		if postURLParse.Scheme == "http" && !config.ALLOW_HTTP_URLS {
			error := "Insecure HTTP links are not allowed, only HTTPS links are allowed."
			http.Error(w, error, http.StatusBadRequest)
			return
		}
		if postURLParse.Scheme != "https" {
			if config.ALLOW_HTTP_URLS {
				http.Error(w, "Only HTTP/HTTPS links are allowed.", http.StatusBadRequest)
			} else {
				http.Error(w, "Only HTTPS links are allowed.", http.StatusBadRequest)
			}
			return
		}
	}

	id, err := s.store.PostNew(&store.UserPostNew{
		UserId: user.Id,
		Title:  postTitle,
		URL:    postURL,
		Body:   postBody,
	})

	if err != nil {
		log.Printf("post internal error: %q", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", id)
}

func (s *Server) GetPostScoreHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	score, err := s.store.PostGetScore(postId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", score)
}

func (s *Server) GetPostUpvotesHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	score, err := s.store.PostGetUpvotes(postId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", score)
}

func (s *Server) GetPostDownvotesHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	score, err := s.store.PostGetDownvotes(postId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", score)
}

func (s *Server) CreatePostUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user := s.getUser(r)
	if user == nil {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}
	upvotes, err := s.store.PostCreateUpvote(user.Id, postId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", upvotes)
}

func (s *Server) CreatePostDownvoteHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user := s.getUser(r)
	if user == nil {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}
	downvotes, err := s.store.PostCreateDownvote(user.Id, postId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", downvotes)
}
