package main

import (
	"sync"
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/stretchr/objx"
)

// templateHandler is an handler represent a single template
type templateHandler struct {
	once sync.Once // once will Must the template only once.
	filename string // filename is the filename of the template
	templ *template.Template // templ is a template.Template object storage so that it will parse only once.
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func(){
		t.templ = template.Must(
			template.ParseFiles(
				filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromURLQuery(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
