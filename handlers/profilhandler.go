package handlers

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var ProfilHtml = template.Must(template.New("profil.html").ParseFiles("templates/profil.html", "templates/login.html", "templates/header.html", "templates/initpage.html"))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {

	err := ProfilHtml.Execute(w, nil)
	if err != nil {
		utils.InternalServError(w)
	}
}
