package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
)

// Obtenir les infos du User depuis la session
func GetUserInfoFromSess(sessId string) (models.User, error) {

	var user models.User
	var username string

	// ** Récupération du username **
	currentSession, err := sessions.GetSession(sessId)
	if err != nil {
		return models.User{}, err
	}

	for key := range currentSession.Data {
		username = key
	}

	// ** Récupération des données du user **

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()

	sql := `SELECT id, username, email, profilpic FROM user where username = ?`
	row := db.QueryRow(sql, username)

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilPic)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserPosts(userId int) ([]models.LastPost, error) {
	var posts []models.LastPost
	var post models.LastPost
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<getuserposts.go> Could not open database: %v\n", err)
		return nil, err
	}
	defer db.Close()

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
		WHERE user_id = ?
    `

	rows, err := db.Query(sqlQuery, userId)
	if err != nil {
		log.Printf("<getuserposts.go> Could not execute the sql query: %v\n", err)
		return []models.LastPost{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&post.MessageID, &post.TopicID, &post.Content, &post.Created, &post.TopicName); err != nil {
			log.Printf("<getuserposts.go> Error scanning message row: %v\n", err)
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Printf("<getuserposts.go> Error during rows iteration: %v\n", err)
		return nil, err
	}

	return posts, nil

}
