package views

import "net/http"

type SignupData struct {
	Err  string
	Goto string
}

func Signup(w http.ResponseWriter, data SignupData, statusCode int) {
	renderPage(w, "signup", data, statusCode)
}
