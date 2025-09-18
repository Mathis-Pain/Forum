package dbutils

import (
	"log"
	"os"
	"strings"
)

// ExtractSql lit le fichier schema.sql et construit la map des tables et colonnes
func ExtractSql(filePath string) map[string][]string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	sqlContent := string(data)
	schema := make(map[string][]string)

	parts := strings.Split(sqlContent, "CREATE TABLE")
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		openParen := strings.Index(part, "(")
		if openParen == -1 {
			continue
		}
		tableName := strings.TrimSpace(part[:openParen])

		closeParen := strings.LastIndex(part, ")")
		if closeParen == -1 {
			continue
		}
		columnsPart := part[openParen+1 : closeParen]

		// Nouvelle logique : gérer les parenthèses
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
				if parentheses == 0 {
					col := strings.TrimSpace(current.String())
					if col != "" {
						columns = append(columns, NormalizeColumn(col))
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
			columns = append(columns, NormalizeColumn(current.String()))
		}

		schema[tableName] = columns
	}

	return schema
}
