package builddb

import (
	"database/sql"
	"fmt"
	"os"
)

// InitDB initialise la base SQLite. Elle crée ou recrée la DB si nécessaire.
func InitDB() (*sql.DB, error) {
	dbPath := "forum.db"
	schemaPath := "forumdbschema.sql"

	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	// Faire un backup si la DB existe
	if dbExists {
		if err := BackupDB(dbPath); err != nil {
			fmt.Println("Backup non effectué:", err)
		}
	}

	recreateDB := false
	if !dbExists {
		recreateDB = true
	} else {
		// Vérifier le schéma existant
		if err := CompareDB(); err != nil {
			fmt.Println("Schéma différent :", err)
			recreateDB = true
		}
	}

	if recreateDB {
		fmt.Println("Création d'une nouvelle base...")

		// Supprimer l'ancienne DB si elle existe
		if dbExists {
			if err := os.Remove(dbPath); err != nil {
				return nil, fmt.Errorf("erreur suppression DB existante: %w", err)
			}
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("erreur ouverture DB: %w", err)
		}

		// Charger le schéma SQL
		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("erreur lecture schema.sql: %w", err)
		}

		// Exécuter le script SQL
		if _, err := db.Exec(string(schema)); err != nil {
			db.Close()
			return nil, fmt.Errorf("erreur exécution schema.sql: %w", err)
		}

		fmt.Println("Base créée avec succès")
		return db, nil
	}

	// Ouvrir la DB existante
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture DB: %w", err)
	}

	return db, nil
}
