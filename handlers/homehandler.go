package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

var HomeHtml = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles("templates/home.html", "templates/login.html", "templates/header.html", "templates/initpage.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := utils.LoginPopUp(r, w)
		if err == nil {
			HomeHtml.Execute(w, data)
		}

		// Connexion réussie (ouverture de session, accès aux boutons, etc, à ajouter ici)

	} else {

		lastPosts, err := utils.GetLastPosts()

		if err == sql.ErrNoRows {
			data := struct {
				LoginData models.LoginData
				Posts     []models.LastPost
			}{
				LoginData: models.LoginData{},
				Posts:     lastPosts,
			}

			err = HomeHtml.Execute(w, data)

		} else if err != nil {
			utils.InternalServError(w)
		}

		data := struct {
			LoginData models.LoginData
			Posts     []models.LastPost
		}{
			LoginData: models.LoginData{},
			Posts:     lastPosts,
		}
		// S'il n'y a pas eu d'envoi du formulaire de connexion, affiche la page d'accueil de base
		err = HomeHtml.Execute(w, data)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template HomeHtml: %v\n", err)
			utils.NotFoundHandler(w)
		}
	}
}
