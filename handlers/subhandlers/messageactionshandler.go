package subhandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
)

func MessageActionsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du topic
	stringID := r.FormValue("topicID")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<messageactionshandler> Erreur dans la récupération de l'ID du sujet : ", err)
		utils.InternalServError(w)
		return
	}

	// Récupère l'ID du message
	stringID = r.FormValue("postID")
	postID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<messageactionshandler> Erreur dans la récupération de l'ID du message : ", err)
		utils.InternalServError(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<messageactionshandler> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	switch r.FormValue("action") {
	case "delete":
		log.Print("Suppression du message, ID : ", postID)
		err := admin.AdminDeleteMessage(topicID, postID, db)
		if err != nil {
			log.Print("<messageactionshandler> Erreur dans la suppression du message : ", err)
			utils.InternalServError(w)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case "warn":
		log.Print("Signalement du message, ID : ", postID)
		err := admin.ModSignalMessage(postID, db)
		if err != nil {
			log.Print("<messageactionshandler> Erreur dans le signalement du message : ", err)
			utils.InternalServError(w)
			return
		}
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "cancel":
		log.Print("Annulation du signalement sur message, ID : ", postID)
		err := admin.AdminCancelSignal(postID, db)
		if err != nil {
			log.Print("<messageactionshandler> Erreur dans l'annulation du signalement' : ", err)
			utils.InternalServError(w)
			return
		}
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "edit":
		err := admin.EditExistingMessage(postID, db)
		log.Print("Modification du message, ID : ", postID)
		if err != nil {
			log.Print("<messageactionshandler> Erreur dans la modification du message' : ", err)
			utils.InternalServError(w)
			return
		}
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		utils.StatusBadRequest(w)
		return
	}
}
