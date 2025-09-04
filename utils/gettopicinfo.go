package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetTopicInfo(db *sql.DB, topicID int) (models.Topic, error) {
	sqlQuery := `SELECT name FROM topics WHERE id = ?`
	row := db.QueryRow(sqlQuery, topicID)

	var topic models.Topic
	err := row.Scan(&topic.Name)
	if err != nil {
		return models.Topic{}, err
	}

	topic.Messages, err = GetMessageList(db, topicID)
	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}
