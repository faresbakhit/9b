package views

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/faresbakhit/9b/internal/config"
)

var t *template.Template
var funcMap = template.FuncMap{
	"split": func(sep string, s string) []string {
		return strings.Split(s, sep)
	},
}

func init() {
	t = template.Must(template.New("render").Funcs(funcMap).ParseGlob(config.TEMPLATES_GLOB_PATTERN))
}

func renderPage(w http.ResponseWriter, name string, data any, statuCode int) error {
	if config.TEMPLATES_LOAD_ON_RENDER {
		var err error
		t, err = template.New("render").Funcs(funcMap).ParseGlob(config.TEMPLATES_GLOB_PATTERN)
		if err != nil {
			log.Printf("error while reloading templates: %q", err)
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statuCode)
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("renderPage: %q", err)
	}
	return err
}
