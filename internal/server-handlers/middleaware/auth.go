package middleaware

import (
	"net/http"

	"github.com/google/uuid"

	serverhandlers "wireguard-bot/internal/server-handlers"
	"wireguard-bot/internal/services"
)

const sessionCookieName = "session"

type Auth struct {
	sessionService services.SessionService
}

func NewAuth(sessionService services.SessionService) *Auth {
	return &Auth{
		sessionService: sessionService,
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

		if err := a.sessionService.Check(r.Context(), sessionID); err != nil {
			http.Redirect(w, r, serverhandlers.LoginPageURI, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
