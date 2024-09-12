package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/services/session"
)

const LoginPageURI = "/login"

type LoginHandler struct {
	userService    services.UserService
	sessionService services.SessionService
	logger         *slog.Logger
}

func NewLoginHandler(
	userService services.UserService, sessionService services.SessionService, logger *slog.Logger,
) *LoginHandler {
	return &LoginHandler{
		userService:    userService,
		sessionService: sessionService,
		logger:         logger,
	}
}

func (h *LoginHandler) Register(router chi.Router) {
	router.Get(LoginPageURI, h.handle)
	router.Post(LoginPageURI, h.handlePost)
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
	username := r.FormValue("username")

	sessID, err := h.sessionService.CheckByUsername(r.Context(), username)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessID.String(),
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			Expires:  time.Now().Add(session.TTL),
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := h.userService.LoginAdmin(r.Context(), username); err != nil {
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
		"Username":    username,
	}
	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		h.logger.ErrorContext(r.Context(), "Template render error.", "err", err)
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
