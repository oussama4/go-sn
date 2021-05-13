package app

import (
	"net/http"
)

// AddUser gets a user from the database and adds it to template data if it's not added yet
func (a *App) AddUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := a.sm.GetInt(r.Context(), "user_id")
		if userID != 0 {
			if a.user == nil {
				user, err := a.userStore.Get(userID)
				if err != nil {
					a.logger.Println(err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				a.user = user
			}
		}
		next(w, r)
	}
}

// LoginRequired rejects any request from an unauthenticated user
func (a *App) LoginRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := a.sm.GetInt(r.Context(), "user_id")
		if userID != 0 {
			a.isAuthenticated = true
			next(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
