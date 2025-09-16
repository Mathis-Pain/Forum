package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var registrationHtml = template.Must(template.New("registration.html").Funcs(funcMap).ParseFiles("templates/registration.html", "templates/login.html", "templates/header.html", "templates/initpage.html"))

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
	profilPic := r.FormValue("userimage")

	if profilPic == "" {
		profilPic = "static/noprofilpic.png"
	}

	// --- Struct pour stocker les erreurs ---
	formData := models.RegisterDataError{
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServError(w)
		return
	}
	// --- Insertion dans la DB ---
	_, err = db.Exec("INSERT INTO user(username, email, password) VALUES(?, ?, ?)", username, email, hashedPassword)
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
