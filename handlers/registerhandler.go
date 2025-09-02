package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

// Handler regroupe toutes les dépendances dont les routes ont besoin

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirmpassword")

	var db *sql.DB

	//Gestion d'erreur
	erreur := utils.ValidName(name)
	if erreur != nil {
		fmt.Fprint(w, erreur)
		w.WriteHeader(http.StatusBadRequest)
	}

	erreur = utils.ValidEmail(email)
	if erreur != nil {
		fmt.Fprint(w, erreur)
		w.WriteHeader(http.StatusBadRequest)
	}

	erreur = utils.ValidPasswd(password, passwordConfirm)
	if erreur != nil {
		fmt.Fprint(w, erreur)
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err := db.Exec("INSERT INTO users(name, email, password) VALUES(?, ?, ?)", name, email, password)
	if err != nil {
		http.Error(w, "Erreur DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Utilisateur ajouté")
}
