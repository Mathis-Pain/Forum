package getdata

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère l'ID et le mot de passe (crypté) d'un utilisateur à partir de son identifiant (mail ou pseudo) pour la connexion
func GetUserInfoFromLogin(db *sql.DB, login string) (models.User, error) {
	// Préparation de la requête SQL : récupérer id, username et password
	sql := `SELECT id, username, password FROM user WHERE username = ?`
	row := db.QueryRow(sql, login)

	var user models.User
	// Parcourt la base de données en cherchant le username correspondant
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Récupère le nom d'utilisateur et la photo de profil d'un utilisateur à partir de son ID
func GetUserInfoFromID(db *sql.DB, ID int) (models.User, error) {
	// Préparation de la requête sql
	sql := `SELECT username FROM user WHERE id = ?`
	row := db.QueryRow(sql, ID)

	var user models.User
	err := row.Scan(&user.Username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
