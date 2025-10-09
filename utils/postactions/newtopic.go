package postactions

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Fonction pour créer un nouveau sujet dans une catégoire
func CreateNewtopic(userID, catID int, topicName, message string) error {
	var newtopic models.Topic

	// Stocke le numéro de la catégorie et le nom du sujet dans la struct
	newtopic.CatID = catID
	newtopic.Name = topicName

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Ajoute les informations du sujet à la base de données (titre, créateur, catégorie)
	err = addTopicToDatabase(db, newtopic, userID)
	if err != nil {
		logMsg := fmt.Sprintln("<newtopic.go> Erreur dans la création d'un nouveau sujet :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	// Récupère l'ID du topic pour pouvoir ajouter le premier message dans la BDD des messages
	newtopic.TopicID, err = getTopicID(db, newtopic.Name, catID)
	if err != nil {
		logMsg := fmt.Sprintln("<newtopic.go> Erreur dans récupération de l'ID du sujet pour créer le premier message :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	// Ajout du premier message du sujet dans la BDD
	err = NewPost(userID, newtopic.TopicID, message, "newtopic")
	if err != nil {
		logMsg := fmt.Sprintln("<newtopic.go> Erreur dans l'ajout du message :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}

// Fonction pour ajouter le nouveau sujet dans la BDD
func addTopicToDatabase(db *sql.DB, newtopic models.Topic, userID int) error {
	sqlUpdate := `INSERT INTO topic (category_id, name, user_id) VALUES(?, ?, ?)`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newtopic.CatID, newtopic.Name, userID)
	if err != nil {
		return err
	}

	return nil
}

// Fonction pour récupérer l'ID du sujet que l'on vient d'ouvrir
func getTopicID(db *sql.DB, name string, catID int) (int, error) {
	var topicID int
	// Au cas où plusieurs sujets auraient le même titre, récupère le sujet :
	// Dans la bonne catégorie, et le plus récent sujet posté avec ce titre
	sqlQuery := `SELECT id FROM topic WHERE name = ? AND category_id = ? 
	ORDER BY created_at DESC LIMIT 1`
	row := db.QueryRow(sqlQuery, name, catID)
	err := row.Scan(&topicID)

	if err != nil {
		return 0, err
	}

	return topicID, nil
}
