package adminhandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// MARK: Modifier un sujet
func EditTopicHandler(r *http.Request, topics []models.Topic, admin string) error {
	// Récupère le nom du sujet, l'ID du sujet et celui de la catégorie
	name := r.FormValue("topicname")
	topicID := r.FormValue("topicID")
	stringID := r.FormValue("catID")

	// Convertit les deux ID au format int pour les comparaisons
	ID, err := strconv.Atoi(topicID)
	if err != nil {
		return err
	}

	catID, err := strconv.Atoi(stringID)
	if err != nil {
		return err
	}

	// Repère le sujet à modifier à partir de son ID
	var topic models.Topic
	for _, current := range topics {
		if current.TopicID == ID {
			topic = current
			break
		}
	}

	log.Println(topic.Messages[0].Author.ID)

	if name == topic.Name && catID == topic.CatID {
		return nil
	}

	logMsg := "ADMIN : "
	notif := ""

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	nameChanged := false
	if name != topic.Name && name != "" {
		nameChanged = true
	}

	topicMoved := false
	if topic.CatID != catID {
		topicMoved = true
	}

	// Si le nom a été modifié, change le nom
	if nameChanged {
		logMsg += fmt.Sprintf("Le sujet \"%s\" a été renommé en \"%s\"", name, topic.Name)
		notif += fmt.Sprintf("Votre sujet \"%s\" a été renommé en \"%s\"", topic.Name, name)
		topic.Name = name
	} else if topicMoved {
		logMsg += fmt.Sprintf("Le sujet \"%s\" a été", topic.Name)
		notif += fmt.Sprintf("Votre sujet \"%s\" a été", topic.Name)

	}

	if topicMoved {
		if nameChanged {
			logMsg += " et "
			notif += " et "
		}
		categ, err := getdata.GetCatDetails(db, catID)
		if err != nil {
			return err
		}
		logMsg += fmt.Sprintf("déplacé dans la catégorie \"%s\"", categ.Name)
		notif += fmt.Sprintf("déplacé dans la catégorie \"%s\"", categ.Name)
	}

	logMsg += fmt.Sprintf(" par %s.", admin)
	notif += fmt.Sprintf(" par %s.", admin)

	// Met à jour le sujet dans la base de données
	sqlUpdate := `UPDATE topic SET name = ?, category_id = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(topic.Name, catID, ID)
	if err != nil {
		return err
	}

	logs.AddLogsToDatabase(logMsg)
	logs.AddNotificationToDatabase("ADMIN", topic.Messages[0].Author.ID, 0, notif)

	return nil
}

// MARK: Supprimer un sujet
func DeleteTopicHandler(stringID string) error {
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <admincatsandtopics.go> Erreur dans la récupération du sujet à supprimer", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Supprime le sujet de la base de données
	sqlUpdate := `DELETE FROM topic WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	// Supprime tous les messages du sujet de la BDD
	err = AdminDeleteMessages(db, ID)
	if err != nil {
		return err
	}

	return nil
}

// MARK: Supprimer les messages
func AdminDeleteMessages(db *sql.DB, ID int) error {
	sqlUpdate := `DELETE FROM message WHERE topic_id = ?`

	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminback.go> Erreur dans la suppression du message n°%d : %v", ID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		logMsg := fmt.Sprint(err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}
