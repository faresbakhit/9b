package views

import (
	"net/http"

	"github.com/faresbakhit/9b/store"
)

type HomeData struct {
	User *store.User
}

func Home(w http.ResponseWriter, data *HomeData, statusCode int) {
	renderPage(w, "home", data, statusCode)
}
