package getdata

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère le nombre de likes et de dislikes d'un sujet
func GetMessageLikesAndDislikes(db *sql.DB, postID int) (models.Message, error) {
	// Préparation de la requête sql
	sqlQuery := `SELECT IFNULL(likes, 0), IFNULL(dislikes, 0) FROM message WHERE post_id = ?`
	row, err := db.Query(sqlQuery, postID)
	if err != nil {
		return models.Message{}, err
	}

	var message models.Message

	message.MessageID = postID

	// Parcourt la base de données et récupère les informations pour rajouter tous les messages dans la slice
	err = row.Scan(&message.Likes, &message.Dislikes)
	if err != nil {
		log.Print("<getmessagelikes.go> Impossible de récupérer les likes et dislikes dans la base de données :", err)
		return models.Message{}, err
	}

	return message, nil
}
