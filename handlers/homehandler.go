package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var HomeHtml = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	err := HomeHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template HomeHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
