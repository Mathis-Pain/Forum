package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetCatDetails(db *sql.DB, catID int) (models.Category, error) {
	sqlQuery := `SELECT name, description FROM category WHERE id = ?`
	row := db.QueryRow(sqlQuery, catID)

	var categ models.Category
	err := row.Scan(&categ.Name, &categ.Description)
	if err != nil {
		return models.Category{}, err
	}

	categ.Topics, err = GetTopicList(db, catID)
	if err != nil {
		return models.Category{}, err
	}

	return categ, nil
}
