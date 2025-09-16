package builddb

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/builddb/dbutils"
)

func CompareDB() error {
	dbPath := "forum.db"

	// Générer expectedSchema depuis le fichier schema.sql
	expectedSchema := dbutils.ExtractSql("schema.sql")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture DB pour comparaison: %w", err)
	}
	defer db.Close()
	var compareErr error

	for table, expectedCols := range expectedSchema {
		createSQL, err := dbutils.GetTableSQL(db, table)
		if err != nil {
			fmt.Printf("Erreur : table '%s' manquante !\n", table)
			compareErr = fmt.Errorf("table '%s' manquante", table)
			continue
		}
		actualCols := dbutils.ExtractColumns(createSQL)
		if !dbutils.CompareColumns(expectedCols, actualCols) {
			fmt.Printf("La table '%s' diffère du schéma attendu\nColonnes attendues : %v\nColonnes réelles : %v\n",
				table, expectedCols, actualCols)
			compareErr = fmt.Errorf("schéma différent pour la table '%s'", table)
		}
	}

	return compareErr
}
