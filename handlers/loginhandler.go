package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
)

var loginHtml = template.Must(template.ParseFiles("templates/login.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// si l'utilisateur demande le formulaire
	case http.MethodGet:
		if err := loginHtml.Execute(w, nil); err != nil {
			utils.InternalServError(w)
			return
		}
		// Si l'utilisateur envoi le formulaire
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			utils.InternalServError(w)
			return
		}
		// Verification username et password non nul
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		// Vérifie login + mot de passe (utils.Authentification s’occupe de la DB)
		db, err := sql.Open("sqlite3", "data/forum.db")
		if err != nil {
			utils.InternalServError(w)
			return
		}
		defer db.Close()

		user, err := utils.Authentification(db, username, password)
		if err != nil {
			http.Error(w, "Nom d’utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}
		//Invalider toutes les sessions existantes
		if err := sessions.InvalidateUserSessions(user.ID); err != nil {
			utils.InternalServError(w)
			return
		}

		// Créer une nouvelle session
		session, err := sessions.CreateSession(user.ID)
		if err != nil {
			utils.InternalServError(w)
			return
		}

		// Ajoute des infos dans la session
		session.Data["user"] = user.Username
		if err := sessions.SaveSessionToDB(session); err != nil {
			utils.InternalServError(w)
			return
		}

		// Pose le cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   false, // ⚠️ false en local, true si HTTPS
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
