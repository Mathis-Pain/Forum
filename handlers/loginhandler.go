package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
)

var loginHtml = template.Must(template.ParseFiles("templates/login.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Affiche le formulaire de login
		if err := loginHtml.Execute(w, nil); err != nil {
			utils.InternalServError(w)
		}

	case http.MethodPost:
		// Parse les données
		err := r.ParseForm()
		if err != nil {
			utils.InternalServError(w)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		// Connexion DB
		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			utils.InternalServError(w)
			return
		}
		defer db.Close()

		// Vérifie login + mot de passe
		user, err := utils.Authentification(db, username, password)
		if err != nil {
			// Mauvais login ou mot de passe
			http.Error(w, "Nom d’utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Crée la session
		sessionData := map[string]interface{}{
			"user": user.Username, // ou user.ID si tu préfères
		}

		sessionID, err := sessions.CreateSession(sessionData)
		if err != nil {
			utils.InternalServError(w)
			return
		}

		// Ajoute le cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   false, // true si HTTPS
			Path:     "/",
		})

		// Redirection vers /home
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
