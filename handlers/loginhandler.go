package handlers

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var loginHtml = template.Must(template.ParseFiles("templates/login.html"))

func SignUpFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// afficher le formulaire
		err := loginHtml.Execute(w, nil)
		if err != nil {
			utils.InternalServError(w)

		}
	}
}
