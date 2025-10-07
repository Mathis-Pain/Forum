package subhandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

func MessageActionsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du topic
	stringID := r.FormValue("topicID")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du sujet : ", err)
		utils.InternalServError(w)
		return
	}

	// Récupère l'ID du message
	stringID = r.FormValue("postID")
	postID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du message : ", err)
		utils.InternalServError(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("ERREUR : <messageactionshandler> Erreur à l'ouverture de la base de données : %v\n", err)
		return
	}
	defer db.Close()

	username := r.FormValue("modname")

	switch r.FormValue("action") {
	case "delete":
		err := admin.AdminDeleteMessage(topicID, postID, db)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans la suppression du message : ", err)
			utils.InternalServError(w)
			return
		}
		log.Printf("ADMIN : Le message %d a été supprimé par %s.", postID, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case "warn":
		warnReason := r.FormValue("warning")
		err := admin.ModSignalMessage(postID, db)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans le signalement du message : ", err)
			utils.InternalServError(w)
			return
		}
		log.Printf("REQUEST : Le message n°%d a été signalé par %s pour la raison suivante : %s\n", postID, username, warnReason)
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "cancel":
		err := admin.AdminCancelSignal(postID, db)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans l'annulation du signalement' : ", err)
			utils.InternalServError(w)
			return
		}
		log.Printf("ADMIN : %s a annulé le signalement du message n°%d.\n", username, postID)
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "edit":
		content := r.FormValue("message-content")
		err := admin.EditExistingMessage(postID, db, content)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans la modification du message' : ", err)
			utils.InternalServError(w)
			return
		}
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		topic, _ := getdata.GetTopicInfo(db, topicID)
		log.Printf("USER : %s a modifié le contenu du message n°%d (sur %s)", username, postID, topic.Name)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "move":
		stringID := r.FormValue("topicindex")
		newtopicID, err := strconv.Atoi(stringID)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du nouveau sujet : ", err)
			utils.StatusBadRequest(w)
			return
		}

		err = admin.MoveMessage(newtopicID, postID, db)
		if err != nil {
			log.Print("ERREUR : <messageactionshandler> Erreur dans le déplacement du message' : ", err)
			utils.InternalServError(w)
			return
		}
		topic, _ := getdata.GetTopicInfo(db, newtopicID)
		log.Printf("ADMIN : Message n°%d déplacé par %s dans le sujet n°%d (%s)\n", postID, username, newtopicID, topic.Name)
		url := fmt.Sprintf("/topic/%d#%d", newtopicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		log.Print("ERREUR : <messageactionshandler.go> Requête invalide sur la page message : ", r.FormValue("action"))
		utils.StatusBadRequest(w)
		return
	}
}
