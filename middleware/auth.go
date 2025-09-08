package middleware

import (
	"context"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
)

// AuthMiddleware protège les routes nécessitant une session
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupère le cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// Pas de cookie → redirige vers login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Récupère la session
		session, err := sessions.GetSession(cookie.Value)
		if err != nil {
			// Cookie invalide ou session expirée
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Ajoute les infos utilisateur dans le contexte
		ctx := context.WithValue(r.Context(), "user", session.Data["user"])
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
