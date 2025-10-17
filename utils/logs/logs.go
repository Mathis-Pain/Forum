package logs

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
)

// AddLogsToDatabase ajoute un message de log dans la base de données de notifications
// Cette fonction traite le message pour en extraire le type et les informations pertinentes
// Format attendu du message : "TYPE : contenu du message"
func AddLogsToDatabase(message string) error {
	// Ouverture de la base de données SQLite dédiée aux notifications et logs
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	// Affichage du message dans la console pour le débogage
	log.Println(message)

	// Extraction du type de log depuis le message (première partie avant " : ")
	logType := retrieveLogType(message)

	// Nettoyage du message : suppression du préfixe de type et des retours à la ligne
	cutLogMessage(&message, logType)

	// Traitement différencié selon le type de log
	if logType == "REQUEST" {
		// Pour les requêtes, on extrait l'ID de l'utilisateur qui a fait la demande
		userID, _ := GetIDFromLog(message)
		sqlUpdate := `INSERT INTO logs (message, type, sender) VALUES (?, ?, ?)`
		_, err = db.Exec(sqlUpdate, message, logType, userID)
	} else {
		// Pour les autres types de logs, pas besoin d'émetteur
		sqlUpdate := `INSERT INTO logs (message, type) VALUES (?, ?)`
		_, err = db.Exec(sqlUpdate, message, logType)
	}

	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur dans l'ajout du log \"%s\" : %v\n", message, err)
		return err
	}

	return nil
}

// retrieveLogType extrait le type de log depuis le message
// Le type est toujours le premier mot du message (avant le premier espace)
// Exemple : "REQUEST : demande de modération" → retourne "REQUEST"
func retrieveLogType(message string) string {
	parts := strings.Split(message, " ")
	return parts[0]
}

// cutLogMessage nettoie le message en supprimant le préfixe de type et les retours à la ligne
// Cette fonction modifie directement le message passé en paramètre (pointeur)
// Exemple : "REQUEST : message\n" devient "message"
func cutLogMessage(message *string, logType string) {
	// Préparation des chaînes à supprimer
	prefix := logType + " : "
	suffix := "\n"

	// Suppression du préfixe (type de log)
	if strings.HasPrefix(*message, prefix) {
		*message, _ = strings.CutPrefix(*message, prefix)
	}

	// Suppression du suffixe (retour à la ligne)
	if strings.HasSuffix(*message, suffix) {
		*message, _ = strings.CutSuffix(*message, suffix)
	}
}

// DisplayLogs récupère et retourne tous les logs de la base de données
// Cette fonction construit également le lien vers le message concerné pour les logs liés à des messages
func DisplayLogs() ([]models.Log, error) {
	var logs []models.Log

	// Ouverture de la base de données de notifications
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <getuserprofil.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return nil, err
	}
	defer db.Close()

	// Récupération de tous les logs
	// IFNULL(sender, 0) : si sender est NULL, on retourne 0
	sqlQuery := `SELECT id, type, message, date, handled, IFNULL(sender, 0)sender FROM logs`
	rows, err := db.Query(sqlQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var log models.Log
	var handled int

	// Parcours de tous les résultats
	for rows.Next() {
		// Scan des données de chaque log
		if err := rows.Scan(&log.ID, &log.LogType, &log.LogMessage, &log.Date, &handled, &log.Requester); err != nil {
			return nil, err
		}

		// Conversion de l'état "traité" : 1 = true, 0 = false
		if handled == 1 {
			log.Handled = true
		} else {
			log.Handled = false
		}

		// Construction du lien vers le message concerné (si applicable)
		log.MessageLink, err = GetMessageLinkFromLog(log.LogMessage)
		if err != nil {
			return nil, err
		}

		// Ajout du log à la liste
		logs = append(logs, log)
	}

	return logs, nil
}

// DeleteLog supprime un log spécifique de la base de données
// Cette fonction est utilisée lorsqu'un administrateur supprime un log depuis le panneau admin
func DeleteLog(logID int) error {
	// Ouverture de la base de données de notifications
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	// Suppression du log par son ID
	sqlUpdate := `DELETE FROM logs WHERE ID = ?`

	_, err = db.Exec(sqlUpdate, logID)
	if err != nil {
		// En cas d'erreur, on enregistre un log de l'erreur elle-même
		logMsg := fmt.Sprintf("ERREUR : <logs.go> Erreur dans la suppression du log n°%d : %v", logID, err)
		AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}

// MarkAsHandled marque un log comme traité dans la base de données quand un administrateur y a répondu
// Recherche le log par son message et l'ID de l'émetteur
func MarkAsHandled(message string, receiverID int) error {
	// Ouverture de la base de données de notifications
	db, err := sql.Open("sqlite3", "./data/notifications/notifications.db")
	if err != nil {
		log.Printf("ERREUR : <logs.go> Erreur à l'ouverture de la base de données : %v\n", err)
		return err
	}
	defer db.Close()

	// Mise à jour du statut "handled" à 1 (traité)
	// On identifie le log par son message ET l'ID de l'émetteur pour éviter les confusions
	sqlUpdate := `UPDATE logs SET handled = 1 WHERE message = ? AND sender = ?`
	_, err = db.Exec(sqlUpdate, message, receiverID)
	if err != nil {
		// En cas d'erreur, on enregistre un log de l'erreur
		logMsg := fmt.Sprint("ERREUR : <logs.go> Erreur dans la mise à jour du log :", err)
		AddLogsToDatabase(logMsg)
		return err
	}
	return nil
}

// GetMessageLinkFromLog extrait l'ID du message depuis le log et construit le lien vers ce message
// Cette fonction permet de créer un lien cliquable vers le message concerné par le log
func GetMessageLinkFromLog(log string) (string, error) {
	// Extraction de l'ID du message depuis le texte du log
	MessageID, err := GetIDFromLog(log)
	if MessageID == 0 || err != nil {
		return "", nil
	}

	// Construction du lien vers le message
	messageLink, err := GetMessageLink(MessageID)
	if err != nil {
		return "", err
	}

	return messageLink, nil
}

// GetMessageLink construit l'URL du lien vers un message spécifique
// Format du lien : /topic/{topic_id}#{message_id}
func GetMessageLink(MessageID int) (string, error) {
	// Ouverture de la base de données principale du forum
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Récupération de l'ID du topic auquel appartient le message
	var topicID int
	sqlQuery := `SELECT topic_id FROM message WHERE id = ?`
	row := db.QueryRow(sqlQuery, MessageID)
	err = row.Scan(&topicID)

	// Si le message n'existe pas (par exemple s'il a été supprimé), on retourne une chaîne vide
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	// Construction du lien
	messageLink := fmt.Sprintf("/topic/%d#%d", topicID, MessageID)

	return messageLink, nil
}

// GetIDFromLog extrait un ID numérique depuis un message de log
// Cette fonction recherche un nombre précédé du symbole '°'
// Exemple : "Message signalé n°42 par l'utilisateur" → retourne 42
func GetIDFromLog(log string) (int, error) {
	var ID int
	foundID := false

	// Parcours de chaque caractère du message
	for _, char := range log {
		// Détection du marqueur '°' qui précède l'ID
		if char == '°' {
			foundID = true
		}

		// Une fois le marqueur trouvé, on lit tous les chiffres qui suivent
		if foundID {
			// Vérification que le caractère est un chiffre (0-9)
			if char >= '0' && char <= '9' {
				// Construction du nombre chiffre par chiffre
				ID *= 10
				ID += int(char) - 48
			}
		}
	}

	return ID, nil
}
