package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/services"
)

const LoginPageUri = "/login"

type LoginHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

func NewLoginHandler(userService services.UserService, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *LoginHandler) Register(router chi.Router) {
	router.Get(LoginPageUri, h.handle)
	router.Post(LoginPageUri, h.handlePost)
}

func (h *LoginHandler) handle(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/login.html")
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Template parsing error.", "err", err)
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	if err := tmp.ExecuteTemplate(w, "base", nil); err != nil {
		h.logger.ErrorContext(r.Context(), "Template render error.", "err", err)
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var errs []string
	if err := h.userService.LoginAdmin(r.Context(), r.FormValue("username")); err != nil {
		h.logger.ErrorContext(r.Context(), "Login admin error.", "err", err)
		errs = append(errs, err.Error())
	}

	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/login.html")
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Template parsing error.", "err", err)
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	pageData := map[string]any{
		"Errors":      errs,
		"MessageSent": len(errs) == 0,
	}
	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		h.logger.ErrorContext(r.Context(), "Template render error.", "err", err)
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
