package logs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func AddNotificationToDatabase(notifType string, userID, postID int, message string) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	typeInt := ConvertNotifType(notifType)

	if typeInt == 0 {
		return errors.New("type de notification invalide")
	}

	sqlUpdate := `INSERT INTO notifications (receiver_id, type, message, post_id) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(sqlUpdate, userID, typeInt, message, postID)
	if err != nil {
		log.Printf("ERREUR : <notifications.go> Erreur dans l'ajout de la notification \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}

func DisplayNotifications(userID int) (models.Notifications, error) {
	var notifications models.Notifications

	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return models.Notifications{}, err
	}
	defer db.Close()

	sqlQuery := `SELECT id, type, message, read, post_id FROM notifications WHERE receiver_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		return models.Notifications{}, err
	}
	defer rows.Close()

	var notif models.Notif
	var read int
	var postID int
	// Parcourir les résultats
	for rows.Next() {
		if err := rows.Scan(&notif.ID, &notif.NotifType, &notif.NotifMessage, &read, &postID); err != nil {
			if err == sql.ErrNoRows {
				return models.Notifications{}, nil
			}
			return models.Notifications{}, err
		}

		if read == 1 {
			notif.Read = true
		} else {
			notif.Read = false
			notifications.NotRead += 1
		}

		notif.MessageLink, err = GetMessageLink(postID)
		if err != nil {
			return models.Notifications{}, err
		}

		notifications.Notifs = append(notifications.Notifs, notif)
	}

	return notifications, nil
}

func MarkAsRead(notifID int) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	sqlUpdate := `UPDATE notifications SET read = 1 WHERE ID = ?`
	_, err = db.Exec(sqlUpdate, notifID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <logs.go> Erreur dans la mise à jour du log :", err)
		AddLogsToDatabase(logMsg)
		return err
	}
	return nil
}

func DeleteNotif(notifID int) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	sqlUpdate := `DELETE FROM notifications WHERE ID = ?`

	_, err = db.Exec(sqlUpdate, notifID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <logs.go> Erreur dans la suppression du log n°%d : %v", notifID, err)
		AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}
