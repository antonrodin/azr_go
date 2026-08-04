package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antonrodin/azr/internal/authsvc"
	"github.com/antonrodin/azr/internal/models"
	"github.com/antonrodin/azr/internal/photos"
	"github.com/antonrodin/azr/internal/tmpl"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
)

// Límite por petición de subida. Fotos de móvil rondan 3-15 MB; esto deja
// margen de sobra sin permitir barbaridades.
const maxUploadBytes = 100 << 20 // 100 MB

var allowedExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true,
	"heic": true, "heif": true, "webp": true,
}

const presignTTL = 15 * time.Minute

// adminPage es el dato común de las vistas de admin (siempre en español).
type adminPage struct {
	Lang  string
	Flash string
	Error string
}

// photoView es una foto lista para pintar: la foto más las URLs resueltas
// (presignadas si es privada, públicas si no).
type photoView struct {
	*models.Photo
	ThumbURL string
	ImageURL string
	OrigURL  string
}

func (app *AppRepository) photoURL(r *http.Request, p *models.Photo, name string) string {
	key := photos.Key(p.Status, p.ID, name)
	if p.Status == "public" {
		return app.Storage.PublicURL(key)
	}
	u, err := app.Storage.Presign(r.Context(), key, presignTTL)
	if err != nil {
		log.Println("presign:", err)
		return ""
	}
	return u
}

func (app *AppRepository) render(w http.ResponseWriter, name string, data interface{}) {
	if err := tmpl.Get(name).Execute(w, data); err != nil {
		log.Println("render", name+":", err)
		http.Error(w, err.Error(), 500)
	}
}

// --- Login / Logout ---

type loginPage struct {
	adminPage
	Email string
	Next  string
}

func (app *AppRepository) LoginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, "admin/login.go.tmpl", loginPage{
		adminPage: adminPage{Lang: "es"},
		Next:      safeNext(r.URL.Query().Get("next")),
	})
}

func (app *AppRepository) LoginPost(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := r.PostFormValue("password")
	next := safeNext(r.PostFormValue("next"))

	token, err := app.Auth.Login(r.Context(), email, password)
	if err != nil {
		msg := "Eso no era. Prueba otra vez."
		if !errors.Is(err, authsvc.ErrInvalidCredentials) {
			log.Println("login:", err)
			msg = "El servicio de auth no responde. Culpa suya, no tuya."
		}
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, "admin/login.go.tmpl", loginPage{
			adminPage: adminPage{Lang: "es", Error: msg},
			Email:     email,
			Next:      next,
		})
		return
	}

	// Renovar el id de sesión al autenticar (session fixation).
	if err := app.Session.RenewToken(r.Context()); err != nil {
		log.Println("renew session:", err)
	}
	app.Session.Put(r.Context(), "token", token)

	http.Redirect(w, r, next, http.StatusSeeOther)
}

// safeNext valida el destino post-login: solo rutas internas de este sitio.
// Rechaza URLs absolutas y protocol-relative ("//host") para evitar
// redirecciones abiertas a sitios externos. Por defecto, el listado de fotos.
func safeNext(next string) string {
	if strings.HasPrefix(next, "/") && !strings.HasPrefix(next, "//") {
		return next
	}
	return "/admin/fotos"
}

func (app *AppRepository) Logout(w http.ResponseWriter, r *http.Request) {
	if err := app.Session.Destroy(r.Context()); err != nil {
		log.Println("logout:", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Listado ---

type adminPhotosPage struct {
	adminPage
	Photos []photoView
}

func (app *AppRepository) AdminPhotos(w http.ResponseWriter, r *http.Request) {
	all, err := app.Photos.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	page := adminPhotosPage{
		adminPage: adminPage{Lang: "es", Flash: app.Session.PopString(r.Context(), "flash")},
	}
	for _, p := range all {
		page.Photos = append(page.Photos, photoView{
			Photo:    p,
			ThumbURL: app.photoURL(r, p, "400.webp"),
		})
	}

	app.render(w, "admin/fotos.go.tmpl", page)
}

// --- Subida ---

func (app *AppRepository) UploadForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, "admin/subir.go.tmpl", adminPage{Lang: "es"})
}

// wantsJSON detecta la subida por fetch (JS) frente al form clásico.
func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func uploadError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": msg})
		return
	}
	http.Error(w, msg, status)
}

func (app *AppRepository) UploadPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		uploadError(w, r, http.StatusRequestEntityTooLarge, "La foto es demasiado grande (máx. 100 MB).")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		uploadError(w, r, http.StatusBadRequest, "No ha llegado ninguna foto.")
		return
	}
	defer file.Close()

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(header.Filename), "."))
	if ext == "jpeg" {
		ext = "jpg"
	}
	if !allowedExts[ext] {
		uploadError(w, r, http.StatusBadRequest, "Formato no soportado: "+ext)
		return
	}

	id, err := app.processUpload(r, file, ext)
	if err != nil {
		log.Println("subida:", err)
		uploadError(w, r, http.StatusInternalServerError, "Algo se ha roto procesando la foto: "+err.Error())
		return
	}

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id, "name": header.Filename})
		return
	}
	app.Session.Put(r.Context(), "flash", "Foto subida. Ponle título, anda.")
	http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", id), http.StatusSeeOther)
}

// processUpload ejecuta el pipeline completo: temporal → exiftool → vips →
// fila en BD → S3. Si S3 falla a medias, deshace lo que pueda.
func (app *AppRepository) processUpload(r *http.Request, file io.Reader, ext string) (int, error) {
	dir, err := os.MkdirTemp("", "fotolog-*")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(dir)

	origName := "orig." + ext
	origPath := filepath.Join(dir, origName)
	dst, err := os.Create(origPath)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(dst, file)
	dst.Close()
	if err != nil {
		return 0, err
	}

	meta, err := photos.ExtractMeta(origPath)
	if err != nil {
		return 0, err
	}

	derivs, err := photos.BuildDerivatives(origPath, dir)
	if err != nil {
		return 0, err
	}

	p := &models.Photo{
		// Slug provisional único; el definitivo se pone con el id.
		Slug:        fmt.Sprintf("tmp-%d", time.Now().UnixNano()),
		TakenAt:     meta.TakenAt,
		CameraMake:  meta.CameraMake,
		CameraModel: meta.CameraModel,
		Lens:        meta.Lens,
		Focal:       meta.Focal,
		Aperture:    meta.Aperture,
		Shutter:     meta.Shutter,
		ISO:         meta.ISO,
		HasGPS:      meta.HasGPS,
		GPSLat:      meta.GPSLat,
		GPSLng:      meta.GPSLng,
		Width:       meta.Width,
		Height:      meta.Height,
		OrigExt:     ext,
		OrigBytes:   written,
		ExifJSON:    meta.ExifJSON,
	}

	id, err := app.Photos.Insert(p)
	if err != nil {
		return 0, err
	}
	if err := app.Photos.UpdateSlug(id, fmt.Sprintf("foto-%d", id)); err != nil {
		return 0, err
	}

	ctx := r.Context()
	uploads := append([]photos.Derivative{
		{Name: origName, Path: origPath, ContentType: photos.ContentTypeByExt(ext)},
	}, derivs...)

	for _, d := range uploads {
		if err := app.Storage.UploadFile(ctx, photos.Key("private", id, d.Name), d.Path, d.ContentType); err != nil {
			// Deshacer: fuera lo subido y fuera la fila.
			if cerr := app.Storage.DeletePhoto(ctx, id, "private"); cerr != nil {
				log.Println("limpieza S3 tras fallo:", cerr)
			}
			if derr := app.Photos.Delete(id); derr != nil {
				log.Println("limpieza BD tras fallo:", derr)
			}
			return 0, fmt.Errorf("subir %s a S3: %w", d.Name, err)
		}
	}

	return id, nil
}

// --- Edición ---

type photoEditPage struct {
	adminPage
	Photo photoView
}

func (app *AppRepository) photoFromURL(r *http.Request) (*models.Photo, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}
	return app.Photos.Get(id)
}

func (app *AppRepository) PhotoEditForm(w http.ResponseWriter, r *http.Request) {
	p, err := app.photoFromURL(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	app.render(w, "admin/foto.go.tmpl", photoEditPage{
		adminPage: adminPage{Lang: "es", Flash: app.Session.PopString(r.Context(), "flash")},
		Photo: photoView{
			Photo:    p,
			ImageURL: app.photoURL(r, p, "1200.webp"),
			OrigURL:  app.photoURL(r, p, "orig."+p.OrigExt),
		},
	})
}

func (app *AppRepository) PhotoEditPost(w http.ResponseWriter, r *http.Request) {
	p, err := app.photoFromURL(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	p.TitleES = strings.TrimSpace(r.PostFormValue("title_es"))
	p.TitleEN = strings.TrimSpace(r.PostFormValue("title_en"))
	p.TitleRU = strings.TrimSpace(r.PostFormValue("title_ru"))
	p.DescES = strings.TrimSpace(r.PostFormValue("desc_es"))
	p.DescEN = strings.TrimSpace(r.PostFormValue("desc_en"))
	p.DescRU = strings.TrimSpace(r.PostFormValue("desc_ru"))
	p.ShowGPS = r.PostFormValue("show_gps") == "1"

	// El slug sale del título en español y se congela al publicar
	// (las URLs públicas no se rompen).
	if p.Status == "private" && p.TitleES != "" {
		if base := slug.Make(p.TitleES); base != "" && base != p.Slug {
			taken, err := app.Photos.SlugTakenByOther(base, p.ID)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if taken {
				base = fmt.Sprintf("%s-%d", base, p.ID)
			}
			p.Slug = base
		}
	}

	// Autotraducción (Fase 3): rellena EN/RU a partir del español. Por
	// defecto solo completa lo que esté vacío; "retraducir" sobrescribe.
	flash := "Guardado."
	retranslate := r.PostFormValue("retranslate") == "1"
	if msg := app.autoTranslate(r, p, retranslate); msg != "" {
		flash = msg
	}

	if err := app.Photos.UpdateTexts(p.ID, p); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.Session.Put(r.Context(), "flash", flash)
	http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
}

// autoTranslate rellena las traducciones EN/RU de p a partir del español.
// Modifica p in situ y devuelve el mensaje de flash a mostrar (o "" para el
// mensaje por defecto). Nunca falla el guardado: si la API se cae, avisa y
// deja los textos como estaban.
func (app *AppRepository) autoTranslate(r *http.Request, p *models.Photo, overwrite bool) string {
	if !app.OpenAI.Enabled() {
		return ""
	}
	// Nada que traducir si no hay origen en español.
	if p.TitleES == "" && p.DescES == "" {
		return ""
	}

	// ¿Hace falta? Solo si vamos a sobrescribir o hay algún destino vacío
	// con su origen presente (así no gastamos una llamada en balde).
	need := overwrite ||
		(p.TitleES != "" && (p.TitleEN == "" || p.TitleRU == "")) ||
		(p.DescES != "" && (p.DescEN == "" || p.DescRU == ""))
	if !need {
		return ""
	}

	tr, err := app.OpenAI.Translate(r.Context(), p.TitleES, p.DescES)
	if err != nil {
		log.Println("autotraducción:", err)
		return "Guardado, pero la traducción automática falló (revísalo o reintenta)."
	}

	// Rellena cada campo si toca sobrescribir o si está vacío.
	fill := func(dst *string, val string) {
		if val == "" {
			return
		}
		if overwrite || *dst == "" {
			*dst = strings.TrimSpace(val)
		}
	}
	fill(&p.TitleEN, tr.TitleEN)
	fill(&p.TitleRU, tr.TitleRU)
	fill(&p.DescEN, tr.DescEN)
	fill(&p.DescRU, tr.DescRU)

	if overwrite {
		return "Guardado y retraducido a inglés y ruso."
	}
	return "Guardado y traducido a inglés y ruso."
}

// --- Dictado por voz (Fase 4) ---

// maxAudioBytes es el límite de OpenAI para transcripción.
const maxAudioBytes = 25 << 20 // 25 MB

var audioExts = map[string]bool{
	"webm": true, "mp4": true, "m4a": true, "mp3": true,
	"wav": true, "ogg": true, "oga": true, "mpeg": true, "mpga": true, "flac": true,
}

// Transcribe recibe un audio grabado en el navegador, lo transcribe a español
// y lo pule en plan SEO. Devuelve JSON {ok, text} para que el JS rellene el
// textarea como borrador (no guarda nada).
func (app *AppRepository) Transcribe(w http.ResponseWriter, r *http.Request) {
	if !app.OpenAI.Enabled() {
		jsonError(w, http.StatusServiceUnavailable, "Transcripción no disponible (sin OpenAI).")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAudioBytes)
	if err := r.ParseMultipartForm(4 << 20); err != nil {
		jsonError(w, http.StatusRequestEntityTooLarge, "El audio es demasiado grande (máx. 25 MB).")
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "No ha llegado audio.")
		return
	}
	defer file.Close()

	text, err := app.OpenAI.Transcribe(r.Context(), file, transcribeFilename(header.Filename, header.Header.Get("Content-Type")))
	if err != nil {
		log.Println("transcribir:", err)
		jsonError(w, http.StatusBadGateway, "No se pudo transcribir: "+err.Error())
		return
	}
	if text == "" {
		jsonError(w, http.StatusUnprocessableEntity, "No se entendió nada en el audio. Repite más claro.")
		return
	}

	// Pulir en plan SEO; si el pulido falla, devolvemos la transcripción cruda.
	polished, perr := app.OpenAI.Polish(r.Context(), text)
	if perr != nil || polished == "" {
		if perr != nil {
			log.Println("pulir:", perr)
		}
		polished = text
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "text": polished})
}

// transcribeFilename da a OpenAI un nombre con extensión válida para que
// detecte el formato del audio.
func transcribeFilename(orig, contentType string) string {
	if ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(orig), ".")); audioExts[ext] {
		return "audio." + ext
	}
	switch {
	case strings.Contains(contentType, "webm"):
		return "audio.webm"
	case strings.Contains(contentType, "mp4"), strings.Contains(contentType, "m4a"), strings.Contains(contentType, "aac"):
		return "audio.mp4"
	case strings.Contains(contentType, "ogg"):
		return "audio.ogg"
	case strings.Contains(contentType, "wav"):
		return "audio.wav"
	case strings.Contains(contentType, "mpeg"), strings.Contains(contentType, "mp3"):
		return "audio.mp3"
	default:
		return "audio.webm"
	}
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": msg})
}

// --- Publicar / Despublicar / Borrar ---

func (app *AppRepository) PhotoPublish(w http.ResponseWriter, r *http.Request) {
	p, err := app.photoFromURL(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if p.Status == "public" {
		http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
		return
	}
	if p.TitleES == "" {
		app.Session.Put(r.Context(), "flash", "Sin título en español no se publica nada. Normas de la casa.")
		http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
		return
	}

	if err := app.Storage.MovePhoto(r.Context(), p.ID, "private", "public"); err != nil {
		http.Error(w, "mover ficheros en S3: "+err.Error(), 500)
		return
	}
	if err := app.Photos.UpdateStatus(p.ID, "public"); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.Session.Put(r.Context(), "flash", "Publicada. Ya es de dominio… bueno, CC BY.")
	http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
}

func (app *AppRepository) PhotoUnpublish(w http.ResponseWriter, r *http.Request) {
	p, err := app.photoFromURL(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if p.Status == "private" {
		http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
		return
	}

	if err := app.Storage.MovePhoto(r.Context(), p.ID, "public", "private"); err != nil {
		http.Error(w, "mover ficheros en S3: "+err.Error(), 500)
		return
	}
	if err := app.Photos.UpdateStatus(p.ID, "private"); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.Session.Put(r.Context(), "flash", "Despublicada. Vuelve a ser solo tuya.")
	http.Redirect(w, r, fmt.Sprintf("/admin/fotos/%d", p.ID), http.StatusSeeOther)
}

func (app *AppRepository) PhotoDelete(w http.ResponseWriter, r *http.Request) {
	p, err := app.photoFromURL(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := app.Storage.DeletePhoto(r.Context(), p.ID, p.Status); err != nil {
		http.Error(w, "borrar de S3: "+err.Error(), 500)
		return
	}
	if err := app.Photos.Delete(p.ID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.Session.Put(r.Context(), "flash", "Borrada. Como si nunca hubiera existido.")
	http.Redirect(w, r, "/admin/fotos", http.StatusSeeOther)
}
