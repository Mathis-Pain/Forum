package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"golang.org/x/crypto/bcrypt"
)

// Fonction de connection
func Authentification(db *sql.DB, username string, password string) (models.User, error) {
	// Récupère l'ID et le mot de passe (crypté) à partir de l'identifiant
	user, err := getdata.GetUserInfoFromLogin(db, username)
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
