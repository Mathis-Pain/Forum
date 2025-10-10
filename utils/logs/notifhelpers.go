package logs

import "database/sql"

// Fonction pour trouver l'utilisateur à qui envoyer la notification pour les likes et les dislikes
func GetUserToNotify(messageID int, db *sql.DB) (int, error) {
	sqlQuery := `SELECT user_id FROM message WHERE id = ?`
	row := db.QueryRow(sqlQuery, messageID)

	var authorID int
	err := row.Scan(&authorID)
	if err != nil {
		return 0, err
	}

	return authorID, nil
}

// Fonction pour convertir le type de notification en donnée chiffrée pour la base de données
func ConvertNotifType(notifType string) int {
	switch notifType {
	case "ADMIN":
		return 1
	case "REQUEST":
		return 2
	case "ANSWER":
		return 3
	case "MESSAGE":
		return 4
	case "INTERACTION":
		return 5
	}

	return 0
}
