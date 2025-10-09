package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Fonction pour vérifier s'il y a eu une modification d'une catégorie par le formulaire
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

// Fonction pour supprimer tous les messages d'un sujet particulier
func AdminDeleteMessages(db *sql.DB, ID int) error {
	sqlUpdate := `DELETE FROM message WHERE topic_id = ?`

	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminback.go> Erreur dans la suppression du message n°%d : %v", ID, err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		logMsg := fmt.Sprint(err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}

func AskedToBeMod(db *sql.DB, ID int) error {
	sqlUpdate := `UPDATE user SET role_id = 5 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}
