package builddb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Mathis-Pain/Forum/builddb/dbutils"
)

// InitDB initialise la base SQLite. Elle crée ou recrée la DB si nécessaire.
func InitDB(dbPath string, schemaPath string) (*sql.DB, error) {
	// log.Print("Analyse de la base de données à l'adresse : ", dbPath)

	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	recreateDB := false
	if !dbExists {
		// DB inexistante → création nécessaire
		recreateDB = true
	} else {
		// Vérifier le schéma existant
		expectedSchema := dbutils.ExtractSql(schemaPath)
		if err := CompareDB(dbPath, expectedSchema); err != nil {
			fmt.Println("Schéma différent :", err)

			// Faire le backup uniquement si le schéma est différent
			if err := BackupDB(dbPath); err != nil {
				fmt.Println("Backup non effectué:", err)
			}

			recreateDB = true
		}
	}

	if recreateDB {
		fmt.Println("Création d'une nouvelle base de données...")

		// Supprimer l'ancienne DB si elle existe
		if dbExists {
			if err := os.Remove(dbPath); err != nil {
				return nil, fmt.Errorf("Erreur dans la suppression de la DB existante : %w", err)
			}
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("Erreur à l'ouverture de la BDD : %w", err)
		}

		// Charger le schéma SQL
		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("Erreur à la lecture du schema sql : %w", err)
		}

		// Exécuter le script SQL
		if _, err := db.Exec(string(schema)); err != nil {
			db.Close()
			return nil, fmt.Errorf("Erreur à l'exécution du schema sql: %w", err)
		}

		fmt.Println("Base de données créée avec succès")
		DefaultDatabase(db, dbPath)
		return db, nil
	}

	// Ouvrir la DB existante (schéma correct)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture DB: %w", err)
	}

	if dbPath == "./data/forum.db" {
		fmt.Println("DB du forum correcte, aucun backup nécessaire")
	} else {
		fmt.Println("DB des notifications correcte, aucun backup nécessaire")
	}

	return db, nil
}
