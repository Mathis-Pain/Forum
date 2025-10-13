package adminhandlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// MARK: Pseudo à partir de l'ID
func ConvertIDtoUsername(stringID string) (string, error) {
	ID, err := strconv.Atoi(stringID)

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	user, err := getdata.GetUserInfoFromID(db, ID)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

// MARK: Vérifie si la catégorie a été modifiée
func AdminIsCatModified(r *http.Request, categories []models.Category) (models.Category, bool, error) {
	// Récupère les données du formulaire
	name := r.FormValue("name")
	description := r.FormValue("description")
	stringID := r.FormValue("catID")

	// Convertit l'ID récupéré de la catégorie en int pour les comparaisons
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return models.Category{}, false, err
	}

	// Récupère la catégorie concernée par la modification
	var categ models.Category

	for _, current := range categories {
		if current.ID == ID {
			categ = current
			break
		}
	}

	// Compare le nom et la description. Si les deux sont les mêmes qu'avant, c'est que la catégorie n'a pas été modifiée
	if categ.Name == name && categ.Description == description {
		return categ, false, nil
	}

	return categ, true, nil

}
