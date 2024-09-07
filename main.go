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
	defer s.Close()

	mux := http.NewServeMux()

	mux.HandleFunc(s.HomeGETPattern(), s.HomeGET)
	mux.HandleFunc(s.SignupPOSTPattern(), s.SignupPOST)
	mux.HandleFunc(s.LoginPOSTPattern(), s.LoginPOSTHandler)
	mux.HandleFunc(s.LogoutPOSTPattern(), s.LogoutPOSTHandler)
	mux.HandleFunc(s.ChangePasswordPOSTPattern(), s.ChangePasswordPOSTHandler)
	mux.HandleFunc(s.UserGETPattern(), s.UserGETHandler)
	mux.HandleFunc(s.CreatePostGETPattern(), s.CreatePostGETHandler)
	mux.HandleFunc(s.CreatePostPOSTPattern(), s.CreatePostPOSTHandler)

	mux.Handle(config.HTTP_PUBLIC_ROUTE, http.StripPrefix(config.HTTP_PUBLIC_ROUTE,
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
