package app

import (
	"net/http"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	d := make(map[string]interface{})
	d["IsAuthenticated"] = false
	a.html(w, "home.page.html", d)
}

func (a *App) signup(w http.ResponseWriter, r *http.Request) {
	a.html(w, "signup.page.html", nil)
}
