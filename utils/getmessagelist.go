package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère la liste complète des messages (date de création, auteur, contenu) pour un sujet
func GetMessageList(db *sql.DB, topicID int) ([]models.Message, error) {
	// Préparation de la requête sql
	sqlQuery := `SELECT created_at,user_id, content FROM message WHERE topic_id = ?`
	rows, err := db.Query(sqlQuery, topicID)
	if err != nil {
		return nil, err
	}

	var messages []models.Message

	// Parcourt la base de données et récupère les informations pour rajouter tous les messages dans la slice
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Created, &message.Author, &message.Content); err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil, err
		}
		messages = append(messages, message)

	}

	return messages, nil
}
