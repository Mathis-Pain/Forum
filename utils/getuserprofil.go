package utils

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
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

	for _, name := range currentSession.Data {
		username = name.(string)
	}

	// ** Récupération des données du user **

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()

	sql := `SELECT id, username, email, profilpic FROM user WHERE username = ?`
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
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}
	defer db.Close()

	sqlQuery := `
        SELECT
            m.id,
            m.topic_id,
            m.content,
            m.created_at,
            t.name
        FROM message m
        JOIN topic t ON m.topic_id = t.id
		WHERE m.user_id = ?
        ORDER BY m.created_at DESC
    `

	rows, err := db.Query(sqlQuery, userId)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de la requête SQL: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return []models.LastPost{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&post.MessageID, &post.TopicID, &post.Content, &post.Created, &post.TopicName); err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}
		post.Author, err = getdata.GetUserInfoFromID(db, userId)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de GetUserInfoFromID: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}

	return posts, nil

}

func GetUserLikes(userId int) ([]models.LastPost, error) {
	var posts []models.LastPost
	var post models.LastPost
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}
	defer db.Close()

	sqlQuery := `
        SELECT
            l.message_id,
            m.topic_id,
            m.content,
            m.created_at,
            t.name
        FROM like l
		JOIN message m ON l.message_id = m.id
        JOIN topic t ON m.topic_id = t.id
		WHERE l.user_id = ?
        ORDER BY m.created_at DESC
    `

	rows, err := db.Query(sqlQuery, userId)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de la requête SQL: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return []models.LastPost{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&post.MessageID, &post.TopicID, &post.Content, &post.Created, &post.TopicName); err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}
		post.Author, err = getdata.GetUserInfoFromID(db, userId)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de GetUserInfoFromID: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}

	return posts, nil
}

func GetUserDislikes(userId int) ([]models.LastPost, error) {
	var posts []models.LastPost
	var post models.LastPost
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}
	defer db.Close()

	sqlQuery := `
        SELECT
            d.message_id,
            m.topic_id,
            m.content,
            m.created_at,
            t.name
        FROM dislike d
		JOIN message m ON d.message_id = m.id
        JOIN topic t ON m.topic_id = t.id
		WHERE d.user_id = ?
        ORDER BY m.created_at DESC
    `

	rows, err := db.Query(sqlQuery, userId)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de la requête SQL: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return []models.LastPost{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&post.MessageID, &post.TopicID, &post.Content, &post.Created, &post.TopicName); err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}
		post.Author, err = getdata.GetUserInfoFromID(db, userId)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de GetUserInfoFromID: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}

	return posts, nil
}

func GetUserTopics(userId int) ([]models.LastPost, error) {
	// 1. Initialize variables
	var topics []models.LastPost // Renamed 'posts' to 'topics' for clarity

	// 2. Database connection and error handling
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}
	defer db.Close() // Ensure the connection is closed

	// 3. SQL Query to find topics started by the user
	sqlQuery := `
        SELECT 
            t.id AS topic_id,
            t.name AS topic_name,
            m.id AS message_id,
            m.content AS message_content,
            m.created_at AS message_created_at
        FROM topic t
        JOIN message m ON t.id = m.topic_id
        WHERE m.user_id = ?
        AND m.id = (
            SELECT MIN(id) 
            FROM message 
            WHERE topic_id = t.id
        )
        ORDER BY m.created_at DESC
    `

	rows, err := db.Query(sqlQuery, userId)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de la requête SQL: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topic models.LastPost

		if err := rows.Scan(
			&topic.TopicID,
			&topic.TopicName,
			&topic.MessageID,
			&topic.Content,
			&topic.Created,
		); err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données (Scan): %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}

		topic.Author, err = getdata.GetUserInfoFromID(db, userId)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans l'exécution de GetUserInfoFromID: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		logMsg := fmt.Sprintf("ERREUR : <getuserprofil.go> Erreur dans le parcours de la base de données (rows.Err): %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return nil, err
	}

	return topics, nil
}
