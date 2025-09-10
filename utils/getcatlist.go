package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetCatList() ([]models.Category, error) {
	var category models.Category
	var categories []models.Category

	// --- Ouverture de la db ---

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return []models.Category{}, err
	}
	defer db.Close()

	// Préparer la requête
	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		return []models.Category{}, err
	}
	defer rows.Close()

	// Parcourir les résultats
	for rows.Next() {
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return []models.Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
