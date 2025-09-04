package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func GetTopicList(db *sql.DB, catID int) ([]models.Topic, error) {
	sqlQuery := `SELECT id, name FROM topic WHERE category_id = ?`
	rows, err := db.Query(sqlQuery, catID)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic

	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(&topic.TopicID, &topic.Name); err != nil {
			log.Printf("Error scanning topic row: %v", err)
			return nil, err
		}
		topic.Messages, err = GetMessageList(db, topic.TopicID)
		if err == sql.ErrNoRows {
			topic.Messages = []models.Message{}
		} else if err != nil {
			return topics, nil
		}
		topics = append(topics, topic)

	}

	return topics, nil
}
