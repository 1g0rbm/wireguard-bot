package serverhandlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils/flash"
)

type UserEnableHandler struct {
	userService services.UserService
	logger      *slog.Logger
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

	if err := h.userService.Enable(r.Context(), userID); err != nil {
		h.logger.ErrorContext(r.Context(), "Could not to enable user.", "err", err)
		flash.SetFlash(w, "user-enable-flash", "error")
		http.Redirect(w, r, fmt.Sprintf("/user/%d", userID), http.StatusSeeOther)
		return
	}

	flash.SetFlash(w, "user-enable-flash", "success")
	http.Redirect(w, r, fmt.Sprintf("/user/%d", userID), http.StatusSeeOther)
}
