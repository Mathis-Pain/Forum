package builddb

import (
	"database/sql"
	"os"

	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// expectedSchema contiendra le schéma attendu
var expectedSchema map[string][]string

// ExtractSql lit le fichier schema.sql et construit la map des tables et colonnes
func ExtractSql(string) map[string][]string {
	data, err := os.ReadFile("forumdbschema.sql")
	if err != nil {
		log.Fatal(err)
	}

	sqlContent := string(data)
	schema := make(map[string][]string)

	// Séparer le contenu par "CREATE TABLE"
	parts := strings.Split(sqlContent, "CREATE TABLE")
	for _, part := range parts[1:] { // ignorer le premier split vide
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Extraire le nom de la table
		openParen := strings.Index(part, "(")
		if openParen == -1 {
			continue
		}
		tableName := strings.TrimSpace(part[:openParen])

		// Extraire les colonnes et contraintes
		closeParen := strings.LastIndex(part, ")")
		if closeParen == -1 {
			continue
		}
		columnsPart := part[openParen+1 : closeParen]

		// Nettoyer et séparer les lignes
		lines := strings.Split(columnsPart, ",")
		var columns []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				columns = append(columns, line)
			}
		}

		schema[tableName] = columns
	}

	return schema
}

// récupère le SQL de création d'une table depuis SQLite
func getTableSQL(db *sql.DB, tableName string) (string, error) {
	var sqlStmt string
	row := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name=?", tableName)
	err := row.Scan(&sqlStmt)
	if err != nil {
		return "", err
	}
	return sqlStmt, nil
}

// extrait les colonnes et contraintes d'un CREATE TABLE proprement
func extractColumns(createSQL string) []string {
	start := strings.Index(createSQL, "(")
	end := strings.LastIndex(createSQL, ")")
	if start == -1 || end == -1 || end <= start {
		return nil
	}

	columnsPart := createSQL[start+1 : end]
	var columns []string

	var current strings.Builder
	parentheses := 0

	for _, r := range columnsPart {
		switch r {
		case '(':
			parentheses++
			current.WriteRune(r)
		case ')':
			parentheses--
			current.WriteRune(r)
		case ',':
			// On sépare seulement si on n'est pas dans une parenthèse
			if parentheses == 0 {
				col := strings.TrimSpace(current.String())
				if col != "" {
					columns = append(columns, normalizeColumn(col))
				}
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	// Ajouter le dernier élément
	if current.Len() > 0 {
		col := strings.TrimSpace(current.String())
		if col != "" {
			columns = append(columns, normalizeColumn(col))
		}
	}

	return columns
}

// normalise une définition de colonne ou contrainte
func normalizeColumn(col string) string {
	col = strings.TrimSpace(col)
	col = strings.ToLower(col)              // insensible à la casse
	col = strings.ReplaceAll(col, "\"", "") // enlever guillemets
	col = strings.ReplaceAll(col, "`", "")  // enlever backticks
	return col
}

// compare deux listes de colonnes/contraintes indépendamment de l'ordre
func compareColumns(expected, actual []string) bool {
	expectedMap := make(map[string]bool)
	for _, c := range expected {
		expectedMap[normalizeColumn(c)] = true
	}
	for _, c := range actual {
		nc := normalizeColumn(c)
		if !expectedMap[nc] {
			return false
		}
		delete(expectedMap, nc)
	}
	return len(expectedMap) == 0
}
