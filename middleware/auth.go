package middleware

import (
	"context"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
)

// key type pour le contexte afin d'éviter les collisions
type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware protège les routes nécessitant une session
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		session, err := sessions.GetSession(cookie.Value)
		if err != nil || session.UserID == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Passer UserID via contexte type-safe
		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext récupère le UserID depuis le contexte
func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
