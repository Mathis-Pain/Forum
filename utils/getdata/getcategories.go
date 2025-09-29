package getdata

import (
	"database/sql"
	"log"

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
	rows, err := db.Query("SELECT id, name, description FROM category")
	if err != nil {
		return []models.Category{}, err
	}
	defer rows.Close()

	// Parcourir les résultats
	for rows.Next() {
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
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

	for _, topic := range categ.Topics {
		if len(topic.Messages) == 0 {
			sqlUpdate := `DELETE FROM topic WHERE id = ?`
			stmt, err := db.Prepare(sqlUpdate)
			if err != nil {
				log.Print("<admincatsandtopics.go> Erreur dans la suppression du sujet", err)
				return categ, err
			}
			defer stmt.Close()
			_, err = stmt.Exec(topic.TopicID)
			if err != nil {
				return categ, err
			}
		}

	}

	if err != nil {
		return models.Category{}, err
	}

	return categ, nil
}
