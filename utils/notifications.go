package utils

import (
	"database/sql"
	"errors"
	"log"
)

func AddNotificationToDatabase(notifType string, userID int, message string) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	typeInt := convertType(notifType)

	if typeInt == 0 {
		return errors.New("type de notification invalide")
	}

	sqlUpdate := `INSERT INTO notifications (receiver_id, type, message) VALUES (?, ?, ?)`
	_, err = db.Exec(sqlUpdate, userID, typeInt, message)
	if err != nil {
		log.Printf("ERREUR : <notifications.go> Erreur dans l'ajout de la notification \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}

func convertType(notifType string) int {
	switch notifType {
	case "ADMIN":
		return 1
	case "REQUEST":
		return 2
	case "ANSWER":
		return 3
	case "MESSAGE":
		return 4
	case "INTERACTION":
		return 5
	}

	return 0
}
