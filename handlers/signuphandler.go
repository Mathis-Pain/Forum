package handlers

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/utils"
)

var signUpHtml = template.Must(template.ParseFiles("templates/signup.html"))

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// afficher le formulaire
		err := signUpHtml.Execute(w, nil)
		if err != nil {
			utils.InternalServError(w)

		}
	}
}
