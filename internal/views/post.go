package views

import (
	"net/http"

	"github.com/faresbakhit/9b/internal/store"
)

type PostData struct {
	LoggedIn bool
	Post     *store.PostGet
}

func Post(w http.ResponseWriter, data *PostData, statusCode int) {
	renderPage(w, "post", data, statusCode)
}
