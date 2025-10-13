package adminhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// MARK: Modifier une catégorie
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

// MARK: Supprimer une catégorie
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

// MARK: Ajouter une catégorie
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
