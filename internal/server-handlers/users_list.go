package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"wireguard-bot/internal/services"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/repository/user"
)

type UsersListHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

type usersListPageData struct {
	Users []user.Model
}

func NewUsersListHandler(userService services.UserService, logger *slog.Logger) *UsersListHandler {
	return &UsersListHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UsersListHandler) Register(router chi.Router) {
	router.Get("/users", h.handle)
}

func (h *UsersListHandler) handle(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.List(r.Context())
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Get users list error.", "err", err)
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/users_list.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	pageData := usersListPageData{
		Users: users,
	}
	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
