package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/antonrodin/azr/internal/authsvc"
	"github.com/antonrodin/azr/internal/blog"
	"github.com/antonrodin/azr/internal/handlers"
	"github.com/antonrodin/azr/internal/models/mysqlite"
	"github.com/antonrodin/azr/internal/openai"
	"github.com/antonrodin/azr/internal/photos"
	"github.com/antonrodin/azr/internal/tmpl"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	Session *scs.SessionManager
}

var session *scs.SessionManager

// jwtSecret lo usa el middleware Auth. Debe coincidir con el del servicio auth.
var jwtSecret []byte

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbFile := os.Getenv("DB_FILE")
	env := os.Getenv("ENV")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("falta JWT_SECRET en .env")
	}
	jwtSecret = []byte(secret)

	authService := os.Getenv("AUTH_SERVICE")
	if authService == "" {
		return fmt.Errorf("falta AUTH_SERVICE en .env")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Println("Connected to SQLite successfully.")

	// Plantillas cacheadas en prod; en dev se reparsean en cada petición.
	tmpl.Cache = env == "prod"

	// Sesiones en SQLite (misma BD) para que sobrevivan a los reinicios de
	// supervisor y no haya que volver a hacer login tras cada despliegue.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		data BLOB NOT NULL,
		expiry REAL NOT NULL
	);
	CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);`); err != nil {
		return fmt.Errorf("crear tabla sessions: %w", err)
	}

	log.Println("Initialize sessions")
	session = scs.New()
	session.Store = sqlite3store.New(db)
	// Sesión larga: la principal usuaria es la subida de fotos desde el
	// móvil y no apetece meter la contraseña cada día.
	session.Lifetime = 30 * 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = env == "prod"
	session.Cookie.HttpOnly = true

	photoModel := &mysqlite.PhotoModel{DB: db}
	if err := photoModel.EnsureSchema(); err != nil {
		return fmt.Errorf("crear esquema de photos: %w", err)
	}

	promptModel := &mysqlite.PromptModel{DB: db}
	if err := promptModel.EnsureSchema(); err != nil {
		return fmt.Errorf("crear esquema de prompts: %w", err)
	}
	if err := promptModel.Seed(handlers.PromptRegistry); err != nil {
		return fmt.Errorf("sembrar prompts: %w", err)
	}

	storage, err := photos.NewStorage(context.Background())
	if err != nil {
		return fmt.Errorf("almacenamiento S3: %w", err)
	}
	log.Println("S3 storage ready.")

	// Cliente OpenAI (Fases 3 y 4: traducir y transcribir voz). Sin
	// OPENAI_API_KEY queda deshabilitado y el admin simplemente no ofrece esas
	// funciones; el resto funciona igual.
	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_MODEL"), os.Getenv("OPENAI_TRANSCRIBE_MODEL"))
	if oai.Enabled() {
		log.Println("OpenAI client ready.")
	} else {
		log.Println("OpenAI client disabled (no OPENAI_API_KEY).")
	}

	posts, err := blog.Load(handlers.BlogDir)
	if err != nil {
		return fmt.Errorf("cargar blog: %w", err)
	}
	log.Println("Blog loaded.")

	app := app{
		Session: session,
	}

	appRepo := handlers.AppRepository{
		Session: session,
		Pages:   &mysqlite.PageModel{DB: db},
		Photos:  photoModel,
		Prompts: promptModel,
		Storage: storage,
		Auth:    authsvc.New(authService, "azr"),
		OpenAI:  oai,
		Blog:    posts,
	}

	handlers.NewRepo(&appRepo)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           app.routes(),
		ReadHeaderTimeout: 10 * time.Second,
		// Generoso: subidas de fotos/audio desde el móvil con mala cobertura.
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  2 * time.Minute,
	}

	log.Printf("Server is running on http://localhost:%s \n", port)

	return server.ListenAndServe()
}
