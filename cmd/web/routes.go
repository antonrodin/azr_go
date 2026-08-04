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
	router.Get("/", handlers.App.Home)       // español por defecto
	router.Get("/{lang}", handlers.App.Home) // /es, /en, /ru

	// Login contra el servicio auth
	router.Get("/login", handlers.App.LoginForm)
	router.Post("/login", handlers.App.LoginPost)

	// Zona privada del fotolog
	router.Route("/admin", func(r chi.Router) {
		r.Use(Auth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/admin/fotos", http.StatusSeeOther)
		})
		r.Get("/logout", handlers.App.Logout)

		r.Get("/fotos", handlers.App.AdminPhotos)
		r.Get("/fotos/subir", handlers.App.UploadForm)
		r.Post("/fotos", handlers.App.UploadPost)
		r.Get("/fotos/{id}", handlers.App.PhotoEditForm)
		r.Post("/fotos/{id}", handlers.App.PhotoEditPost)
		r.Post("/fotos/{id}/publicar", handlers.App.PhotoPublish)
		r.Post("/fotos/{id}/despublicar", handlers.App.PhotoUnpublish)
		r.Post("/fotos/{id}/borrar", handlers.App.PhotoDelete)
	})

	// Static files
	fileServer := http.FileServer(http.Dir("./public/"))
	router.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return router
}
