package server_handlers

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (r *RootHandler) Register(router chi.Router) {
	router.Get("/", r.handle)
}

func (r *RootHandler) handle(w http.ResponseWriter, req *http.Request) {
	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/root.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	if err := tmp.ExecuteTemplate(w, "base", map[string]string{"Username": "Aaaa"}); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
