package getdata

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère la liste complète des messages (date de création, auteur, contenu) pour un sujet
func GetMessageList(db *sql.DB, topicID int) ([]models.Message, error) {
	// Préparation de la requête sql
	sqlQuery := `SELECT created_at, user_id, content, IFNULL(likes, 0), IFNULL(dislikes, 0), id, warning FROM message WHERE topic_id = ?`
	rows, err := db.Query(sqlQuery, topicID)
	if err != nil {
		return nil, err
	}

	var messages []models.Message

	// Parcourt la base de données et récupère les informations pour rajouter tous les messages dans la slice
	for rows.Next() {
		var message models.Message
		user_id := 0
		if err := rows.Scan(&message.Created, &user_id, &message.Content, &message.Likes, &message.Dislikes, &message.MessageID, &message.Warning); err != nil {
			log.Printf("ERREUR : <getmessagelist.go> Erreur dans le parcours de la base de données : %v", err)
			return nil, err
		}

		message.Author, err = GetUserInfoFromID(db, user_id)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func FormatDate(messages []models.Message) []models.Message {
	for i := 0; i < len(messages); i++ {
		parts := strings.Split(messages[i].Created, " ")
		date := parts[0]
		parts = strings.Split(date, "-")

		if len(parts) != 3 {
			log.Printf("ERREUR : <getmessagelist.go> Erreur dans le format de la date sur le message n°%d. Doit être YYYY-MM-DD, est : %s", messages[i].MessageID, messages[i].Created)
			return nil
		}
		day := parts[2]
		month := parts[1]
		year := parts[0]

		messages[i].Created = day + "/" + month + "/" + year
	}

	return messages
}
