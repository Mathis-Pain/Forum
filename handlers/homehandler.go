package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

var HomeHtml = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<homehandler.go> Could not open database : %v\n", err)
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
				HomeHtml.Execute(w, NotValid)
			}
		}
	} else {
		err := HomeHtml.Execute(w, nil)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template HomeHtml: %v\n", err)
			utils.NotFoundHandler(w)
		}
	}
}
