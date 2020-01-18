package app

import (
	"net/http"
	"os"
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
