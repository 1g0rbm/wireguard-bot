package flash

import (
	"net/http"
	"net/url"
)

func GetFlash(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}

	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	msg, _ := url.QueryUnescape(cookie.Value)
	return msg
}

func SetFlash(w http.ResponseWriter, name, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: url.QueryEscape(message),
		Path:  "/",
	})
}
