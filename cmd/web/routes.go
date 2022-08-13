package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joefazee/hnews/public"
	"net/http"
)

func (a *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(a.CSRFTokenRequired)
	mux.Use(a.LoadSession)

	if a.debug {
		mux.Use(middleware.Logger)
	}

	// register routes
	mux.Get("/", a.homeHandler)
	mux.Get("/comments/{postId}", a.commentHandler)
	mux.Post("/comments/{postId}", a.commentPostHandler)

	mux.Get("/login", a.loginHandler)
	mux.Post("/login", a.loginPostHandler)
	mux.Get("/signup", a.signUpHandler)
	mux.Post("/signup", a.signPostUpHandler)
	mux.Get("/logout", a.authRequired(a.logoutHandler))

	mux.Get("/vote", a.authRequired(a.voteHandler))
	mux.Get("/submit", a.authRequired(a.submitHandler))
	mux.Post("/submit", a.authRequired(a.submitPostHandler))

	fileServer := http.FileServer(http.FS(public.Files))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
