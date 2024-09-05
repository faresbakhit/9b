package views

import (
	"net/http"
	"time"

	"github.com/faresbakhit/9b/store"
)

type UserData struct {
	Username  string
	Posts     []*store.UserPost
	CreatedAt *time.Time
	IsSelf    bool
	Err       string
}

func User(w http.ResponseWriter, data UserData, statusCode int) {
	renderPage(w, "user", data, statusCode)
}
