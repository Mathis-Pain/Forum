package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var templ = template.Must(template.ParseFiles("templates/signup.html"))

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirmpassword")

	var db *sql.DB

	//Gestion d'erreur
	erreur := utils.ValidName(name)

	data := struct {
		Error error
	}{
		Error: erreur,
	}

	if erreur != nil {
		templ.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)
	}

	erreur = utils.ValidEmail(email)
	if erreur != nil {
		templ.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)
	}

	erreur = utils.ValidPasswd(password, passwordConfirm)
	if erreur != nil {
		templ.Execute(w, data)
		w.WriteHeader(http.StatusBadRequest)

	}

	_, err := db.Exec("INSERT INTO users(name, email, password) VALUES(?, ?, ?)", name, email, password)
	if err != nil {
		http.Error(w, "Erreur DB: "+err.Error(), http.StatusInternalServerError)
		utils.InternalServError(w)
		return
	}

	//est-ce que vraiment ça va marcher ou alors il faut mettre un if err == nil pour être sûr
	fmt.Fprint(w, "Utilisateur ajouté") //renvoyer ça dans un template pour le stylisé ?
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
