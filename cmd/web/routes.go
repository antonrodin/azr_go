package main

import (
	"net/http"

	"github.com/antonrodin/azr/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(LoadSession)

	// Rutes
	router.Get("/", handlers.App.Home)        // español por defecto
	router.Get("/{lang}", handlers.App.Home)  // /es, /en, /ru

	// Static files
	fileServer := http.FileServer(http.Dir("./public/"))
	router.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return router
}
