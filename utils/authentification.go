package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"

	"golang.org/x/crypto/bcrypt"
)

// Fonction de connection
func Authentification(db *sql.DB, username string, password string) (models.User, error) {
	if username == "" || password == "" {
		mylog := fmt.Errorf("tous les champs sont requis")
		logMsg := fmt.Sprintln("ERREUR : <authentification.go> ", mylog)
		logs.AddLogsToDatabase(logMsg)
		return models.User{}, mylog
	}
	// Récupère l'ID et le mot de passe (crypté) à partir de l'identifiant
	user, err := getdata.GetUserInfoFromLogin(db, username)
	if errors.Is(err, sql.ErrNoRows) {
		// Si aucun utilisateur n'est trouvé avec cet identifiant (mail ou pseudo), renvoie une erreur
		logMsg := fmt.Sprintf("ERREUR : <authentification.go> Tentative de connexion échouée : L'utilisateur %s n'existe pas.\n", username)
		logs.AddLogsToDatabase(logMsg)
		return models.User{}, fmt.Errorf("nom d'utilisateur incorrect")
	} else if err != nil {
		// Erreur dans la base de données
		mylog := fmt.Errorf("(db) Impossible de récupérer les données utilisateur dans la base de données : %v", err)
		logMsg := fmt.Sprintln("ERREUR : <authentification.go> ", mylog)
		logs.AddLogsToDatabase(logMsg)
		return models.User{}, mylog
	}

	// Fonction bcrypt pour comparer le mot de passe entré par l'utilisateur avec celui présent dans la base de données
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <authentification.go> : Mot de passe incorrect")
		logs.AddLogsToDatabase(logMsg)
		return models.User{}, fmt.Errorf("mot de passe incorrect")
	}

	return user, err
}
