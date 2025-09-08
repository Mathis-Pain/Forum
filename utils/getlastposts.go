package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère les 7 derniers messages postés sur le forum pour pouvoir les afficher sur la page d'accueil
// Le format LastPost est un format Message + titre du sujet
func GetLastPosts() ([]models.LastPost, error) {
	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<getlastposts.go> Could not open database: %v\n", err)
		return nil, err
	}
	defer db.Close()

	// Préparation de la requête sql :
	// - Joint la section "message" et la section "topic" pour récupérer le titre du sujet et les infos du message en une seule requête
	// - Récupère l'ID du message et celui du sujet, le contenu du message, la date de création, l'auteur du message et le titre du sujet
	// - Commence par le plus récent et s'arrête maximum à 7 messages
	sqlQuery := `
        SELECT
            m.id,
            m.topic_id,
            m.content,
            m.created_at,
            m.user_id,
            t.name
        FROM message m
        JOIN topic t ON m.topic_id = t.id
        ORDER BY m.created_at DESC
        LIMIT 7
    `

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Printf("<getlastposts.go> Error querying messages: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var messagesWithTopics []models.LastPost

	// Parcourt la base de données pour récupérer les informations
	for rows.Next() {
		var mw models.LastPost
		if err := rows.Scan(&mw.MessageID, &mw.TopicID, &mw.Content, &mw.Created, &mw.Author, &mw.TopicName); err != nil {
			log.Printf("<getlastposts.go> Error scanning message row: %v\n", err)
			return nil, err
		}
		messagesWithTopics = append(messagesWithTopics, mw)
	}

	if err = rows.Err(); err != nil {
		log.Printf("<getlastposts.go> Error during rows iteration: %v\n", err)
		return nil, err
	}

	return messagesWithTopics, nil
}
