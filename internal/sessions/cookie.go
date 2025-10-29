package sessions

import (
	"net/http"
	"time"
)

// SetCookie crée un cookie de session
func SetCookie(w http.ResponseWriter, name, value string, secure bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/", // même Path partout
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   24 * 60 * 60, // 24 heures
		HttpOnly: true,         // protège côté JS
		Secure:   secure,       // true si HTTPS, false en local
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// DeleteCookie supprime un cookie existant
func DeleteCookie(w http.ResponseWriter, name string, secure bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",             // doit être identique à SetCookie
		Expires:  time.Unix(0, 0), // date passée = suppression
		MaxAge:   -1,              // suppression immédiate
		HttpOnly: true,            // doit correspondre à SetCookie
		Secure:   secure,          // doit correspondre à SetCookie
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
