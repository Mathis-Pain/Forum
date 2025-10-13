package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Récupère l'ID de la page (pour les catégories et les sujets)
func GetPageID(r *http.Request) int {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		return 0
	}

	ID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0
	}

	return ID
}

// Récupération de l'ID de la personne devant recevoir la notification
func AddSenderID(logMessage, username string) error {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messageactionshandler> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	user, err := getdata.GetUserInfoFromLogin(db, username)
	if err != nil {
		return err
	}

	sqlUpdate := `UPDATE logs SET sender = ? WHERE message = ?`

	_, err = db.Exec(sqlUpdate, user.ID, logMessage)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <logs.go> Erreur dans l'ajout de l'ID du modérateur :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}

	return nil
}
