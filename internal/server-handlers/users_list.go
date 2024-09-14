package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/services"
	userService "wireguard-bot/internal/services/user"
	"wireguard-bot/internal/utils"
)

type UsersListHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

type usersListPageData struct {
	Users          []user.Model
	UsernameFilter string
	StateFilter    string
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
	pageData := usersListPageData{
		UsernameFilter: r.URL.Query().Get("username"),
		StateFilter:    r.URL.Query().Get("state"),
	}

	opts := make([]utils.FilterOption, 0)
	if pageData.UsernameFilter != "" {
		opts = append(opts, userService.WithUsername(pageData.UsernameFilter))
	}
	if pageData.StateFilter == user.DisabledState {
		opts = append(opts, userService.WithDisabled())
	} else if pageData.StateFilter == user.EnabledState {
		opts = append(opts, userService.WithEnabled())
	}

	users, err := h.userService.List(r.Context(), opts...)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Get users list error.", "err", err)
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}
	pageData.Users = users

	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/users_list.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
