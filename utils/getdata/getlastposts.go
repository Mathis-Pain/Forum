package getdata

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Fonction pour afficher un "preview" d'un message dans l'affichage des catégories
func Preview(s string, length int) string {
	// Si le message fait plus de la longueur du preview (par exemple 300 caractères), coupe le message et ajoute "..."
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}

// Récupère les 7 derniers messages postés sur le forum pour pouvoir les afficher sur la page d'accueil
// Le format LastPost est un format Message + titre du sujet
func GetLastPosts() ([]models.LastPost, error) {
	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("ERREUR : <getlastposts.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return nil, err
	}
	defer db.Close()

	// Préparation de la requête sql :
	// - Joint la section "message" et la section "topic" pour récupérer le titre du sujet et les infos du message en une seule requête
	// - Récupère l'ID du message et celui du sujet, le contenu du message, la date de création, l'auteur du message et le titre du sujet
	// - Commence par le plus récent et s'arrête maximum à 7 messages
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
    `

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Printf("ERREUR : <getlastposts.go> Erreur dans la récupération des messages : %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var messagesWithTopics []models.LastPost

	// Parcourt la base de données pour récupérer les informations
	for rows.Next() {
		var mw models.LastPost
		var user_id int
		if err := rows.Scan(&mw.MessageID, &mw.TopicID, &mw.Content, &mw.Created, &user_id, &mw.TopicName); err != nil {
			log.Printf("ERREUR : <getlastposts.go> Erreur dans le parcours de la base de données : %v\n", err)
			return nil, err
		}
		mw.Author, err = GetUserInfoFromID(db, user_id)
		toAppend := true
		for i := 0; i < len(messagesWithTopics); i++ {
			if messagesWithTopics[i].TopicName == mw.TopicName && messagesWithTopics[i].Content == mw.Content {
				toAppend = false
			}
		}
		if toAppend {
			messagesWithTopics = append(messagesWithTopics, mw)
		}
		if len(messagesWithTopics) == 7 {
			break
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("ERREUR : <getlastposts.go> Erreur dans la conversation des données de la base de données : %v\n", err)
		return nil, err
	}

	return messagesWithTopics, nil
}

func LastMonthPost() ([]models.LastPost, int, error) {

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("ERREUR : <getlastposts.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return nil, 0, err
	}
	defer db.Close()

	sqlQuery := `
        SELECT
            m.id,
            m.topic_id,
            m.content,
            m.created_at,
            m.user_id,
			IFNULL(m.likes, 0),
			IFNULL(m.dislikes, 0),
			u.username,
            t.name
        FROM message m
        JOIN topic t ON m.topic_id = t.id
		JOIN user u ON m.user_id = u.id
		WHERE m.created_at >= DATETIME('now', '-30 days')
        ORDER BY m.created_at DESC
    `
	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Print("ERREUR : <getlastposts.go> Erreur dans la récupération des derniers messages :", err)
		return nil, 0, err
	}
	defer rows.Close()

	var lastMontPosts []models.LastPost
	for rows.Next() {
		var currentPost models.LastPost

		err := rows.Scan(&currentPost.MessageID, &currentPost.TopicID, &currentPost.Content, &currentPost.Created, &currentPost.Author.ID, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Author.Username, &currentPost.TopicName)
		if err != nil {
			log.Print("ERREUR : <getlastposts.go> Erreur dans le parcours de la base de données :", err)
			return nil, 0, err
		}
		lastMontPosts = append(lastMontPosts, currentPost)
	}

	return lastMontPosts, len(lastMontPosts), nil
}
