package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/repository"
	"wireguard-bot/internal/repository/users2servers"
	"wireguard-bot/internal/utils/flash"
)

type UserPageHandler struct {
	users2serversRepo repository.Users2Servers
	logger            *slog.Logger
}

type userPageData struct {
	User  *users2servers.UsersServers
	Flash string
}

func NewUserPageHandler(users2serversRepo repository.Users2Servers, logger *slog.Logger) *UserPageHandler {
	return &UserPageHandler{
		users2serversRepo: users2serversRepo,
		logger:            logger,
	}
}

func (h *UserPageHandler) Register(router chi.Router) {
	router.Get("/user/{user_id}", h.handle)
}

func (h *UserPageHandler) handle(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	fullInfo, err := h.users2serversRepo.GetFullInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	pageData := userPageData{
		User:  fullInfo,
		Flash: flash.GetFlash(w, r, "user-enable-flash"),
	}

	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/user_page.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
