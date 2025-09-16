package dbutils

import "strings"

// normalise une définition de colonne ou contrainte
func NormalizeColumn(col string) string {
	col = strings.TrimSpace(col)
	col = strings.ToLower(col)              // insensible à la casse
	col = strings.ReplaceAll(col, "\"", "") // enlever guillemets
	col = strings.ReplaceAll(col, "`", "")  // enlever backticks
	return col
}
