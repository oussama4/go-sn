package app

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gocraft/dbr/v2"
	"github.com/oussama4/go-sn/models"
	"github.com/oussama4/go-sn/pkg/forms"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := false

	userID := a.sm.GetInt(r.Context(), "user_id")
	if userID != 0 {
		isAuthenticated = true
	}

	a.td["IsAuthenticated"] = isAuthenticated
	a.html(w, "home.page.html", a.td)
}

// serves the signup page
func (a *App) signup(w http.ResponseWriter, r *http.Request) {
	a.html(w, "signup.page.html", nil)
}

// TODO: simplify this handler
func (a *App) handleSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	f := forms.New(r.PostForm)
	f.Required("name", "email", "pass1", "pass2")
	f.MinLength("name", 5)
	f.MinLength("pass1", 8)
	f.ValidEmail("email")
	f.StringsMatch("pass1", "pass2")

	if !f.Valid() {
		a.td["form"] = f
		a.html(w, "signup.page.html", a.td)
		return
	}

	err = a.userStore.Insert(f.Get("name"), f.Get("email"), f.Get("pass1"))
	if err == models.ErrUsernameExist {
		f.Errors.Add("name", "username already in use")
		a.td["form"] = f
		a.html(w, "signup.page.html", a.td)
		return
	} else if err == models.ErrEmailExist {
		f.Errors.Add("email", "email already in use")
		a.td["form"] = f
		a.html(w, "signup.page.html", a.td)
		return
	} else if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	f := forms.New(r.PostForm)
	f.Required("email", "pass")
	if !f.Valid() {
		a.td["form"] = f
		a.html(w, "home.page.html", a.td)
		return
	}

	id, err := a.userStore.Authenticate(f.Get("email"), f.Get("pass"))
	if err == models.ErrInvalidCredentials || err == dbr.ErrNotFound {
		f.Errors.Add("email", "Invalid credentials")
		a.td["form"] = f
		a.html(w, "home.page.html", a.td)
		return
	} else if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	a.sm.Put(r.Context(), "user_id", id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := a.sm.Destroy(r.Context())
	if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *App) HandleProfile(w http.ResponseWriter, r *http.Request) {
	a.html(w, "profile.page.html", a.td)
}

// HandleOtherProfile handles a request from the current authenticated user to other profiles
func (a *App) HandleOtherProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
	userID2 := a.sm.GetInt(r.Context(), "user_id")

	user, err := a.userStore.Get(userID)
	if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	connected, err := a.connStore.AreConnected(userID, userID2)
	if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	a.td["connected"] = connected
	a.td["other_user"] = user
	a.html(w, "profile.page.html", a.td)
}
