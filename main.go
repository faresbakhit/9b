package main

import (
	"log"
	"net/http"
	"time"

	"github.com/faresbakhit/9b/internal/config"
	"github.com/faresbakhit/9b/internal/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.HomeHandler)
	mux.HandleFunc("POST /signup", s.SignupHandler)
	mux.HandleFunc("POST /login", s.LoginHandler)
	mux.HandleFunc("POST /logout", s.LogoutHandler)
	mux.HandleFunc("POST /change-password", s.ChangePasswordHandler)
	mux.HandleFunc("POST /posts", s.CreatePostHandler)

	mux.Handle("GET "+config.HTTP_PUBLIC_ROUTE, http.StripPrefix(config.HTTP_PUBLIC_ROUTE,
		http.FileServer(http.Dir(config.HTTP_PUBLIC_DIRECTORY))))

	srv := &http.Server{
		Addr:         config.HTTP_SERVER_ADDR,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("server listening on %s", srv.Addr)

	if config.TLS_CERTIFICATE_FILE == "" && config.TLS_PRIVATE_KEY_FILE == "" {
		log.Fatal(srv.ListenAndServe())
	}

	log.Fatal(srv.ListenAndServeTLS(config.TLS_CERTIFICATE_FILE, config.TLS_PRIVATE_KEY_FILE))
}
