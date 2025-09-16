package builddb

import (
	"database/sql"
	"fmt"
)

func CompareDB() error {
	dbPath := "forum.db"

	// Générer expectedSchema depuis le fichier schema.sql
	expectedSchema := ExtractSql("schema.sql")

	// Ouvrir la base SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture DB pour comparaison: %w", err)
	}
	defer db.Close()

	var compareErr error

	// Comparer chaque table
	for table, expectedCols := range expectedSchema {
		createSQL, err := getTableSQL(db, table)
		if err != nil {
			fmt.Printf("Erreur : table '%s' manquante !\n", table)
			compareErr = fmt.Errorf("table '%s' manquante", table)
			continue
		}
		actualCols := extractColumns(createSQL)
		if !compareColumns(expectedCols, actualCols) {
			fmt.Printf("La table '%s' diffère du schéma attendu\nColonnes attendues : %v\nColonnes réelles : %v\n",
				table, expectedCols, actualCols)
			compareErr = fmt.Errorf("schéma différent pour la table '%s'", table)
		}
	}

	return compareErr
}
