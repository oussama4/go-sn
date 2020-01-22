package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/oussama4/go-sn/models"
)

type App struct {
	logger    *log.Logger
	templates map[string]*template.Template
	sm        *scs.SessionManager
	userStore models.UserStore
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

	sessionMan := scs.New()
	sessionMan.Store = postgresstore.New(db.DB)
	app := App{
		logger:    l,
		templates: cache,
		sm:        sessionMan,
		userStore: us,
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
