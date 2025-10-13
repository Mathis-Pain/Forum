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

// Fonctions pour la gestion des actions de message (supprimer, modifier, déplacer un message)
func MessageActionsHandler(w http.ResponseWriter, r *http.Request) {
	// MARK: Initialisation
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

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messageactionshandler> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return
	}
	defer db.Close()

	// Récupère les infos de l'utilisateur connecté pour les notifications et les logs
	_, _, currentUser, err := BuildHeader(r, w, db)

	// Appel des fonctions selon l'action renvoyée par le formulaire
	switch r.FormValue("action") {

	// MARK: Suppression
	// d'un message
	case "delete":
		err := admin.AdminDeleteMessage(topicID, postID, db, currentUser)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la suppression du message : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}

		logMsg := fmt.Sprintf("ADMIN : Le message %d a été supprimé par %s.", postID, currentUser.Username)
		logs.AddLogsToDatabase(logMsg)

		// Redirection vers le sujet s'il existe encore, sinon redirige vers la page d'accueil
		if _, err := getdata.GetTopicInfo(db, topicID); err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
		return

		// MARK: Signalement
		// d'un message par un modérateur
	case "warn":
		warnReason := r.FormValue("warning")
		err := admin.ModSignalMessage(postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans le signalement du message : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		// Ajout du log
		logMsg := fmt.Sprintf("WARN : Le message n°%d a été signalé par %s pour la raison suivante : %s\n", postID, currentUser.Username, warnReason)
		logs.AddLogsToDatabase(logMsg)
		AddSenderID(logMsg, currentUser.Username)

		// Redirection avec les données modifiées
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return

		// MARK: Annulation du signalement
		// par un administrateur
	case "cancel":
		err := admin.AdminCancelSignal(postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans l'annulation du signalement' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
		// Ajout du log
		logMsg := fmt.Sprintf("ADMIN : %s a annulé le signalement du message n°%d.\n", currentUser.Username, postID)
		logs.AddLogsToDatabase(logMsg)

		// Redirection avec les données modifiées
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return

		// MARK: Modification
		// d'un message
	case "edit":
		content := r.FormValue("message-content")
		err := admin.EditExistingMessage(postID, db, content)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la modification du message' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}

		// Redirection avec les données modifiées
		url := fmt.Sprintf("/topic/%d#%d", topicID, postID)
		topic, _ := getdata.GetTopicInfo(db, topicID)

		// Ajout du log
		logMsg := fmt.Sprintf("USER : %s a modifié le contenu du message n°%d (sur \"%s\")", currentUser.Username, postID, topic.Name)
		logs.AddLogsToDatabase(logMsg)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return

		// MARK: Déplacement
		// d'un message
	case "move":
		// Récupère les informations du sujet sur lequel déplacer le message
		stringID := r.FormValue("topicindex")
		newtopicID, err := strconv.Atoi(stringID)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans la récupération de l'ID du nouveau sujet : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.StatusBadRequest(w)
			return
		}

		// Déplace le message
		err = admin.MoveMessage(newtopicID, postID, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <messageactionshandler> Erreur dans le déplacement du message' : ", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}

		// Ajout du log
		topic, _ := getdata.GetTopicInfo(db, newtopicID)
		logMsg := fmt.Sprintf("ADMIN : Message n°%d déplacé par %s dans le sujet n°%d (%s)\n", postID, currentUser.Username, newtopicID, topic.Name)
		logs.AddLogsToDatabase(logMsg)

		// Redirection avec les données modifiées
		url := fmt.Sprintf("/topic/%d#%d", newtopicID, postID)
		http.Redirect(w, r, url, http.StatusSeeOther)
		return

		// MARK: Erreur
	default:
		// Si l'action n'est pas prise en charge, renvoie un Bad Request et un log d'erreur
		logMsg := fmt.Sprint("ERREUR : <messageactionshandler.go> Requête invalide sur la page message : ", r.FormValue("action"))
		logs.AddLogsToDatabase(logMsg)
		utils.StatusBadRequest(w)
		return
	}
}
