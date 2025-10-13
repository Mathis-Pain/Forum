package postactions

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Fonction pour ajouter les posts dans les table likes et dislikes de la base de données
func AddLikesAndDislikes(db *sql.DB, postID, userID int, table string) error {
	var sqlUpdate string
	user, _ := getdata.GetUserInfoFromID(db, userID)

	switch table {
	case "like":
		sqlUpdate = `INSERT INTO like (user_id, message_id) VALUES(?, ?)`

	case "dislike":
		sqlUpdate = `INSERT INTO dislike (user_id, message_id) VALUES(?, ?)`
	}
	_, err := db.Exec(sqlUpdate, userID, postID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <updatelikes.go> Erreur dans l'ajout du like/dislike sur le post %d : %v\n", postID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	userToNotify, err := logs.GetUserToNotify(postID, db)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <updatelikes.go> Erreur dans la récupération de l'utilisateur à notifier : ", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	notif := fmt.Sprintf("Un utilisateur (%s) a réagi à votre ", user.Username)
	logs.AddNotificationToDatabase("INTERACTION", userToNotify, postID, notif)

	logMsg := fmt.Sprintf("USER : %s a ajouté un %s sur le post n°%d", user.Username, table, postID)
	logs.AddLogsToDatabase(logMsg)

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
		logMsg := fmt.Sprintf("ERREUR : <updatelikes.go> Erreur dans la suppression du like/dislike sur le post %d : %v", postID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	n, _ := result.RowsAffected()

	if n != 0 {
		user, _ := getdata.GetUserInfoFromID(db, userID)
		logMsg := fmt.Sprintf("USER : L'utilisateur %s a supprimé un %s sur le post n°%d", user.Username, table, postID)
		logs.AddLogsToDatabase(logMsg)
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
