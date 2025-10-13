package getdata

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/logs"
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

func GetMessageAuthor(db *sql.DB, postID int) (int, error) {
	var authorID int

	sqlQuery := `SELECT user_id FROM message WHERE id = ?`
	row := db.QueryRow(sqlQuery, postID)
	err := row.Scan(&authorID)
	if err != nil {
		return 0, err
	}

	return authorID, nil
}

func FormatDateAllMessages(messages []models.Message) []models.Message {
	for i := 0; i < len(messages); i++ {
		messages[i].Created = FormatDate(messages[i].Created)
		if strings.Contains(messages[i].Created, "ERREUR") {
			logMsg := messages[i].Created + " (message n°%d)"
			logs.AddLogsToDatabase(logMsg)
			return nil
		}
	}
	return messages
}

func FormatDate(dateToConvert string) string {
	parts := strings.Split(dateToConvert, " ")
	date := parts[0]
	parts = strings.Split(date, "-")

	if len(parts) != 3 {
		logMsg := fmt.Sprint("ERREUR : <getmessagelist.go> Erreur dans le format de la date sur le message. Doit être YYYY-MM-DD, est : ", dateToConvert)
		return logMsg
	}
	day := parts[2]
	month := parts[1]
	year := parts[0]

	date = day + "/" + month + "/" + year

	return date
}
