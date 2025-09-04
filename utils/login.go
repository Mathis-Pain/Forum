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

func Authentification(db *sql.DB, login string, password string) error {
	user, err := GetUserInfoFromLogin(db, login)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("<login.go l.35 : login failed, User %v not found\n", login)
		return fmt.Errorf("incorrect password or username")
	} else if err != nil {
		mylog := fmt.Errorf("could not recover user infos from the database %v", err)
		log.Println("ERROR <login.go> l39 :", mylog)
		return mylog
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("<login.go> l.45 : password does not match")
		return fmt.Errorf("incorrect password or username")
	}

	return nil
}

func LoginPopUp(r *http.Request, w http.ResponseWriter) (models.LoginData, error) {
	// Tentative de connexion
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<homehandler.go> Could not open database: %v\n", err)
		InternalServError(w)
		return models.LoginData{}, err
	}
	defer db.Close()

	login := r.FormValue("login")
	password := r.FormValue("password")
	err = Authentification(db, login, password)

	if err != nil {
		// Erreur dans la tentative de connexion
		if strings.Contains(err.Error(), "db") {
			InternalServError(w)
		} else {
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
