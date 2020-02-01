package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gocraft/dbr/v2"
	"github.com/oussama4/go-sn/models"
	"github.com/oussama4/go-sn/pkg/forms"
)

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	userID := a.sm.GetInt(r.Context(), "user_id")
	if userID != 0 {
		a.isAuthenticated = true
	}

	a.html(w, "home.page.html", M{})
}

// serves the signup page
func (a *App) signup(w http.ResponseWriter, r *http.Request) {
	a.html(w, "signup.page.html", M{})
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
		data := M{"form": f}
		a.html(w, "signup.page.html", data)
		return
	}

	err = a.userStore.Insert(f.Get("name"), f.Get("email"), f.Get("pass1"))
	if err == models.ErrUsernameExist {
		f.Errors.Add("name", "username already in use")
		data := M{"form": f}
		a.html(w, "signup.page.html", data)
		return
	} else if err == models.ErrEmailExist {
		f.Errors.Add("email", "email already in use")
		data := M{"form": f}
		a.html(w, "signup.page.html", data)
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
		data := M{"form": f}
		a.html(w, "home.page.html", data)
		return
	}

	id, err := a.userStore.Authenticate(f.Get("email"), f.Get("pass"))
	if err == models.ErrInvalidCredentials || err == dbr.ErrNotFound {
		f.Errors.Add("email", "Invalid credentials")
		data := M{"form": f}
		a.html(w, "home.page.html", data)
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
	a.html(w, "profile.page.html", M{})
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

	data := M{
		"connected":  connected,
		"other_user": user,
	}
	a.html(w, "profile.page.html", data)
}

func (a *App) HandleSettingsPage(w http.ResponseWriter, r *http.Request) {
	a.html(w, "settings.page.html", M{})
}

func (a *App) HandleSettings(w http.ResponseWriter, r *http.Request) {
	fn := ""
	f, fh, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err == nil {
		defer f.Close()
		fn = filepath.Join("img/avatars", fh.Filename)
	}

	// fields that are going to be updated
	updated := DeleteEmpty(map[string]string{
		"name":   r.PostFormValue("name"),
		"bio":    r.PostFormValue("bio"),
		"avatar": fn,
	})

	if len(updated) == 0 {
		data := M{"update_error": "you didn't provide anything new"}
		a.html(w, "settings.page.html", data)
		return
	}

	userID := a.sm.GetInt(r.Context(), "user_id")
	if err = a.userStore.Update(userID, updated); err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if fn != "" {
		out, err := os.Create(filepath.Join(os.Getenv("STATIC_PATH"), fn))
		defer out.Close()
		_, err = io.Copy(out, f)
		if err != nil {
			a.logger.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}

// HandleCreateActivity handles the creation of activities of type CREATE
func (a *App) HandleCreateActivity(w http.ResponseWriter, r *http.Request) {
	fn := ""
	noFile := false
	txt := r.PostFormValue("txt")
	txt = strings.TrimSpace(txt)
	userID := a.sm.GetInt(r.Context(), "user_id")
	f, fh, err := r.FormFile("img")
	if err == http.ErrMissingFile {
		noFile = true
	} else if err != nil {
		a.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err == nil {
		defer f.Close()
		fn = filepath.Join("img/post", fh.Filename)
	}

	if txt == "" && noFile {
		data := M{"empty_payload": true}
		a.html(w, "home.page.html", data)
		return
	} else {
		ac := models.NewCreateActivity(userID, models.CREATE, txt, fn)
		if err := ac.Save(a.db); err != nil {
			a.logger.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// upload the image
		if !noFile {
			out, err := os.Create(filepath.Join(os.Getenv("STATIC_PATH"), fn))
			defer out.Close()
			_, err = io.Copy(out, f)
			if err != nil {
				a.logger.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
