package dbutils

import "database/sql"

// récupère le SQL de création d'une table depuis SQLite dans forum.db
func GetTableSQL(db *sql.DB, tableName string) (string, error) {
	var sqlStmt string
	row := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name=?", tableName)
	err := row.Scan(&sqlStmt)
	if err != nil {
		return "", err
	}
	return sqlStmt, nil
}
