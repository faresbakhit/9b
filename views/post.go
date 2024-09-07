package views

import (
	"net/http"

	"github.com/faresbakhit/9b/store"
)

type CreatePostData struct {
	User *store.User
	Err  string
}

func CreatePost(w http.ResponseWriter, data *CreatePostData, statusCode int) {
	renderPage(w, "create_post", data, statusCode)
}
