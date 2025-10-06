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
	sql := `SELECT username, profilpic, role_id FROM user WHERE id = ?`
	row := db.QueryRow(sql, ID)

	var user models.User
	err := row.Scan(&user.Username, &user.ProfilPic, &user.Status)

	user.Status = SetUserStatus(user.Status)
	user.ID = ID

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Transforme l'ID du role en nom de rôle
func SetUserStatus(role string) string {
	switch role {
	case "1":
		role = "Admin"
	case "2":
		role = "Modérateur"
	case "3":
		role = "Membre"
	case "4":
		role = "Banni"
	default:
		// Le rôle 5 contient un espace et sert à différencier les membres ayant demandé à rejoindre la modération
		// des membres qui n'ont pas fait cette demande
		role = "Membre "
	}

	return role
}

// Transforme le nom du statut en chiffre pour le stocker dans la base de données
func CodeUserStatus(role string) string {
	switch role {
	case "Admin":
		role = "1"
	case "Modérateur":
		role = "2"
	case "Membre":
		role = "3"
	case "Banni":
		role = "4"
	default:
		role = "3"
	}

	return role
}
