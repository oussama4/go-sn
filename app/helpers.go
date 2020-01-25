package app

import (
	"bytes"
	"net/http"
	"os"
	"strings"
)

type staticFileServer struct {
	fs http.FileSystem
}

func (sfs staticFileServer) Open(name string) (http.File, error) {
	f, err := sfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

// Static handles static files requests
func Static(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileserver := http.FileServer(staticFileServer{http.Dir(path)})
		http.StripPrefix("/static", fileserver).ServeHTTP(w, r)
	}
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

// DeleteEmpty takes a map of strings and  returns it after removing empty string values
func DeleteEmpty(m map[string]string) map[string]string {
	for k, v := range m {
		if strings.TrimSpace(v) == "" {
			delete(m, k)
		}
	}
	return m
}
