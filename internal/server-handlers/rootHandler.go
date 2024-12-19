package serverhandlers

import (
	"html/template"
	"net/http"
	"wireguard-bot/internal/utils/sessionctx"

	"github.com/go-chi/chi/v5"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (rh *RootHandler) Register(router chi.Router) {
	router.Get("/", rh.handle)
}

func (rh *RootHandler) handle(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/root.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	username := sessionctx.ExtractUsername(r.Context())
	if err := tmp.ExecuteTemplate(w, "base", map[string]string{"Username": username}); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
