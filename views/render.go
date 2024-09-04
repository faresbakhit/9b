package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/faresbakhit/9b/config"
)

var t *template.Template

func init() {
	t = template.Must(template.ParseGlob(config.PAGES_GLOB_PATTERN))
}

func renderPage(w http.ResponseWriter, name string, data any, statuCode int) error {
	if config.PAGES_RELOAD {
		var err error
		t, err = template.ParseGlob(config.PAGES_GLOB_PATTERN)
		if err != nil {
			log.Printf("error while reloading templates: %q", err)
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statuCode)
	return t.ExecuteTemplate(w, name, data)
}
