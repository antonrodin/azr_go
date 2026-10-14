package handlers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/antonrodin/azr/internal/authsvc"
	"github.com/antonrodin/azr/internal/blog"
	"github.com/antonrodin/azr/internal/i18n"
	"github.com/antonrodin/azr/internal/models/mysqlite"
	"github.com/antonrodin/azr/internal/openai"
	"github.com/antonrodin/azr/internal/photos"
	"github.com/go-chi/chi/v5"
)

var App *AppRepository

type AppRepository struct {
	Session *scs.SessionManager
	Pages   *mysqlite.PageModel
	Photos  *mysqlite.PhotoModel
	Prompts *mysqlite.PromptModel
	Storage *photos.Storage
	Auth    *authsvc.Client
	OpenAI  *openai.Client
	Blog    *blog.Store
}

func NewRepo(a *AppRepository) {
	App = a
}

// brandView alimenta la barra lateral de branding (partial "components/brand"),
// compartida entre la portada y el fotolog. Embebe i18n.Content para exponer
// los claims y enlaces; añade los destinos del conmutador de idioma (que
// dependen de la página) y el botón opcional "Subir foto".
type brandView struct {
	i18n.Content
	LangES     string
	LangEN     string
	LangRU     string
	ShowUpload bool
	UploadURL  string
}

// makeBrand construye la barra lateral. es/en/ru son los destinos del
// conmutador de idioma para la página actual.
func makeBrand(c i18n.Content, es, en, ru string, showUpload bool) brandView {
	return brandView{
		Content:    c,
		LangES:     es,
		LangEN:     en,
		LangRU:     ru,
		ShowUpload: showUpload,
		UploadURL:  "/admin/fotos/subir",
	}
}

// homePage es el dato de la portada: el contenido i18n (promovido para las
// plantillas) más la barra lateral.
type homePage struct {
	i18n.Content
	Brand brandView
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

	data := homePage{
		Content: c,
		Brand:   makeBrand(c, "/", "/en", "/ru", false),
	}

	app.render(w, "home.go.tmpl", data)
}

func (app *AppRepository) Contacto(w http.ResponseWriter, r *http.Request) {
	app.render(w, "contacto/contacto.go.tmpl", nil)
}

func (app *AppRepository) Cookies(w http.ResponseWriter, r *http.Request) {
	app.render(w, "contacto/cookies.go.tmpl", nil)
}

func (app *AppRepository) Legal(w http.ResponseWriter, r *http.Request) {
	app.render(w, "contacto/legal.go.tmpl", nil)
}
