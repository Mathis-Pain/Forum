package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func GetTopicList(db *sql.DB, catID int) ([]models.Topic, error) {
	sql := `SELECT topicid, name, author FROM topics WHERE catid = ?`
	rows, err := db.Query(sql, catID)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic

	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(&topic.TopicID, &topic.Name, &topic.Author); err != nil {
			log.Printf("Error scanning topic row: %v", err)
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
