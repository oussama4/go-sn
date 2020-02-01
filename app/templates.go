package app

import (
	"html/template"
	"path/filepath"
)

// M contains data that is going to be apllied to a template
type M map[string]interface{}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		t, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		t, err = t.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		/* t, err = t.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		} */

		cache[name] = t
	}
	return cache, nil
}
