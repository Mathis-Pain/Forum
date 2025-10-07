package utils

import (
	"database/sql"
	"log"
)

func AddLogsToDatabase(message string) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	sqlUpdate := `INSERT INTO logs (message) VALUES (?)`
	_, err = db.Exec(sqlUpdate, message)
	if err != nil {
		log.Printf("ERREUR : <notifications.go> Erreur dans l'ajout de la notification \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}
