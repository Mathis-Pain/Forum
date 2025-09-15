package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère le nom, l'ID et la liste des messages pour chaque sujet présent dans la catégorie
func GetTopicList(db *sql.DB, catID int) ([]models.Topic, error) {
	// Préparation de la requête slq
	sqlQuery := `SELECT id, name FROM topic WHERE category_id = ?`
	rows, err := db.Query(sqlQuery, catID)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic

	// Parcourt le fichier et stocke chaque sujet dans la slice topics
	for rows.Next() {
		var topic models.Topic
		if err := rows.Scan(&topic.TopicID, &topic.Name); err != nil {
			log.Printf("<gettopiclist.go> Error scanning topic row: %v", err)
			return nil, err
		}
		// Récupère la liste des messages du sujet
		topic.Messages, err = GetMessageList(db, topic.TopicID)
		if err == sql.ErrNoRows {
			topic.Messages = []models.Message{}
		} else if err != nil {
			return topics, nil
		}
		topic.LastPost = len(topic.Messages) - 1
		if topic.LastPost < 0 {
			topic.LastPost = 0
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
