package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

var HomeHtml = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := HomeHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template HomeHtml: %v\n", err)
		utils.NotFoundHandler(w)
	}

}
