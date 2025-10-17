package logs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

// AddNotificationToDatabase insère une nouvelle notification dans la base de données.
// Elle prend en paramètres le type de notification, l'ID de l'utilisateur receveur, l'ID du post associé et le message de la notification.
func AddNotificationToDatabase(notifType string, userID, postID int, message string) error {
	// Connexion à la base de données SQLite des notifications.
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	// Convertit le type de notification en chiffre pour la base de données
	typeInt := ConvertNotifType(notifType)

	if typeInt == 0 {
		return errors.New("type de notification invalide")
	}

	// Ajoute la notification à la base de données
	sqlUpdate := `INSERT INTO notifications (receiver_id, type, message, post_id) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(sqlUpdate, userID, typeInt, message, postID)
	if err != nil {
		log.Printf("ERREUR : <notifications.go> Erreur dans l'ajout de la notification \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}

// DisplayNotifications récupère et affiche toutes les notifications d'un utilisateur.
// C'est la petite boîte aux lettre dans le header
func DisplayNotifications(userID int) (models.Notifications, error) {
	var notifications models.Notifications

	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return models.Notifications{}, err
	}
	defer db.Close()

	// Récupère toutes les notifications qui ont été envoyées à l'utilisateur connecté
	sqlQuery := `SELECT id, type, message, read, post_id FROM notifications WHERE receiver_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		return models.Notifications{}, err
	}
	defer rows.Close()

	var notif models.Notif
	var read int
	var postID int

	for rows.Next() {
		if err := rows.Scan(&notif.ID, &notif.NotifType, &notif.NotifMessage, &read, &postID); err != nil {
			// Si aucune ligne n'est trouvée, c'est que l'utilisateur n'a pas de notification, donc pas d'erreur à renvoyer
			if err == sql.ErrNoRows {
				return models.Notifications{}, nil
			}
			// En cas d'autre erreur
			return models.Notifications{}, err
		}

		// Convertit 'read' (qui est au format 0 ou 1 dans la base de données) en booléen 'Read' (true/false) pour le modèle
		if read == 1 {
			notif.Read = true
		} else {
			notif.Read = false
			notifications.NotRead += 1
		}

		// Récupère le lien du message associé à la notification s'il y en a un
		notif.MessageLink, err = GetMessageLink(postID)
		if err != nil {
			return models.Notifications{}, err
		}

		// Ajoute la notification formatée à la liste des notifications.
		notifications.Notifs = append(notifications.Notifs, notif)
	}

	// Retourne la structure complète des notifications et nil pour succès.
	return notifications, nil
}

// MarkAsRead met à jour le statut d'une notification pour la marquer comme lue.
func MarkAsRead(notifID int) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	// Définit 'read' à 1 pour la notification (0 = non lue, 1 = lue)
	sqlUpdate := `UPDATE notifications SET read = 1 WHERE ID = ?`
	// 3. Exécute la mise à jour.
	_, err = db.Exec(sqlUpdate, notifID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <logs.go> Erreur dans la mise à jour du log :", err)
		AddLogsToDatabase(logMsg)
		return err
	}
	return nil
}

// DeleteNotif supprime une notification de la base de données.
// Il récupère simplement l'ID de la notification pour ensuite la supprimer.
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
