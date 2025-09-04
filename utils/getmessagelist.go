package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func GetMessageList(db *sql.DB, topicID int) ([]models.Message, error) {
	sqlQuery := `SELECT created_at,user_id,likes, dislikes, content FROM message WHERE topic_id = ?`
	rows, err := db.Query(sqlQuery, topicID)
	if err != nil {
		return nil, err
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Created, &message.Author, &message.Likes, &message.Dislikes, &message.Content); err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil, err
		}
		messages = append(messages, message)

	}

	return messages, nil
}
