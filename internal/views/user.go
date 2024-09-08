package views

import (
	"net/http"

	"github.com/faresbakhit/9b/internal/store"
)

type UserData struct {
	User      *store.User
	UserOther *store.User
	Posts     []*store.UserPost
	IsSelf    bool
	Err       string
}

func User(w http.ResponseWriter, data UserData, statusCode int) {
	renderPage(w, "user", data, statusCode)
}
