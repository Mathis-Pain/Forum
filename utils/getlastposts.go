package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func GetLastPosts() ([]models.LastPost, error) { // Return a slice of messages, each with its topic name
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<getlastposts.go> Could not open database: %v\n", err)
		return nil, err
	}
	defer db.Close()

	// Fetch the last 7 messages and their topic names in a single query
	// We join the message table with the topic table
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

	for rows.Next() {
		var mw models.LastPost
		if err := rows.Scan(&mw.MessageID, &mw.TopicID, &mw.Content, &mw.Created, &mw.Author, &mw.TopicName); err != nil {
			log.Printf("<getlastposts.go> Error scanning message row: %v\n", err)
			return nil, err
		}
		messagesWithTopics = append(messagesWithTopics, mw)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		log.Printf("<getlastposts.go> Error during rows iteration: %v\n", err)
		return nil, err
	}

	return messagesWithTopics, nil
}
