package getdata

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetCatList() ([]models.Category, error) {
	var category models.Category
	var categories []models.Category

	// --- Ouverture de la db ---

	db, err := sql.Open("sqlite3", "./data/forum.db")
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

// Récupère le titre, la description et la liste des sujets d'une catégorie
func GetCatDetails(db *sql.DB, catID int) (models.Category, error) {
	// Création de la requête sql
	sqlQuery := `SELECT id, name, IFNULL(description, '') as description FROM category WHERE id = ?`
	row := db.QueryRow(sqlQuery, catID)

	// Parcourt la  base de données jusqu'à trouver la catégorie et récupérer les informations
	var categ models.Category
	err := row.Scan(&categ.ID, &categ.Name, &categ.Description)
	if err != nil {
		return models.Category{}, err
	}

	// Appelle la fonction pour récupérer la liste des sujets
	categ.Topics, err = GetTopicList(db, catID)
	if err != nil {
		return models.Category{}, err
	}

	return categ, nil
}
