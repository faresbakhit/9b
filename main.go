package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	// UI
	mux.HandleFunc("GET /", s.HomeHandler)
	mux.HandleFunc("GET /posts/{id}", s.GetPostHandler)
	// Auth
	mux.HandleFunc("POST /signup", s.SignupHandler)
	mux.HandleFunc("POST /login", s.LoginHandler)
	mux.HandleFunc("POST /logout", s.LogoutHandler)
	mux.HandleFunc("POST /change-password", s.ChangePasswordHandler)
	// Posts
	mux.HandleFunc("POST /posts", s.CreatePostHandler)
	mux.HandleFunc("GET /posts/{id}/score", s.GetPostScoreHandler)
	mux.HandleFunc("GET /posts/{id}/upvotes", s.GetPostUpvotesHandler)
	mux.HandleFunc("GET /posts/{id}/downvotes", s.GetPostDownvotesHandler)
	mux.HandleFunc("POST /posts/{id}/upvotes", s.CreatePostUpvoteHandler)
	mux.HandleFunc("POST /posts/{id}/downvotes", s.CreatePostDownvoteHandler)

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
		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe: %v", err)
			}
		}()
	} else {
		go func() {
			if err := srv.ListenAndServeTLS(config.TLS_CERTIFICATE_FILE, config.TLS_PRIVATE_KEY_FILE); err != http.ErrServerClosed {
				log.Fatalf("ListenAndServeTLS: %v", err)
			}
		}()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	sig := <-signalChan
	log.Printf("%v: shutting down", sig)

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Print("shutdown error: ", err)
		defer os.Exit(1)
		return
	}

	defer os.Exit(0)
}
