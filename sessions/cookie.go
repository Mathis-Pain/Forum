package sessions

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, name, value string, secure bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   secure, // configurable
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
}
