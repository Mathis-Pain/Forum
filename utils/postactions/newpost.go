package postactions

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

func NewPost(userID, topicID int, message string, mode string) error {
	var newpost models.Message
	newpost.Author.ID = userID
	newpost.TopicID = topicID
	newpost.Content = message

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlQuery := `SELECT username, profilpic FROM user WHERE id = ?`
	row := db.QueryRow(sqlQuery, userID)

	err = row.Scan(&newpost.Author.Username, &newpost.Author.ProfilPic)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <newpost.go> : Impossible de récupérer les données de l'utilisateur %d : %v\n", userID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	err = addPostToDatabase(db, newpost, mode)

	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <newpost.go> Erreur lors de la création du nouveau message : ", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}

func addPostToDatabase(db *sql.DB, newpost models.Message, mode string) error {
	sqlUpdate := `INSERT INTO message (topic_id, content, user_id) VALUES(?, ?, ?)`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newpost.TopicID, newpost.Content, newpost.Author.ID)
	if err != nil {
		return err
	}

	topic, _ := getdata.GetTopicInfo(db, newpost.TopicID)
	if mode != "newtopic" {
		logMsg := fmt.Sprintf("USER : L'utilisateur %s a posté une réponse sur le sujet \"%s\" (ID : %d)\n", newpost.Author.Username, topic.Name, newpost.TopicID)
		logs.AddLogsToDatabase(logMsg)
	}

	return nil
}
