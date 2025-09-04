package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/utils"
)

var LoginHtml = template.Must(template.ParseFiles("templates/login.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<loginhandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := r.FormValue("password")
		err = utils.Authentification(db, login, password)
		if err != nil {
			if strings.Contains(err.Error(), "db") {
				utils.InternalServError(w)
			} else {
				NotValid := "Mot de passe ou nom d'utilisateur incorrect. Veuillez réessayer."
				LoginHtml.Execute(w, NotValid)
			}
		}
	} else {
		err := LoginHtml.Execute(w, nil)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template LoginHtml: %v\n", err)
			utils.NotFoundHandler(w)
		}
	}
}
