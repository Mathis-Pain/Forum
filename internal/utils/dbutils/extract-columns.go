package dbutils

import (
	"strings"
)

// extrait les colonnes et contraintes d'un CREATE TABLE proprement
func ExtractColumns(createSQL string) []string {
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
		col := strings.TrimSpace(current.String())
		if col != "" {
			columns = append(columns, NormalizeColumn(col))
		}
	}

	return columns
}
