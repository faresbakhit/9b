package views

import "net/http"

func NotFound(w http.ResponseWriter, message string) {
	renderPage(w, "404.html", message, http.StatusNotFound)
}

func Unauthorized(w http.ResponseWriter) {
	renderPage(w, "401.html", nil, http.StatusUnauthorized)
}
