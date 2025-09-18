package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// Fonction pour créer un nouveau sujet dans une catégoire
func CreateNewtopic(userID, catID int, topicName, message string) error {
	var newtopic models.Topic

	// Stocke le numéro de la catégorie et le nom du sujet dans la struct
	newtopic.CatID = catID
	newtopic.Name = topicName

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Println("<newtopic.go> Could not open database : ", err)
		return err
	}
	defer db.Close()

	// Ajoute les informations du sujet à la base de données (titre, créateur, catégorie)
	err = addTopicToDatabase(db, newtopic)
	if err != nil {
		log.Println("<newtopic.go> Erreur dans la création d'un nouveau sujet :", err)
		return err
	}

	// Récupère l'ID du topic pour pouvoir ajouter le premier message dans la BDD des messages
	newtopic.TopicID, err = getTopicID(db, newtopic.Name)
	if err != nil {
		log.Println("<newtopic.go> Erreur dans la création d'un nouveau sujet :", err)
		return err
	}

	// Ajout du premier message du sujet dans la BDD
	err = NewPost(userID, newtopic.TopicID, message)
	if err != nil {
		log.Println("<newtopic.go> Erreur dans la création d'un nouveau sujet :", err)
		return err
	}

	return nil
}

// Fonction pour ajouter le nouveau sujet dans la BDD
func addTopicToDatabase(db *sql.DB, newtopic models.Topic) error {
	sqlUpdate := `INSERT INTO topic (category_id, name, user_id) VALUES(?, ?, ?)`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(newtopic.CatID, newtopic.Name, newtopic.Messages[0].Author.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	log.Printf("<newtopic.go> %d nouveau sujet : %s. Ajouté à la catégorie %d par l'utilisateur n°%d)", n, newtopic.Name, newtopic.CatID, newtopic.Messages[0].Author.ID)

	return nil
}

// Fonction pour récupérer l'ID du sujet que l'on vient d'ouvrir
func getTopicID(db *sql.DB, name string) (int, error) {
	var topicID int
	sqlQuery := "SELECT id FROM topics WHERE name = ?"
	row := db.QueryRow(sqlQuery, name)

	err := row.Scan(&topicID)
	if err != nil {
		return 0, err
	}

	return topicID, nil
}
