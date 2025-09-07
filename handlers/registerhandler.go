package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var registrationHtml = template.Must(template.ParseFiles("templates/registration.html"))

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("confirmpassword")

		db, err := sql.Open("sqlite3", "./database/forum.db")
		if err != nil {
			http.Error(w, "Impossible d'ouvrir la DB", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		//Gestion d'erreur

		erreur := utils.ValidName(username)

		type FormDataError struct {
			NameError  string
			EmailError string
			PassError  string
		}
		data := struct {
			Error error
		}{
			Error: erreur,
		}

		if erreur != nil {
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, data)
			return
		}

		erreur = utils.ValidEmail(email)
		if erreur != nil {
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, data)
			return
		}

		erreur = utils.ValidPasswd(password, passwordConfirm)
		if erreur != nil {
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, data)
			return
		}

		_, err = db.Exec("INSERT INTO user(username, email, password) VALUES(?, ?, ?)", username, email, password)
		if err != nil {
			http.Error(w, "Erreur DB: "+err.Error(), http.StatusInternalServerError)
			utils.InternalServError(w)
			return
		}

		//est-ce que vraiment ça va marcher ou alors il faut mettre un if err == nil pour être sûr
		fmt.Fprint(w, "Utilisateur ajouté") //renvoyer ça dans un template pour le stylisé ?
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {

		err := registrationHtml.Execute(w, nil)
		if err != nil {
			utils.InternalServError(w)
			return
		}
	}
}
