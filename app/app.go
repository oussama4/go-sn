package app

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	logger    *log.Logger
	templates map[string]*template.Template
}

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/static/*", Static("./ui/static"))
	r.Get("/", a.index)
	r.Get("/signup", a.signup)

	return r
}

// html renders an html template
func (a *App) html(w http.ResponseWriter, name string, data interface{}) {
	t, ok := a.templates[name]
	if !ok {
		a.logger.Printf("template %s does not exist", name)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, data)
	if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Start start the http server
func Start() {
	l := log.New(os.Stdout, "LOGGER: ", log.Ldate|log.Ltime|log.Lshortfile)
	cache, err := newTemplateCache("./ui/templates")
	if err != nil {
		l.Fatalln(err)
	}
	app := App{
		logger:    l,
		templates: cache,
	}

	srv := &http.Server{
		Addr:    os.Getenv("ADDRESS"),
		Handler: app.routes(),
	}

	app.logger.Printf("starting server on %s", os.Getenv("ADDRESS"))
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatalln(err)
	}
}
