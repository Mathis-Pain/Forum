package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Création de la base de données si elle n'existe pas encore
func InitDB() *sql.DB {
	dbPath := "./database/forum.db"

	// Vérifier si la DB existe déjà
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("Base non trouvée, création d'une nouvelle...")

		// Créer le fichier DB
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Erreur ouverture DB:", err)
		}
		defer db.Close()

		// Charger le schéma SQL
		schema, err := os.ReadFile("./database/schema.sql")
		if err != nil {
			log.Fatal("Erreur lecture schema.sql:", err)
		}

		// Exécuter le script SQL
		_, err = db.Exec(string(schema))
		if err != nil {
			log.Fatal("Erreur exécution schema.sql:", err)
		}

		fmt.Println("Base créée avec succès")
	}

	// Ouvrir la DB (existe déjà ou vient d'être créée)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Erreur ouverture DB:", err)
	}

	return db
}
