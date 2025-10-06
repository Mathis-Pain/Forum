package postactions

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// Fonction pour ajouter les posts dans les table likes et dislikes de la base de données
func AddLikesAndDislikes(db *sql.DB, postID, userID int, table string) error {
	var sqlUpdate string
	switch table {
	case "like":
		sqlUpdate = `INSERT INTO like (user_id, message_id) VALUES(?, ?)`
	case "dislike":
		sqlUpdate = `INSERT INTO dislike (user_id, message_id) VALUES(?, ?)`
	}
	_, err := db.Exec(sqlUpdate, userID, postID)
	if err != nil {
		log.Printf("ERREUR : <updatelikes.go> Erreur dans l'ajout du like/dislike sur le post %d : %v\n", postID, err)
		return err
	}

	user, _ := getdata.GetUserInfoFromID(db, userID)
	log.Printf("USER : L'utilisateur %s a ajouté un %s sur le post n°%d", user.Username, table, postID)

	return nil
}

// Fonction pour supprimer les posts dans les tables likes et dislikes de la base de données
func RemoveLikesAndDislikes(db *sql.DB, postID, userID int, table string) error {
	var sqlUpdate string
	switch table {
	case "like":
		sqlUpdate = `DELETE FROM like WHERE user_id = ? AND message_id = ?`
	case "dislike":
		sqlUpdate = `DELETE FROM dislike WHERE user_id = ? AND message_id = ?`
	}
	result, err := db.Exec(sqlUpdate, userID, postID)
	if err != nil {
		log.Printf("ERREUR : <updatelikes.go> Erreur dans la suppression du like/dislike sur le post %d : %v", postID, err)
		return err
	}

	n, _ := result.RowsAffected()

	if n != 0 {
		user, _ := getdata.GetUserInfoFromID(db, userID)
		log.Printf("USER : L'utilisateur %s a supprimé un %s sur le post n°%d", user.Username, table, postID)
	}

	return nil
}

// Met à jour le nombre de likes et de dislikes dans la table message pour le post
func UpdateLikesAndDislikes(db *sql.DB, postID, userID, likes, dislikes int, table string) error {
	sqlUpdate := `UPDATE message SET dislikes = ?, likes = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(dislikes, likes, postID)
	if err != nil {
		return err
	}

	return nil
}
