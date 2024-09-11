package middleaware

import (
	"net/http"
	"strconv"

	serverhandlers "wireguard-bot/internal/server-handlers"
	"wireguard-bot/internal/services"
)

const (
	sessionCookieName = "session"
	queryUserID       = "user_id"
)

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
		var userID int64
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			userID, err = strconv.ParseInt(r.URL.Query().Get(queryUserID), 10, 64)
			if err != nil || userID == 0 {
				http.Redirect(w, r, serverhandlers.LoginPageUri, http.StatusFound)
				return
			}
		} else {
			userID, err = strconv.ParseInt(cookie.Value, 10, 64)
			if err != nil || userID == 0 {
				http.Redirect(w, r, serverhandlers.LoginPageUri, http.StatusFound)
				return
			}
		}

		if err := a.sessionService.Check(r.Context(), userID); err != nil {
			http.Redirect(w, r, serverhandlers.LoginPageUri, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
