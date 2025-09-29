package admin

import (
	"database/sql"
	"fmt"
	"log"
)

func AdminDeleteMessage(topicID, postID int, db *sql.DB) error {
	// Supprime le message de la base de données
	sqlUpdate := `DELETE FROM message WHERE id = ?`
	_, err := db.Exec(sqlUpdate, postID)
	if err != nil {
		log.Printf("<adminmessage.go> Erreur dans la suppression du message %d : %v", postID, err)
		return err
	}

	logMsg := fmt.Sprintf("Suppression du message %d réussie.", postID)

	var topic string
	sqlQuery := `SELECT id FROM message WHERE topic_id = ?`
	row := db.QueryRow(sqlQuery, topicID)
	err = row.Scan(&topic)
	if err != nil {
		if err == sql.ErrNoRows {
			// S'il n'y a plus de messages dans le sujet, supprime le sujet
			sqlUpdate := `DELETE FROM topic WHERE id = ?`
			_, err := db.Exec(sqlUpdate, topicID)
			if err != nil {
				log.Printf("<adminmessage.go> Erreur dans la suppression du message %d : %v", postID, err)
				return err
			}
			logMsg += fmt.Sprintf(" Le sujet %d ne contient plus aucun message et a été supprimé.", postID)
		} else {
			return err
		}
	}

	log.Print(logMsg)

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

func EditExistingMessage(postID int, db *sql.DB) error {

	return nil
}
