package test

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// The BuildTestDatabase function now correctly uses the passed-in database connection.
func BuildTestDatabase(db *sql.DB) {
	// Check if the passed-in connection is valid
	if db == nil {
		log.Print("<testdatabase.go> Database connection is nil.")
		return
	}

	// This helper function centralizes the logic for executing a query and handling errors.
	// It's a much cleaner way to avoid repeating the same logic for every query.
	execAndLog := func(query, logMessage string) {
		result, err := db.Exec(query)

		// Check the error immediately and handle it before trying to use `result`.
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				fmt.Printf("<testdatabase.go> Skipped insertion due to unique constraint: %s\n", logMessage)
			} else {
				log.Printf("<testdatabase.go> Fatal error during insertion: %s - %v\n", logMessage, err)
				panic(err) // Panic on a fatal error so you can see the full stack trace.
			}
			return // Don't proceed if there was an error
		}

		rows, err := result.RowsAffected()
		if err != nil {
			log.Printf("<testdatabase.go> Error getting rows affected for %s: %v\n", logMessage, err)
			return
		}
		log.Printf("%d %s", rows, logMessage)
	}

	// Inserts
	execAndLog(`INSERT INTO role (name) VALUES ('ADMIN')`, "rôle ajouté dans la base de données")
	execAndLog(`INSERT INTO user (username, email, password, role_id) VALUES ('Zoé', 'mail@mail.com', '1234', '1')`, "utilisateur ajouté dans la base de données")
	execAndLog(`INSERT INTO category (name, description) VALUES ('Plop', 'Cette catégorie s''appelle Plop')`, "catégorie ajoutée dans la base de données")
	execAndLog(`INSERT INTO category (name, description) VALUES ('Deuxième catégorie', 'On commence à manquer d''idées pour les titres')`, "catégorie ajoutée dans la base de données")
	execAndLog(`INSERT INTO topic (name, category_id, user_id) VALUES ('Pourquoi plop ?', '1', '1')`, "sujet ajouté dans la base de données")
	execAndLog(`INSERT INTO topic (name, category_id, user_id) VALUES ('Pour les titres de sujet aussi on manque d''idées', '2', '1')`, "sujet ajouté dans la base de données")
	execAndLog(`INSERT INTO topic (name, category_id, user_id) VALUES ('On va faire un deuxième sujet pour la forme', '2', '1')`, "sujet ajouté dans la base de données")
	execAndLog(`INSERT INTO message (topic_id, content, user_id, likes, dislikes) VALUES ('1', 'Pourquoi cette catégorie s''appelle plop ?', '1', '0', '0')`, "message ajouté dans la base de données")
	execAndLog(`INSERT INTO message (topic_id, content, user_id, likes, dislikes) VALUES ('1', 'Est-ce un message caché ? Un complot ?', '1', '0', '0')`, "message ajouté dans la base de données")
	execAndLog(`INSERT INTO message (topic_id, content, user_id, likes, dislikes) VALUES ('2', 'Il faudrait voir à trouver un vrai thème pour ce forum un jour, quand même.', '1', '0', '0')`, "message ajouté dans la base de données")
	execAndLog(`INSERT INTO message (topic_id, content, user_id, likes, dislikes) VALUES ('3', 'Les tests tout fait, c''est quand même pratique !', '1', '0', '0')`, "message ajouté dans la base de données")
}
