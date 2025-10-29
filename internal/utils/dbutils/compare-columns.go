package dbutils

// compare deux listes de colonnes/contraintes indépendamment de l'ordre entre forum.db actuel et celui souhaité dans forumdbschema.sql
func CompareColumns(expected, actual []string) bool {
	expectedMap := make(map[string]bool)
	for _, c := range expected {
		expectedMap[NormalizeColumn(c)] = true
	}
	for _, c := range actual {
		nc := NormalizeColumn(c)
		if !expectedMap[nc] {
			return false
		}
		delete(expectedMap, nc)
	}
	return len(expectedMap) == 0
}
