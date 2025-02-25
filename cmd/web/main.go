package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/antonrodin/azr/internal/handlers"
	"github.com/antonrodin/azr/internal/models/mysqlite"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	Title   string
	Session *scs.SessionManager
	pages   *mysqlite.PageModel
}

var session *scs.SessionManager

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

	log.Println("Initialize sessions")
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	log.Println("Connected to SQLite successfully.")

	app := app{
		Title:   "Nurburgring. Un pequeño paraiso.",
		Session: session,
		pages: &mysqlite.PageModel{
			DB: db,
		},
	}

	appRepo := handlers.AppRepository{
		Session: session,
		Pages: &mysqlite.PageModel{
			DB: db,
		},
	}

	handlers.NewRepo(&appRepo)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Server is running on http://localhost:%s \n", port)

	return server.ListenAndServe()
}
