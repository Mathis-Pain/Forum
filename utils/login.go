package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
	"golang.org/x/crypto/bcrypt"
)

// Fonction de connection
func Authentification(db *sql.DB, username string, password string) (models.User, error) {
	// Récupère l'ID et le mot de passe (crypté) à partir de l'identifiant
	user, err := GetUserInfoFromLogin(db, username)
	if errors.Is(err, sql.ErrNoRows) {
		// Si aucun utilisateur n'est trouvé avec cet identifiant (mail ou pseudo), renvoie une erreur
		log.Printf("<login.go> : login failed, User %v not found\n", username)
		return models.User{}, fmt.Errorf("incorrect password or username")
	} else if err != nil {
		// Erreur dans la base de données
		mylog := fmt.Errorf("(db) could not recover user infos from the database %v", err)
		log.Println("ERROR <login.go>:", mylog)
		return models.User{}, mylog
	}

	// Fonction bcrypt pour comparer le mot de passe entré par l'utilisateur avec celui présent dans la base de données
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("<login.go> : password does not match")
		return models.User{}, fmt.Errorf("incorrect password or username")
	}

	return user, err
}

// Fonction d'affichage du popup de connexion
func LoginPopUp(r *http.Request, w http.ResponseWriter) (models.LoginData, error) {
	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./database/forum.db")
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
