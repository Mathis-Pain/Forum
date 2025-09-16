package utils

import (
	"database/sql"
	"log"
)

// Fonction pour ajouter les posts dans les table likes et dislikes de la base de données
func AddLikesAndDislikes(db *sql.DB, postID, userID int, table string) error {
	var sqlUpdate string
	switch table {
	case "likes":
		sqlUpdate = `INSERT INTO likes (user_id, message_id) VALUES(?, ?)`
	case "dislikes":
		sqlUpdate = `INSERT INTO dislikes (user_id, message_id) VALUES(?, ?)`
	}
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	result, err := stmt.Exec(userID, postID)
	if err != nil {
		log.Print(err)
		return err
	}
	n, _ := result.RowsAffected()
	log.Printf("<updatelikes.go> %d lignes ont été rajoutées à la table %s (message ID : %d, par l'utilisateur n°%d)", n, table, postID, userID)

	return nil
}

// Fonction pour supprimer les posts dans les tables likes et dislikes de la base de données
func RemoveLikesAndDislikes(db *sql.DB, postID, userID int, table string) error {
	var sqlUpdate string
	switch table {
	case "likes":
		sqlUpdate = `DELETE FROM likes WHERE user_id = ? AND message_id = ?`
	case "dislikes":
		sqlUpdate = `DELETE FROM dislikes WHERE user_id = ? AND message_id = ?`
	}
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	result, err := stmt.Exec(userID, postID)
	if err != nil {
		log.Print(err)
		return err
	}
	n, _ := result.RowsAffected()
	log.Printf("<updatelikes.go> %d lignes ont été supprimées de la table %s (message ID : %d, par l'utilisateur n°%d)", n, table, postID, userID)

	return nil
}

// Met à jour le nombre de likes et de dislikes dans la table message pour le post
func UpdateLikesAndDislikes(db *sql.DB, postID, userID, likes, dislikes int, table string) error {
	sqlUpdate := `UPDATE message SET dislikes = ?, likes = ? WHERE message_id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(dislikes, likes, postID)
	if err != nil {
		log.Print(err)
		return err
	}
	n, _ := result.RowsAffected()
	log.Printf("<updatelikes.go> La table message a été mise à jour sur %d lignes (message ID : %d)", n, postID)

	return nil
}
