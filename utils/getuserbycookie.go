package utils

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Récupère le pseudo et l'ID de l'utilisateur si un utilisateur est en ligne
func GetUserNameAndIDByCookie(r *http.Request, db *sql.DB) (string, int, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération du cookie : ", err)
		logs.AddLogsToDatabase(logMsg)
		return "", 0, err
	}
	session, err := sessions.GetSession(cookie.Value)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération de session : ", err)
		logs.AddLogsToDatabase(logMsg)
		return "", 0, err
	}

	sqlQuery := `SELECT username FROM user WHERE id = ?`
	row := db.QueryRow(sqlQuery, session.UserID)

	var username string

	err = row.Scan(&username)
	if err != nil {
		return "", 0, err
	}

	return username, session.UserID, nil
}
