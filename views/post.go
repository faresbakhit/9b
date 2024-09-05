package views

import "net/http"

type CreatePostData struct {
	Err string
}

func CreatePost(w http.ResponseWriter, data CreatePostData, statusCode int) {
	renderPage(w, "create_post", data, statusCode)
}
