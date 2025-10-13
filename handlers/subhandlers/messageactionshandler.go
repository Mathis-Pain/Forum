package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

func MessageActionsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID du topic
	stringID := r.FormValue("topicID")
	topicID, err := strconv.Atoi(stringID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du sujet : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère l'ID du message
	stringID = r.FormValue("postID")
	postID, err := strconv.Atoi(stringID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du message : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messageactionshandler> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return
	}
	defer db.Close()

	username := r.FormValue("modname")

	_, _, currentUser, err := BuildHeader(r, w, db)

	switch r.FormValue("action") {
	case "delete":
		err := admin.AdminDeleteMessage(topicID, postID, db, currentUser)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la suppression du message : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}

		logMsg := fmt.Sprintf("ADMIN : Le message %d a été supprimé par %s.", postID, username)
		logs.AddLogsToDatabase(logMsg)

		if _, err := getdata.GetTopicInfo(db, topicID); err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
			http.Redirect(w, r, url, http.StatusSeeOther)
		}

		return
	case "warn":
		warnReason := r.FormValue("warning")
		err := admin.ModSignalMessage(postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans le signalement du message : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		logMsg := fmt.Sprintf("WARN : Le message n°%d a été signalé par %s pour la raison suivante : %s\n", postID, username, warnReason)
		logs.AddLogsToDatabase(logMsg)
		AddSenderID(logMsg, username)

		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "cancel":
		err := admin.AdminCancelSignal(postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans l'annulation du signalement' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		logMsg := fmt.Sprintf("ADMIN : %s a annulé le signalement du message n°%d.\n", username, postID)
		logs.AddLogsToDatabase(logMsg)
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "edit":
		content := r.FormValue("message-content")
		err := admin.EditExistingMessage(postID, db, content)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la modification du message' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		topic, _ := getdata.GetTopicInfo(db, topicID)
		logMsg := fmt.Sprintf("USER : %s a modifié le contenu du message n°%d (sur \"%s\")", username, postID, topic.Name)
		logs.AddLogsToDatabase(logMsg)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	case "move":
		stringID := r.FormValue("topicindex")
		newtopicID, err := strconv.Atoi(stringID)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du nouveau sujet : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.StatusBadRequest(w)
			return
		}

		err = admin.MoveMessage(newtopicID, postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans le déplacement du message' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		topic, _ := getdata.GetTopicInfo(db, newtopicID)
		logMsg := fmt.Sprintf("ADMIN : Message n°%d déplacé par %s dans le sujet n°%d (%s)\n", postID, username, newtopicID, topic.Name)
		logs.AddLogsToDatabase(logMsg)
		url := fmt.Sprintf("/topic/%d#%d", newtopicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	default:
		logMsg := fmt.Sprint("ERREUR : <messageactionshandler.go> Requête invalide sur la page message : ", r.FormValue("action"))
		logs.AddLogsToDatabase(logMsg)
		utils.StatusBadRequest(w)
		return
	}
}

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
