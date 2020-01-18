package app

import (
	"net/http"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	a.html(w, "home.page.html", "Hello")
}
