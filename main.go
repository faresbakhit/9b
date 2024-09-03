package main

import (
	"log"
	"net/http"
	"time"

	"github.com/faresbakhit/9b/config"
	"github.com/faresbakhit/9b/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /signup", s.SignupPage)
	mux.HandleFunc("POST /signup", s.Signup)
	mux.HandleFunc("GET /login", s.LoginPage)
	mux.HandleFunc("POST /login", s.Login)
	mux.HandleFunc("GET /u/{name}", s.UserPage)
	mux.HandleFunc("POST /u/{name}", s.UpdatePassword)
	mux.HandleFunc("GET /protected", s.ProtectedPage)

	srv := &http.Server{
		Addr:         config.HTTP_SERVER_ADDR,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Listening on %s", srv.Addr)

	if config.TLS_CERTIFICATE_FILE == "" && config.TLS_PRIVATE_KEY_FILE == "" {
		log.Fatal(srv.ListenAndServe())
	}

	log.Fatal(srv.ListenAndServeTLS(config.TLS_CERTIFICATE_FILE, config.TLS_PRIVATE_KEY_FILE))
}
