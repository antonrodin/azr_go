package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antonrodin/azr/internal/blog"
	"github.com/antonrodin/azr/internal/i18n"
	"github.com/antonrodin/azr/internal/tmpl"
	"github.com/go-chi/chi/v5"
)

// BlogDir es la carpeta de los artículos en Markdown.
const BlogDir = "./content/blog"

// blogStore devuelve los artículos. En dev (sin caché de plantillas) se
// recargan en cada petición para ver los cambios del Markdown al momento.
func (app *AppRepository) blogStore() *blog.Store {
	if !tmpl.Cache {
		if s, err := blog.Load(BlogDir); err == nil {
			return s
		} else {
			log.Println("blog:", err)
		}
	}
	return app.Blog
}

func blogIndexPath(lang string) string { return "/" + lang + "/blog" }

func articlePath(a *blog.Article) string { return blogIndexPath(a.Lang) + "/" + a.Slug }

// RedirectBlog manda las rutas sin idioma (/blog, /blog/{slug}) a la versión
// en español.
func (app *AppRepository) RedirectBlog(w http.ResponseWriter, r *http.Request) {
	target := blogIndexPath(i18n.Default)
	if slug := chi.URLParam(r, "slug"); slug != "" {
		target += "/" + slug
	}
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

// --- Listado ---

type blogIndexItem struct {
	Href        string
	Title       string
	Description string
	Date        string
	DateISO     string
}

type blogIndexPage struct {
	Lang      string
	Brand     brandView
	T         i18n.BlogText
	Canonical string
	AltES     string
	AltEN     string
	AltRU     string
	Articles  []blogIndexItem
}

func (app *AppRepository) BlogIndex(w http.ResponseWriter, r *http.Request) {
	lang := chi.URLParam(r, "lang")
	t, ok := i18n.GetBlog(lang)
	if !ok {
		http.NotFound(w, r)
		return
	}
	c, _ := i18n.Get(lang)

	data := blogIndexPage{
		Lang:      lang,
		Brand:     makeBrand(c, blogIndexPath("es"), blogIndexPath("en"), blogIndexPath("ru"), false),
		T:         t,
		Canonical: siteBase + blogIndexPath(lang),
		AltES:     siteBase + blogIndexPath("es"),
		AltEN:     siteBase + blogIndexPath("en"),
		AltRU:     siteBase + blogIndexPath("ru"),
	}
	for _, a := range app.blogStore().List(lang) {
		data.Articles = append(data.Articles, blogIndexItem{
			Href:        articlePath(a),
			Title:       a.Title,
			Description: a.Description,
			Date:        t.FormatDate(a.Date),
			DateISO:     a.Date.Format("2006-01-02"),
		})
	}

	app.render(w, "blog/index.go.tmpl", data)
}

// --- Artículo ---

type articlePage struct {
	Lang      string
	Brand     brandView
	T         i18n.BlogText
	Article   *blog.Article
	Date      string
	DateISO   string
	Canonical string
	OgImage   string
	JSONLD    string
	// hreflang: solo los idiomas en los que existe el artículo
	Alternates []alternate
	XDefault   string
}

type alternate struct {
	Lang string
	Href string
}

func (app *AppRepository) BlogArticle(w http.ResponseWriter, r *http.Request) {
	lang := chi.URLParam(r, "lang")
	t, ok := i18n.GetBlog(lang)
	if !ok {
		http.NotFound(w, r)
		return
	}
	store := app.blogStore()
	a, ok := store.Get(lang, chi.URLParam(r, "slug"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	c, _ := i18n.Get(lang)

	// Conmutador de idioma: la traducción si existe, si no el listado.
	switchTo := map[string]string{}
	data := articlePage{
		Lang:      lang,
		T:         t,
		Article:   a,
		Date:      t.FormatDate(a.Date),
		DateISO:   a.Date.Format("2006-01-02"),
		Canonical: siteBase + articlePath(a),
		OgImage:   siteBase + a.Image,
	}
	for _, l := range []string{"es", "en", "ru"} {
		if tr, ok := store.Translation(a, l); ok {
			switchTo[l] = articlePath(tr)
			data.Alternates = append(data.Alternates, alternate{Lang: l, Href: siteBase + articlePath(tr)})
			if l == i18n.Default {
				data.XDefault = siteBase + articlePath(tr)
			}
		} else {
			switchTo[l] = blogIndexPath(l)
		}
	}
	data.Brand = makeBrand(c, switchTo["es"], switchTo["en"], switchTo["ru"], false)
	data.JSONLD = articleJSONLD(a, data.Canonical, data.OgImage)

	app.render(w, "blog/show.go.tmpl", data)
}

func articleJSONLD(a *blog.Article, canonical, image string) string {
	author := map[string]interface{}{
		"@type": "Person",
		"@id":   siteBase + "/#person",
		"name":  i18n.Author,
		"url":   siteBase + "/",
	}
	ld := map[string]interface{}{
		"@context":         "https://schema.org",
		"@type":            "BlogPosting",
		"headline":         a.Title,
		"description":      a.Description,
		"inLanguage":       a.Lang,
		"datePublished":    a.Date.Format("2006-01-02"),
		"image":            image,
		"url":              canonical,
		"mainEntityOfPage": canonical,
		"author":           author,
		"publisher":        author,
	}
	// json.Marshal escapa < > & por defecto, así que es seguro dentro de <script>.
	b, err := json.Marshal(ld)
	if err != nil {
		return "{}"
	}
	return string(b)
}
