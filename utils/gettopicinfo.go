package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère les informations d'un sujet à partir de son ID
func GetTopicInfo(db *sql.DB, topicID int) (models.Topic, error) {
	// Préparation de la requête sql
	sqlQuery := `SELECT name FROM topic WHERE id = ?`
	row := db.QueryRow(sqlQuery, topicID)

	var topic models.Topic
	// Récupération du titre du sujet
	err := row.Scan(&topic.Name)
	if err != nil {
		return models.Topic{}, err
	}
	// Récupération de la liste des messages
	topic.Messages, err = GetMessageList(db, topicID)
	if err != nil {
		return models.Topic{}, err
	}

	return topic, nil
}
