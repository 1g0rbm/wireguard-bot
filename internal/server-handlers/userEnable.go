package serverhandlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"wireguard-bot/internal/services"

	"github.com/go-chi/chi/v5"
)

type UserEnableHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

type userEnablePageData struct {
	Status bool
}

func NewUserEnableHandler(userService services.UserService, logger *slog.Logger) *UserEnableHandler {
	return &UserEnableHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserEnableHandler) Register(router chi.Router) {
	router.Post("/user/enable/{user_id}", h.handle)
}

func (h *UserEnableHandler) handle(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	pageData := userEnablePageData{
		Status: true,
	}
	if err := h.userService.Enable(r.Context(), userID); err != nil {
		h.logger.ErrorContext(r.Context(), "Could not to enable user.", "err", err)
		pageData.Status = false
	}

	tmp, err := template.ParseFiles("static/templates/base.html", "static/templates/user_enable.html")
	if err != nil {
		http.Error(w, "Template parsing error.", http.StatusInternalServerError)
		return
	}

	if err := tmp.ExecuteTemplate(w, "base", pageData); err != nil {
		http.Error(w, "Template render error.", http.StatusInternalServerError)
		return
	}
}
