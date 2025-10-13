package subhandlers

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

// Fonction pour modifier une catégorie
func EditCatHandler(r *http.Request, categ models.Category, currentUser models.UserLoggedIn) error {
	// Récupère le nouveau nom et la nouvelle description dans le formulaire
	name := r.FormValue("name")
	description := r.FormValue("description")

	if name == categ.Name && description == categ.Description {
		return nil
	}

	logMsg := "ADMIN :"

	// Modifie le nom et la description s'ils ont été changés
	if name != "" && name != categ.Name {
		logMsg += fmt.Sprintf("La catégorie %s a été renommée en %s par %s.", categ.Name, name, currentUser.Username)
		categ.Name = name

	}
	if categ.Description != description {
		logMsg += fmt.Sprintf("La description de la catégorie \"%s\"  est maintenant : %s", categ.Name, description)
		categ.Description = description
	}

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Met à jour la catégorie dans la base de données
	sqlUpdate := `UPDATE category SET name = ?, description = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(categ.Name, categ.Description, categ.ID)
	if err != nil {
		return err
	}

	logs.AddLogsToDatabase(logMsg)
	return nil
}

// Fonction pour supprimer une catégorie
func DeleteCatHandler(stringID string) error {
	// Récupère l'ID (sous forme de string) et le convertit en int pour les comparaisons
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <admincatsandtopics.go> Erreur dans la récupération de la catégorie à supprimer")
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Supprime la catégorie dans la base de données
	sqlUpdate := `DELETE FROM category WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	// Récupère tous les sujets présents dans la catégorie
	topicList, err := getdata.GetTopicList(db, ID)
	if err != nil {
		return err
	}

	// Supprime de la BDD tous les messages de ces sujets
	for i := 0; i < len(topicList); i++ {
		err := AdminDeleteMessages(db, topicList[i].TopicID)
		if err != nil {
			logMsg := fmt.Sprintln("ERREUR : <admincatsandtopics.go> Erreur dans la suppression des messages")
			logs.AddLogsToDatabase(logMsg)
			return err
		}
	}

	// Supprime ensuite de la BDD les sujets de la catégorie
	sqlUpdate = `DELETE FROM topic WHERE category_id = ?`
	stmt, err = db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	// Confirme la suppression de la catégorie et de tout ce qu'elle contenait
	// logMsg := fmt.Sprint("ADMIN : Catégorie et sujets liés supprimés avec succès.")

	return nil
}

// Fonction pour ajouter une catégorie
func AddCatHandler(r *http.Request) error {
	// Récupère le nom et la description de la nouvelle catégorie
	name := r.FormValue("newcatname")
	description := r.FormValue("newcatdesc")

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Ajoute le nom et la description à la BDD
	sqlUpdate := `INSERT INTO category (name, description) VALUES(?, ?)`
	_, err = db.Exec(sqlUpdate, name, description)
	if err != nil {
		return err
	}

	return nil
}

// Fonction pour modifier un sujet (titre et catégorie)
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

// Fonction pour supprimer un sujet
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
