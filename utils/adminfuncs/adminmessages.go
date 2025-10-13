package admin

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
	"github.com/Mathis-Pain/Forum/utils/postactions"
)

func AdminDeleteMessage(topicID, postID int, db *sql.DB, currentUser models.UserLoggedIn) error {

	authorID, err := getdata.GetMessageAuthor(db, postID)
	if err != nil {
		return err
	}
	topic, err := getdata.GetTopicInfo(db, topicID)
	if err != nil {
		return err
	}
	notif := fmt.Sprintf("Votre message sur le sujet \"%s\" a été supprimé.", topic.Name)

	// Supprime le message de la base de données
	sqlUpdate := `DELETE FROM message WHERE id = ?`
	_, err = db.Exec(sqlUpdate, postID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminmessage.go> Erreur dans la suppression du message %d : %v", postID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	logMsg := fmt.Sprintf("ADMIN : Suppression du message %d réussie.", postID)

	// Supprime tous les likes et dislikes liés à ce message de la base de données
	_, totalUsers, err := GetAllUsers()
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminmessage.go, GetStats> Erreur dans la récupération des utilisateurs", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	for i := 1; i <= totalUsers; i++ {
		postactions.RemoveLikesAndDislikes(db, postID, i, "like")
		postactions.RemoveLikesAndDislikes(db, postID, i, "dislike")
	}

	// Vérifie s'il reste encore des messages sur le sujet
	var scanner string
	sqlQuery := `SELECT id FROM message WHERE topic_id = ?`
	row := db.QueryRow(sqlQuery, topicID)
	err = row.Scan(&scanner)
	if err != nil {
		if err == sql.ErrNoRows {
			// S'il n'y a plus de messages dans le sujet, supprime le sujet
			sqlUpdate := `DELETE FROM topic WHERE id = ?`
			_, err := db.Exec(sqlUpdate, topicID)
			if err != nil {
				logMsg := fmt.Sprintf("ERREUR : <adminmessage.go> Erreur dans la suppression du message %d : %v", postID, err)
				logs.AddLogsToDatabase(logMsg)
				return err
			}
			logMsg += fmt.Sprintf(" Le sujet %d ne contient plus aucun message et a été supprimé.", topicID)
		} else {
			return err
		}
	}

	if authorID != currentUser.ID {
		logs.AddNotificationToDatabase("ADMIN", authorID, postID, notif)
	}
	logs.AddLogsToDatabase(logMsg)

	return nil
}

func AdminCancelSignal(postID int, db *sql.DB) error {
	sqlUpdate := `UPDATE message SET warning = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(0, postID)
	if err != nil {
		return err
	}

	return nil
}

func ModSignalMessage(postID int, db *sql.DB) error {
	sqlUpdate := `UPDATE message SET warning = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(1, postID)
	if err != nil {
		return err
	}

	return nil
}

func EditExistingMessage(postID int, db *sql.DB, content string) error {
	sqlUpdate := `UPDATE message SET content = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(content, postID)
	if err != nil {
		return err
	}

	return nil
}

func MoveMessage(topicID int, postID int, db *sql.DB) error {
	sqlUpdate := `UPDATE message SET topic_id = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(topicID, postID)
	if err != nil {
		return err
	}

	return nil
}
