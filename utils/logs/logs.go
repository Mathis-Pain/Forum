package logs

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
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

	if logType == "REQUEST" {
		userID, _ := GetIDFromLog(message)
		sqlUpdate := `INSERT INTO logs (message, type, sender) VALUES (?, ?, ?)`
		_, err = db.Exec(sqlUpdate, message, logType, userID)
	} else {
		sqlUpdate := `INSERT INTO logs (message, type) VALUES (?, ?)`
		_, err = db.Exec(sqlUpdate, message, logType)
	}

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

func DisplayLogs() ([]models.Log, error) {
	var logs []models.Log

	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return nil, err
	}
	defer db.Close()

	sqlQuery := `SELECT id, type, message, date, handled, IFNULL(sender, 0)sender FROM logs`
	rows, err := db.Query(sqlQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var log models.Log
	var handled int
	// Parcourir les résultats
	for rows.Next() {
		if err := rows.Scan(&log.ID, &log.LogType, &log.LogMessage, &log.Date, &handled, &log.Requester); err != nil {
			return nil, err
		}
		if handled == 1 {
			log.Handled = true
		} else {
			log.Handled = false
		}
		log.MessageLink, err = GetMessageLinkFromLog(log.LogMessage)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func DeleteLog(logID int) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	sqlUpdate := `DELETE FROM logs WHERE ID = ?`

	_, err = db.Exec(sqlUpdate, logID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <logs.go> Erreur dans la suppression du log n°%d : %v", logID, err)
		AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}

func MarkAsHandled(message string, receiverID int) error {
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	sqlUpdate := `UPDATE logs SET handled = 1 WHERE message = ? AND sender = ?`
	_, err = db.Exec(sqlUpdate, message, receiverID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <logs.go> Erreur dans la mise à jour du log :", err)
		AddLogsToDatabase(logMsg)
		return err
	}
	return nil
}

func GetMessageLinkFromLog(log string) (string, error) {
	MessageID, err := GetIDFromLog(log)
	if MessageID == 0 || err != nil {
		return "", nil
	}

	messageLink, err := GetMessageLink(MessageID)
	if err != nil {
		return "", err
	}

	return messageLink, nil
}

func GetMessageLink(MessageID int) (string, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	var topicID int
	sqlQuery := `SELECT topic_id FROM message WHERE id = ?`
	row := db.QueryRow(sqlQuery, MessageID)
	err = row.Scan(&topicID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	messageLink := fmt.Sprintf("/topic/%d#%d", topicID, MessageID)

	return messageLink, nil
}

func GetIDFromLog(log string) (int, error) {
	var ID int
	foundID := false

	for _, char := range log {
		if char == '°' {
			foundID = true
		}
		if foundID {
			if char >= '0' && char <= '9' {
				ID *= 10
				ID += int(char) - 48
			}
		}
	}

	return ID, nil
}
