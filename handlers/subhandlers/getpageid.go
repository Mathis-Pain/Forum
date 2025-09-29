package subhandlers

import (
	"net/http"
	"strconv"
	"strings"
)

// Récupère l'ID de la page (pour les catégories et les sujets)
func GetPageID(r *http.Request) int {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		return 0
	}

	ID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}

	return ID
}
