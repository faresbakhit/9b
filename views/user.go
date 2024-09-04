package views

import (
	"net/http"
	"time"
)

type UserData struct {
	Username  string
	CreatedAt *time.Time
	IsSelf    bool
	Err       string
}

func User(w http.ResponseWriter, data UserData, statusCode int) {
	renderPage(w, "user.html", data, statusCode)
}
