package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/antonrodin/azr/internal/authsvc"
	"github.com/antonrodin/azr/internal/handlers"
	"github.com/antonrodin/azr/internal/models/mysqlite"
	"github.com/antonrodin/azr/internal/photos"
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

	log.Println("Initialize sessions")
	session = scs.New()
	// Sesión larga: la principal usuaria es la subida de fotos desde el
	// móvil y no apetece meter la contraseña cada día.
	session.Lifetime = 30 * 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = env == "prod"
	session.Cookie.HttpOnly = true

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	defer db.Close()
	log.Println("Connected to SQLite successfully.")

	photoModel := &mysqlite.PhotoModel{DB: db}
	if err := photoModel.EnsureSchema(); err != nil {
		return fmt.Errorf("crear esquema de photos: %w", err)
	}

	storage, err := photos.NewStorage(context.Background())
	if err != nil {
		return fmt.Errorf("almacenamiento S3: %w", err)
	}
	log.Println("S3 storage ready.")

	app := app{
		Session: session,
	}

	appRepo := handlers.AppRepository{
		Session: session,
		Pages:   &mysqlite.PageModel{DB: db},
		Photos:  photoModel,
		Storage: storage,
		Auth:    authsvc.New(authService, "azr"),
	}

	handlers.NewRepo(&appRepo)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Server is running on http://localhost:%s \n", port)

	return server.ListenAndServe()
}
