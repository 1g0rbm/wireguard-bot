package middleaware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	serverhandlers "wireguard-bot/internal/server-handlers"
	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils/sessionctx"
)

const (
	sessionCookieName  = "session"
	contextUsernameKey = "username"
)

type Auth struct {
	userService services.UserService
	logger      *slog.Logger
}

func NewAuth(userService services.UserService, logger *slog.Logger) *Auth {
	return &Auth{
		userService: userService,
		logger:      logger,
	}
}

func (a *Auth) HandleFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Redirect(w, r, serverhandlers.LoginPageURI, http.StatusFound)
			return
		}

		sessionID, err := uuid.Parse(c.Value)
		if err != nil {
			http.Redirect(w, r, serverhandlers.LoginPageURI, http.StatusFound)
			return
		}

		user, err := a.userService.FindLoggedIn(r.Context(), sessionID)
		if err != nil {
			http.Redirect(w, r, serverhandlers.LoginPageURI, http.StatusFound)
			a.logger.ErrorContext(r.Context(), "Authorization middleware error.", "err", err)
			return
		}
		if user == nil {
			a.logger.ErrorContext(r.Context(), "Authorization middleware error.", "err", err)
			return
		}

		next.ServeHTTP(w, r.WithContext(sessionctx.WithUsername(r.Context(), user.Username)))
	})
}
