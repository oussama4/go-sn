package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/gocraft/dbr/v2"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/oussama4/go-sn/models"
)

type App struct {
	logger    *log.Logger
	db *dbr.Connection
	templates map[string]*template.Template
	sm        *scs.SessionManager
	userStore models.UserStore
	connStore models.ConnectionStore
	td        map[string]interface{}
}

// Start start the http server
func Start() {
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

	data := make(map[string]interface{})
	data["user"] = nil

	app := App{
		logger:    l,
		db: db,
		templates: cache,
		sm:        sessionMan,
		userStore: us,
		connStore: cs,
		td:        data,
	}

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
