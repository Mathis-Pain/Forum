package utils

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
)

// Fonction d'affichage du popup de connexion
func LoginPopUp(r *http.Request, w http.ResponseWriter) (models.LoginData, error) {
	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "/data/forum.db")
	if err != nil {
		log.Printf("<login.go> Could not open database: %v\n", err)
		InternalServError(w)
		return models.LoginData{}, err
	}
	defer db.Close()

	// Récupération des informations du formulaire
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Vérifie que l'utilisateur est enregistré et que le mot de passe correspond
	_, err = Authentification(db, username, password)

	if err != nil {
		if strings.Contains(err.Error(), "db") {
			// En cas d'erreur dans la base de données
			InternalServError(w)
		} else {
			// En cas d'erreur qui ne vient pas de la base de données
			data := models.LoginData{
				Message:   "Mot de passe ou nom d'utilisateur incorrect. Veuillez réessayer.",
				ShowLogin: true,
			}

			return data, nil
		}
		return models.LoginData{}, err
	}

	return models.LoginData{}, nil
}
