package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

var HomeHtml = template.Must(template.ParseFiles("templates/home.html", "templates/login.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := utils.LoginPopUp(r, w)
		if err == nil {
			HomeHtml.Execute(w, data)
		}

		// Connexion réussie (ouverture de session, accès aux boutons, etc, à ajouter ici)

	} else {
		data := struct {
			Message   string
			ShowLogin bool
		}{
			Message:   "",    // No message on initial load
			ShowLogin: false, // The modal should be hidden initially
		}
		// S'il n'y a pas eu d'envoi du formulaire de connexion, affiche la page d'accueil de base
		err := HomeHtml.Execute(w, data)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template HomeHtml: %v\n", err)
			utils.NotFoundHandler(w)
		}
	}
}
