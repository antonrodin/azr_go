package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/antonrodin/azr/internal/i18n"
	"github.com/antonrodin/azr/internal/models"
	"github.com/antonrodin/azr/internal/photos"
	"github.com/go-chi/chi/v5"
)

const (
	siteBase     = "https://azr.es"
	perPage      = 20
	fotologImage = "2048.jpg" // imagen de referencia para OG y sitemap
)

// --- Helpers de rutas del fotolog ---

func fotologBase(lang string) string {
	switch lang {
	case "en":
		return "/en/fotolog"
	case "ru":
		return "/ru/fotolog"
	default:
		return "/fotolog"
	}
}

// galleryPath construye la ruta de la galería, preservando la página (>1) y
// el filtro de búsqueda (q) si los hay.
func galleryPath(lang string, page int, q string) string {
	b := fotologBase(lang)
	if page > 1 {
		b = fmt.Sprintf("%s/pagina/%d", b, page)
	}
	if q != "" {
		b += "?q=" + url.QueryEscape(q)
	}
	return b
}

func photoPath(lang, slug string) string {
	return fotologBase(lang) + "/" + slug
}

// resolveLang saca el idioma del segmento {lang} (vacío = español) y valida
// que esté soportado.
func resolveLang(r *http.Request) (string, bool) {
	lang := chi.URLParam(r, "lang")
	if lang == "" {
		lang = i18n.Default
	}
	_, ok := i18n.GetFotolog(lang)
	return lang, ok
}

// localTitle/localDesc devuelven el texto en el idioma pedido con respaldo
// en español (nunca dejamos una foto sin título si el ES existe).
func localTitle(p *models.Photo, lang string) string {
	switch lang {
	case "en":
		if p.TitleEN != "" {
			return p.TitleEN
		}
	case "ru":
		if p.TitleRU != "" {
			return p.TitleRU
		}
	}
	return p.TitleES
}

func localDesc(p *models.Photo, lang string) string {
	switch lang {
	case "en":
		if p.DescEN != "" {
			return p.DescEN
		}
	case "ru":
		if p.DescRU != "" {
			return p.DescRU
		}
	}
	return p.DescES
}

func (app *AppRepository) pubURL(id int, name string) string {
	return app.Storage.PublicURL(photos.Key("public", id, name))
}

// --- Galería ---

type galleryPhoto struct {
	Href   string // ruta al detalle en el idioma actual
	Title  string
	Src    string // 400.webp
	SrcSet string
	Width  int
	Height int
}

type galleryPage struct {
	Lang      string // para <html lang> del layout
	Brand     brandView
	T         i18n.FotologText
	Canonical string
	NoIndex   bool // true en páginas de resultados de búsqueda
	Photos    []galleryPhoto
	// Panel de control
	Query        string // término de búsqueda actual
	Searching    bool   // hay filtro activo
	SearchAction string // action del formulario (base del idioma)
	UploadURL    string
	Page         int
	LastPage     int
	PrevURL      string
	NextURL      string
	// Conmutador de idioma (misma página/búsqueda en cada idioma)
	LangES string
	LangEN string
	LangRU string
	// hreflang de esta galería
	AltES string
	AltEN string
	AltRU string
}

func (app *AppRepository) Fotolog(w http.ResponseWriter, r *http.Request) {
	lang, ok := resolveLang(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	t, _ := i18n.GetFotolog(lang)

	page := 1
	if s := chi.URLParam(r, "n"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			http.NotFound(w, r)
			return
		}
		page = n
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))

	total, err := app.Photos.CountPublic(q)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	lastPage := (total + perPage - 1) / perPage
	if lastPage == 0 {
		lastPage = 1
	}
	if page > lastPage {
		http.NotFound(w, r)
		return
	}

	list, err := app.Photos.PublicPage(q, perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Barra lateral compartida con la portada (sin botón de subir: ese vive
	// ahora en el panel de control).
	c, _ := i18n.Get(lang)

	data := galleryPage{
		Lang:         lang,
		Brand:        makeBrand(c, galleryPath("es", page, q), galleryPath("en", page, q), galleryPath("ru", page, q), false),
		T:            t,
		Canonical:    siteBase + galleryPath(lang, page, q),
		NoIndex:      q != "",
		Query:        q,
		Searching:    q != "",
		SearchAction: fotologBase(lang),
		UploadURL:    "/admin/fotos/subir",
		Page:         page,
		LastPage:     lastPage,
		LangES:       galleryPath("es", page, q),
		LangEN:       galleryPath("en", page, q),
		LangRU:       galleryPath("ru", page, q),
		AltES:        siteBase + galleryPath("es", page, q),
		AltEN:        siteBase + galleryPath("en", page, q),
		AltRU:        siteBase + galleryPath("ru", page, q),
	}
	for _, p := range list {
		data.Photos = append(data.Photos, galleryPhoto{
			Href:   photoPath(lang, p.Slug),
			Title:  localTitle(p, lang),
			Src:    app.pubURL(p.ID, "1200.webp"),
			SrcSet: app.pubURL(p.ID, "400.webp") + " 400w, " + app.pubURL(p.ID, "1200.webp") + " 1200w, " + app.pubURL(p.ID, "2048.webp") + " 2048w",
			Width:  p.Width,
			Height: p.Height,
		})
	}
	if page > 1 {
		data.PrevURL = galleryPath(lang, page-1, q)
	}
	if page < lastPage {
		data.NextURL = galleryPath(lang, page+1, q)
	}

	app.render(w, "fotolog/index.go.tmpl", data)
}

// --- Detalle ---

type photoDownload struct {
	Label string
	URL   string
	Note  string
}

type photoPage struct {
	Lang        string // para <html lang> del layout
	T           i18n.FotologText
	Slug        string
	Title       string
	Desc        string
	Canonical   string
	Photo       *models.Photo
	SrcJPG      string // fallback <img> src (1200.jpg)
	SrcSetWebP  string
	SrcSetJPG   string
	OgImage     string
	Downloads   []photoDownload
	ShowGPS     bool
	GPSLat      float64
	GPSLng      float64
	JSONLD      string
	LicenseURL  string
	LicenseName string
	// hreflang y conmutador de idioma
	AltES string
	AltEN string
	AltRU string
}

func (app *AppRepository) FotologPhoto(w http.ResponseWriter, r *http.Request) {
	lang, ok := resolveLang(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	t, _ := i18n.GetFotolog(lang)

	slug := chi.URLParam(r, "slug")
	p, err := app.Photos.GetPublicBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	title := localTitle(p, lang)
	if title == "" {
		title = t.NoTitle
	}
	desc := localDesc(p, lang)
	canonical := siteBase + photoPath(lang, p.Slug)

	data := photoPage{
		Lang:        lang,
		T:           t,
		Slug:        p.Slug,
		Title:       title,
		Desc:        desc,
		Canonical:   canonical,
		Photo:       p,
		SrcJPG:      app.pubURL(p.ID, "1200.jpg"),
		SrcSetWebP:  app.pubURL(p.ID, "400.webp") + " 400w, " + app.pubURL(p.ID, "1200.webp") + " 1200w, " + app.pubURL(p.ID, "2048.webp") + " 2048w",
		SrcSetJPG:   app.pubURL(p.ID, "1200.jpg") + " 1200w, " + app.pubURL(p.ID, "2048.jpg") + " 2048w",
		OgImage:     app.pubURL(p.ID, fotologImage),
		ShowGPS:     p.HasGPS && p.ShowGPS,
		GPSLat:      p.GPSLat,
		GPSLng:      p.GPSLng,
		LicenseURL:  i18n.LicenseURL,
		LicenseName: i18n.LicenseName,
		AltES:       siteBase + photoPath("es", p.Slug),
		AltEN:       siteBase + photoPath("en", p.Slug),
		AltRU:       siteBase + photoPath("ru", p.Slug),
		Downloads: []photoDownload{
			{Label: t.DownloadMedium, URL: app.pubURL(p.ID, "1200.jpg"), Note: "1200 px"},
			{Label: t.DownloadLarge, URL: app.pubURL(p.ID, "2048.jpg"), Note: "2048 px"},
			{Label: t.DownloadOriginal, URL: app.pubURL(p.ID, "orig."+p.OrigExt), Note: humanBytes(p.OrigBytes)},
		},
	}

	data.JSONLD = app.buildJSONLD(p, title, desc, canonical)

	app.render(w, "fotolog/show.go.tmpl", data)
}

// buildJSONLD arma el ImageObject de schema.org. La pareja license +
// acquireLicensePage es lo que dispara el badge "Licenciable" en Google Images.
func (app *AppRepository) buildJSONLD(p *models.Photo, title, desc, canonical string) string {
	ld := map[string]interface{}{
		"@context":   "https://schema.org",
		"@type":      "ImageObject",
		"contentUrl": app.pubURL(p.ID, "2048.jpg"),
		"thumbnailUrl": app.pubURL(p.ID, "400.webp"),
		"name":       title,
		"creator": map[string]interface{}{
			"@type": "Person",
			"name":  i18n.Author,
			"url":   siteBase,
		},
		"copyrightNotice":    i18n.Author,
		"creditText":         i18n.AttributionText,
		"license":            i18n.LicenseURL,
		"acquireLicensePage": canonical,
		"width":              p.Width,
		"height":             p.Height,
	}
	if desc != "" {
		ld["description"] = desc
	}
	if p.TakenAt != "" {
		ld["dateCreated"] = p.TakenAt
	}
	// json.Marshal escapa < > & por defecto, así que es seguro dentro de <script>.
	b, err := json.Marshal(ld)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}

// --- Sitemap de imágenes y robots ---

func (app *AppRepository) Sitemap(w http.ResponseWriter, r *http.Request) {
	pics, err := app.Photos.AllPublic()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">` + "\n")

	// Home y galería en los tres idiomas.
	for _, path := range []string{"/", "/en", "/ru", "/fotolog", "/en/fotolog", "/ru/fotolog"} {
		b.WriteString("  <url><loc>" + xmlEscape(siteBase+path) + "</loc></url>\n")
	}

	// Blog: listado y artículos en cada idioma.
	posts := app.blogStore()
	for _, lang := range []string{"es", "en", "ru"} {
		b.WriteString("  <url><loc>" + xmlEscape(siteBase+blogIndexPath(lang)) + "</loc></url>\n")
		for _, a := range posts.List(lang) {
			b.WriteString("  <url><loc>" + xmlEscape(siteBase+articlePath(a)) + "</loc>")
			b.WriteString("<lastmod>" + a.Date.Format("2006-01-02") + "</lastmod>")
			if a.Image != "" {
				b.WriteString("<image:image><image:loc>" + xmlEscape(siteBase+a.Image) + "</image:loc></image:image>")
			}
			b.WriteString("</url>\n")
		}
	}

	// Cada foto pública, en los tres idiomas, con su imagen.
	for _, p := range pics {
		imgLoc := app.pubURL(p.ID, fotologImage)
		for _, lang := range []string{"es", "en", "ru"} {
			loc := siteBase + photoPath(lang, p.Slug)
			title := localTitle(p, lang)
			b.WriteString("  <url><loc>" + xmlEscape(loc) + "</loc>")
			b.WriteString("<image:image><image:loc>" + xmlEscape(imgLoc) + "</image:loc>")
			if title != "" {
				b.WriteString("<image:title>" + xmlEscape(title) + "</image:title>")
			}
			b.WriteString("<image:license>" + xmlEscape(i18n.LicenseURL) + "</image:license>")
			b.WriteString("</image:image></url>\n")
		}
	}

	b.WriteString("</urlset>\n")

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(b.String()))
}

func (app *AppRepository) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /login\n\nSitemap: %s/sitemap.xml\n", siteBase)
}

func xmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return r.Replace(s)
}
