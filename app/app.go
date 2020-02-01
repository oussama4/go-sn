package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	"github.com/oussama4/go-sn/models"
)

type App struct {
	isAuthenticated bool
	logger          *log.Logger
	db              *dbr.Connection
	user            *models.User
	sm              *scs.SessionManager
	templates       map[string]*template.Template
	userStore       models.UserStore
	connStore       models.ConnectionStore
}

// Start creates an App instance and starts an http server
func Start() {
	// create the app
	l := log.New(os.Stdout, "LOGGER: ", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		l.Fatalln(err)
	}

	cache, err := newTemplateCache(os.Getenv("TEMPLATES_PATH"))
	if err != nil {
		l.Fatalln(err)
	}

	db, err := models.New(l)
	if err != nil {
		l.Fatal(err)
	}
	us := models.NewUserStore(l, db)
	cs := models.NewConnStore(l, db)

	sessionMan := scs.New()
	sessionMan.Store = postgresstore.New(db.DB)

	app := App{
		isAuthenticated: false,
		logger:          l,
		db:              db,
		user:            nil,
		sm:              sessionMan,
		templates:       cache,
		userStore:       us,
		connStore:       cs,
	}

	// start the http server
	srv := &http.Server{
		Addr:    os.Getenv("ADDRESS"),
		Handler: app.sm.LoadAndSave(app.routes()),
	}

	app.logger.Printf("starting server on %s", os.Getenv("ADDRESS"))
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatalln(err)
	}
}
