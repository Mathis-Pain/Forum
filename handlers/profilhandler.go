package handlers

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var ProfilHtml = template.Must(template.ParseFiles("templates/profil.html"))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {

	err := ProfilHtml.Execute(w, nil)
	if err != nil {
		utils.InternalServError(w)
	}
}
