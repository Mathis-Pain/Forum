package subhandlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// Fonction pour modifier une catégorie
func EditCatHandler(r *http.Request, categ models.Category) error {
	// Récupère le nouveau nom et la nouvelle description dans le formulaire
	name := r.FormValue("name")
	description := r.FormValue("description")

	// Modifie le nom et la description s'ils ont été changés
	if name != "" {
		categ.Name = name
	}
	if description != "" {
		categ.Description = description
	}

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("ERREUR : <admincatsandtopics.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Met à jour la catégorie dans la base de données
	sqlUpdate := `UPDATE category SET name = ?, description = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(categ.Name, categ.Description, categ.ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// Fonction pour supprimer une catégorie
func DeleteCatHandler(stringID string) error {
	// Récupère l'ID (sous forme de string) et le convertit en int pour les comparaisons
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Println("ERREUR : <admincatsandtopics.go> Erreur dans la récupération de la catégorie à supprimer")
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Println("ERREUR : <admincatsandtopics.go> Erreur à l'ouverture de la base de données :")
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
			log.Println("ERREUR : <admincatsandtopics.go> Erreur dans la suppression des messages")
			return err
		}
	}

	// Supprime ensuite de la BDD les sujets de la catégorie
	sqlUpdate = `DELETE FROM topic WHERE category_id = ?`
	stmt, err = db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		log.Print(err)
		return err
	}

	// Confirme la suppression de la catégorie et de tout ce qu'elle contenait
	// log.Print("ADMIN : Catégorie et sujets liés supprimés avec succès.")

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
func EditTopicHandler(r *http.Request, topics []models.Topic) error {
	// Récupère le nom du sujet, l'ID du sujet et celui de la catégorie
	name := r.FormValue("topicname")
	topicID := r.FormValue("topicID")
	stringID := r.FormValue("catID")

	// Convertit les deux ID au format int pour les comparaisons
	ID, err := strconv.Atoi(topicID)
	if err != nil {
		return nil
	}

	catID, err := strconv.Atoi(stringID)
	if err != nil {
		return nil
	}

	// Repère le sujet à modifier à partir de son ID
	var topic models.Topic
	for _, current := range topics {
		if current.TopicID == ID {
			topic = current
			break
		}
	}

	// Si le nom a été modifié, change le nom
	if name != "" {
		topic.Name = name
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("ERREUR : <admincatsandtopics.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

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

	return nil
}

// Fonction pour supprimer un sujet
func DeleteTopicHandler(stringID string) error {
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("ERREUR : <admincatsandtopics.go> Erreur dans la récupération du sujet à supprimer", err)
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("ERREUR : <admincatsandtopics.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Supprime le sujet de la base de données
	sqlUpdate := `DELETE FROM topic WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print("ERREUR : <admincatsandtopics.go> Erreur dans la suppression du sujet", err)
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
		log.Print("ERREUR : <admincatsandtopics.go> Erreur dans la suppression des messages", err)
		return err
	}

	// Confirmation des modifications
	// log.Print("Sujets et messages supprimés avec succès.")

	return nil
}
