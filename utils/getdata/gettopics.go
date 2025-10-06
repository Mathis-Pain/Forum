package getdata

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
			log.Printf("ERREUR : <gettopiclist.go> Erreur dans le parcours de la base de données : %v", err)
			return nil, err
		}

		// Récupère la liste des messages du sujet
		topic.Messages, err = GetMessageList(db, topic.TopicID)

		if err == sql.ErrNoRows {
			topic.Messages = []models.Message{}
			return topics, err
		} else if err != nil {
			return topics, err
		}

		topic.LastPost = len(topic.Messages) - 1
		if topic.LastPost < 0 {
			topic.LastPost = 0
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

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
