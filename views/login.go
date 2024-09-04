package views

import "net/http"

type LoginData struct {
	Err  string
	Goto string
}

func Login(w http.ResponseWriter, data LoginData, statusCode int) {
	renderPage(w, "login.html", data, statusCode)
}
