package utils

// Fonction pour afficher un "preview" d'un message dans l'affichage des catégories
func Preview(s string, length int) string {
	// Si le message fait plus de la longueur du preview (par exemple 300 caractères), coupe le message et ajoute "..."
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}
