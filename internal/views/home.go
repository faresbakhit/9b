package views

import (
	"net/http"

	"github.com/faresbakhit/9b/internal/store"
)

type HomeData struct {
	LoggedIn bool
	Posts    []*store.UserPostListResult
}

func Home(w http.ResponseWriter, data *HomeData, statusCode int) {
	renderPage(w, "home", data, statusCode)
}
