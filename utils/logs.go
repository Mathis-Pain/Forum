package utils

import (
	"database/sql"
	"log"
	"strings"
)

func AddLogsToDatabase(message string) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	log.Println(message)
	logType := retrieveLogType(message)
	cutLogMessage(&message, logType)

	sqlUpdate := `INSERT INTO logs (message, type) VALUES (?, ?)`
	_, err = db.Exec(sqlUpdate, message, logType)
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur dans l'ajout du log \"%s\" : %v\n", message, err)
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
	suffix := "\n"

	if strings.HasPrefix(*message, prefix) {
		*message, _ = strings.CutPrefix(*message, prefix)
	}

	if strings.HasSuffix(*message, suffix) {
		*message, _ = strings.CutSuffix(*message, suffix)
	}
}
