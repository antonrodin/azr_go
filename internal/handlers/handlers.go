package handlers

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/antonrodin/azr/internal/i18n"
	"github.com/antonrodin/azr/internal/models/mysqlite"
	"github.com/antonrodin/azr/internal/tmpl"
	"github.com/go-chi/chi/v5"
)

var App *AppRepository

type AppRepository struct {
	Session *scs.SessionManager
	Pages   *mysqlite.PageModel
}

func NewRepo(a *AppRepository) {
	App = a
}

func (app *AppRepository) Home(w http.ResponseWriter, r *http.Request) {
	// El idioma viene del segmento de ruta ({lang}); si no hay, usamos el
	// idioma por defecto (español, ruta "/").
	lang := chi.URLParam(r, "lang")
	if lang == "" {
		lang = i18n.Default
	}

	c, ok := i18n.Get(lang)
	if !ok {
		http.NotFound(w, r)
		return
	}

	err := tmpl.Get("home.go.tmpl").Execute(w, c)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (app *AppRepository) Contacto(w http.ResponseWriter, r *http.Request) {

	err := tmpl.Get("contacto/contacto.go.tmpl").Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (app *AppRepository) Cookies(w http.ResponseWriter, r *http.Request) {

	err := tmpl.Get("contacto/cookies.go.tmpl").Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (app *AppRepository) Legal(w http.ResponseWriter, r *http.Request) {

	err := tmpl.Get("contacto/legal.go.tmpl").Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
