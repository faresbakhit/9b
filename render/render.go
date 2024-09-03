package render

import (
	"html/template"
	"net/http"
	"time"

	"github.com/faresbakhit/9b/config"
)

var t *template.Template

func init() {
	t = template.Must(template.ParseGlob(config.PAGES_GLOB_PATTERN))
}

func NotFound(w http.ResponseWriter, message string) {
	renderPage(w, "404.html", message, http.StatusNotFound)
}

func Unauthorized(w http.ResponseWriter) {
	renderPage(w, "401.html", nil, http.StatusUnauthorized)
}

type SignupPageData struct {
	Err  string
	Goto string
}

func SignupPage(w http.ResponseWriter, data SignupPageData, statusCode int) {
	renderPage(w, "signup.html", data, statusCode)
}

type LoginPageData struct {
	Err  string
	Goto string
}

func LoginPage(w http.ResponseWriter, data LoginPageData, statusCode int) {
	renderPage(w, "login.html", data, statusCode)
}

type UserPageData struct {
	Name     string
	JoinedAt *time.Time
	IsSelf   bool
	Err      string
}

func UserPage(w http.ResponseWriter, data UserPageData, statusCode int) {
	renderPage(w, "user.html", data, statusCode)
}

func renderPage(w http.ResponseWriter, name string, data any, statuCode int) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statuCode)
	return t.ExecuteTemplate(w, name, data)
}
