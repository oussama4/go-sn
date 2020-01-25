package app

import (
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/static/*", Static(os.Getenv("STATIC_PATH")))
	r.Get("/", a.AddUser(a.index))
	r.Get("/signup", a.signup)
	r.Post("/signup", a.handleSignup)
	r.Post("/login", a.HandleLogin)
	r.Get("/logout", a.handleLogout)
	r.Get("/profile", a.LoginRequired(a.AddUser(a.HandleProfile)))
	r.Get("/profile/{userID}", a.LoginRequired(a.AddUser(a.HandleOtherProfile)))
	r.Get("/profile/settings", a.LoginRequired(a.AddUser(a.HandleSettingsPage)))
	r.Post("/profile/edit", a.LoginRequired(a.AddUser(a.HandleSettings)))

	return r
}
