package authhandlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

var funcMap = template.FuncMap{
	"preview": getdata.Preview,
}

var HomeHtml = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles(
	"templates/home.html", "templates/login.html", "templates/header.html", "templates/initpage.html",
))

// MARK: Connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			utils.InternalServError(w)
			return
		}
		// Récupération des données
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Vérifie login + mot de passe (utils.Authentification s’occupe de la DB)
		db, err := sql.Open("sqlite3", "./data/forum.db")
		if err != nil {
			utils.InternalServError(w)
			return
		}
		defer db.Close()

		user, loginErr := utils.Authentification(db, username, password)
		if loginErr != nil {
			if strings.Contains(loginErr.Error(), "db") {
				// En cas d'erreur dans la base de données
				utils.InternalServError(w)
			} else {
				// En cas d'erreur qui ne vient pas de la base de données
				// Création d'une session temporaire et anonyme
				err := InitSession(w, 0, "LoginErr", loginErr.Error())
				if err != nil {
					utils.InternalServError(w)
					return
				}

				// Redirection vers la page d'origine
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
		}

		// Invalider toutes les sessions existantes
		if err := sessions.InvalidateUserSessions(user.ID); err != nil {
			utils.InternalServError(w)
			return
		}

		err = InitSession(w, user.ID, "user", user.Username)
		if err != nil {
			utils.InternalServError(w)
			return
		}

		http.Redirect(w, r, referer, http.StatusSeeOther)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// MARK: Création session
// Crée, sauvegarde une session, permet d'y insérer une donnée et pose le cookie
func InitSession(w http.ResponseWriter, id int, fieldName string, fieldData any) error {
	session, err := sessions.CreateSession(id)
	if err != nil {
		return err
	}

	session.Data[fieldName] = fieldData
	if err := sessions.SaveSessionToDB(session); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false, // false en local, true si HTTPS
		Path:     "/",
	})
	return nil
}
