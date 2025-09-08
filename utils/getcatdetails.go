package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

// Récupère le titre, la description et la liste des sujets d'une catégorie
func GetCatDetails(db *sql.DB, catID int) (models.Category, error) {
	// Création de la requête sql
	sqlQuery := `SELECT name, IFNULL(description, '') as description FROM category WHERE id = ?`
	row := db.QueryRow(sqlQuery, catID)

	// Parcourt la  base de données jusqu'à trouver la catégorie et récupérer les informations
	var categ models.Category
	err := row.Scan(&categ.Name, &categ.Description)
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
