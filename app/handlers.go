package app

import (
	"net/http"

	"github.com/oussama4/go-sn/models"
	"github.com/oussama4/go-sn/pkg/forms"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	d := make(map[string]interface{})
	d["IsAuthenticated"] = false
	a.html(w, "home.page.html", d)
}

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
		data := make(map[string]interface{})
		data["form"] = f
		a.html(w, "signup.page.html", data)
		return
	}

	err = a.userStore.Insert(f.Get("name"), f.Get("email"), f.Get("pass1"))
	if err == models.ErrUsernameExist {
		f.Errors.Add("name", "username already in use")
		data := make(map[string]interface{})
		data["form"] = f
		a.html(w, "signup.page.html", data)
		return
	} else if err == models.ErrEmailExist {
		f.Errors.Add("email", "email already in use")
		data := make(map[string]interface{})
		data["form"] = f
		a.html(w, "signup.page.html", data)
		return
	} else if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
