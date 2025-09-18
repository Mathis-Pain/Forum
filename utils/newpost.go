package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func NewPost(userID, topicID int, message string) error {
	var newpost models.Message
	newpost.Author.ID = userID
	newpost.TopicID = topicID
	newpost.Content = message

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println("<newpost.go> Could not open database : ", err)
		return err
	}
	defer db.Close()

	sqlQuery := `SELECT username, profilpic FROM user WHERE id = ?`
	row := db.QueryRow(sqlQuery, userID)

	err = row.Scan(&newpost.Author.Username, &newpost.Author.ProfilPic)
	if err != nil {
		log.Printf("<newpost.go> : Impossible de récupérer les données de l'utilisateur %d : %v\n", userID, err)
		return err
	}
	err = addPostToDatabase(db, newpost)

	if err != nil {
		log.Println("<newpost.go> Erreur lors de la création du nouveau message : ", err)
		return err
	}

	return nil
}

func addPostToDatabase(db *sql.DB, newpost models.Message) error {
	sqlUpdate := `INSERT INTO message (topic_id, content, user_id) VALUES(?, ?, ?)`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(newpost.TopicID, newpost.Content, newpost.Author.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	log.Printf("<newpost.go> %d nouveau message ajouté au topic %d par l'utilisateur n°%d)", n, newpost.TopicID, newpost.Author.ID)

	return nil
}
