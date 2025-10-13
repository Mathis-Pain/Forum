package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Gère les requêtes à l'administration
// Pour l'instant il y a juste un bouton "rejoindre la modération"
func RequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		action := r.FormValue("action")

		switch action {
		case "requestmod":
			db, err := sql.Open("sqlite3", "./data/forum.db")
			if err != nil {
				logMsg := fmt.Sprintf("ERREUR : <requesthandler.go> Could not open database : %v\n", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}
			defer db.Close()

			username := r.FormValue("username")
			userID, _ := strconv.Atoi(r.FormValue("id"))
			url := "/profil"

			err = AskedToBeMod(db, userID)

			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <requesthandler.go> Erreur dans l'envoi de la requête à l'administrateur, ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			logMsg := fmt.Sprintf("REQUEST : L'utilisateur %s (n°%d) a demandé à rejoindre la modération.", username, userID)
			logs.AddLogsToDatabase(logMsg)

			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	}

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
