package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func AddLogsToDatabase(message string) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	logType := retrieveLogType(message)
	cutLogMessage(&message, logType)

	sqlUpdate := `INSERT INTO logs (message, type) VALUES (?, ?)`
	_, err = db.Exec(sqlUpdate, message, logType)
	if err != nil {
		log.Printf("ERREUR : <notifications.go> Erreur dans l'ajout de la notification \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}

func retrieveLogType(message string) string {
	parts := strings.Split(message, " ")
	return parts[0]
}

func cutLogMessage(message *string, logType string) {
	prefix := logType + " : "

	fmt.Println(prefix)
	if strings.HasPrefix(*message, prefix) {
		*message, _ = strings.CutPrefix(*message, prefix)
	}
}
