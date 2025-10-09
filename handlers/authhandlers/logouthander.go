package authhandlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère la session depuis la requête
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <logouthandler.go> Erreur lors de la récupération de la session :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Ouvre la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <logouthandler.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Supprime le cookie de session côté navigateur
	sessions.DeleteCookie(w, "session_id", false) // false si local, true si HTTPS

	// Supprime la session côté serveur/DB
	sessions.DeleteSession(session.ID)

	// Redirige vers la page d’accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
