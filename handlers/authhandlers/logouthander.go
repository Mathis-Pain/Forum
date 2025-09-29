package authhandlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
)

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Récupère la session depuis la requête
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Println("Erreur lors de la récupération de la session :", err)
		utils.InternalServError(w)
		return
	}

	// Ouvre la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Println("Erreur à l'ouverture de la base de données :", err)
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
