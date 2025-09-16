package builddb

import (
	"database/sql"
	"fmt"
)

func CompareDB() error {
	dbPath := "forum.db"
	expectedSchema := ExtractSql("schema.sql")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture DB pour comparaison: %w", err)
	}
	defer db.Close()

	missingTables := []string{}
	diffTables := []string{}

	for table, expectedCols := range expectedSchema {
		createSQL, err := getTableSQL(db, table)
		if err != nil {
			missingTables = append(missingTables, table)
			continue
		}
		actualCols := extractColumns(createSQL)
		if !compareColumns(expectedCols, actualCols) {
			diffTables = append(diffTables, table)
		}
	}

	if len(missingTables) > 0 || len(diffTables) > 0 {
		return fmt.Errorf(
			"tables manquantes: %v, tables différentes: %v",
			missingTables, diffTables,
		)
	}
	return nil
}
