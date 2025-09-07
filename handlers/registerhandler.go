package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

var registrationHtml = template.Must(template.ParseFiles("templates/registration.html"))

// Struct pour transmettre les erreurs au template
type FormDataError struct {
	NameError  string
	EmailError string
	PassError  string
}

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// GET : afficher le formulaire vide
		if err := registrationHtml.Execute(w, nil); err != nil {
			utils.InternalServError(w)
		}
		return
	}

	// --- Récupération des valeurs ---
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirmpassword")

	// --- Struct pour stocker les erreurs ---
	formData := FormDataError{
		NameError:  utils.ValidName(username),
		EmailError: utils.ValidEmail(email),
		PassError:  utils.ValidPasswd(password, passwordConfirm),
	}

	// Si une erreur existe, renvoyer le formulaire avec messages
	if formData.NameError != "" || formData.EmailError != "" || formData.PassError != "" {
		w.WriteHeader(http.StatusBadRequest)
		registrationHtml.Execute(w, formData)
		return
	}

	// --- Ouverture de la DB ---
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// --- Insertion dans la DB ---
	_, err = db.Exec("INSERT INTO user(username, email, password) VALUES(?, ?, ?)", username, email, password)
	if err != nil {
		// Vérification UNIQUE (nom ou email déjà utilisé)
		if err.Error() == "UNIQUE constraint failed: user.username" {
			formData.NameError = "Nom d'utilisateur déjà utilisé"
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, formData)
			return
		} else if err.Error() == "UNIQUE constraint failed: user.email" {
			formData.EmailError = "email déjà utilisé"
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, formData)
			return
		}
		// Toute autre erreur
		utils.InternalServError(w)
		return
	}

	// --- Succès : redirection vers la page d'accueil ---
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
