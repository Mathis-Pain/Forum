package subhandlers

import (
	"database/sql"
	"fmt"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Gère le bouton "rejoindre la modération"
func RequestMod(db *sql.DB, user models.User) {
	err := AskedToBeMod(db, user.ID)

	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <requesthandler.go> Erreur dans l'envoi de la requête à l'administrateur, ", err)
		logs.AddLogsToDatabase(logMsg)
		return
	}

	logMsg := fmt.Sprintf("REQUEST : L'utilisateur %s (n°%d) a demandé à rejoindre la modération.", user.Username, user.ID)
	logs.AddLogsToDatabase(logMsg)

}

// Met à jour le statut d'un membre ayant demandé à rejoindre la modération pour éviter le spam
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
