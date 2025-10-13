package adminhandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var funcShort = template.FuncMap{
	"preview": getdata.Preview,
}

// ************** Pages du panneau d'administration (appelées par adminhandler)

// MARK: Page principale
func AdminHome(categories []models.Category, topics []models.Topic, stats models.Stats, users []models.User, w http.ResponseWriter, currentUser models.UserLoggedIn, postList []models.LastPost, logList []models.Log) {
	data := struct {
		PageName    string
		Categories  []models.Category
		Topics      []models.Topic
		Users       []models.User
		Stats       models.Stats
		PostList    []models.LastPost
		CurrentUser models.UserLoggedIn
		LogList     []models.Log
	}{
		PageName:    "Panneau d'administration",
		Categories:  categories,
		Topics:      topics,
		Users:       users,
		Stats:       stats,
		PostList:    postList,
		CurrentUser: currentUser,
		LogList:     logList,
	}

	pageToLoad := template.Must(template.New("admin.html").Funcs(funcShort).ParseFiles("templates/admin/admin.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html"))

	err := pageToLoad.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la lecture du template adminHome : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
}

// MARK: Liste des sujets
func AdminTopics(topics []models.Topic, categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
	// Données à renvoyer à la page
	data := struct {
		PageName    string
		Topics      []models.Topic
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
		LogList     []models.Log
	}{
		PageName:    "Administration des sujets",
		Topics:      topics,
		Categories:  categories,
		CurrentUser: currentUser,
		Stats:       stats,
		LogList:     logList,
	}

	// Si un formulaire (modification ou suppression de post) a été utilisé
	if r.Method == "POST" {
		if name := r.FormValue("topicname"); name != "" {
			// Si un sujet a été modifié
			err := EditTopicHandler(r, topics, currentUser.Username)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminTopics> Erreur dans la modification du sujet : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

		} else if stringID := r.FormValue("topicToDelete"); stringID != "" {
			// Si on clique pour supprimer un sujet
			err := DeleteTopicHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminTopics> Erreur dans la suppression du sujet : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}
		}
		// Renvoie la page avec les informations mises à jour
		http.Redirect(w, r, "/admin/topiclist", http.StatusSeeOther)
	}

	// Crée le template complet de la page avec les sous-templates
	pageToLoad, err := template.ParseFiles("templates/admin/all-topics.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminhandler.go> Erreur dans la génération du template adminTopics : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Charge la page en fonction des informations récupérées
	err = pageToLoad.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur à l'ouverture de la page adminTopic :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
}

// MARK: Liste des catégories
func AdminCategories(categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
	data := struct {
		PageName    string
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
		LogList     []models.Log
	}{
		PageName:    "Administration des catégories",
		Categories:  categories,
		CurrentUser: currentUser,
		Stats:       stats,
		LogList:     logList,
	}

	// Si un formulaire (créer, modifier ou supprimer) a été utilisé
	if r.Method == "POST" {
		if stringID := r.FormValue("catToDelete"); stringID != "" {
			// Si on clique pour supprimer une catégorie
			err := DeleteCatHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminCategories> Erreur dans la suppression de la catégorie : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : La catégorie %s a été supprimée par %s", stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		} else if newcat := r.FormValue("newcatname"); newcat != "" {
			// Si on crée une nouvelle catégorie
			err := AddCatHandler(r)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminCategories> Erreur dans la création de la catégorie : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : %s a créé une nouvelle catégorie : %s", currentUser.Username, newcat)
			logs.AddLogsToDatabase(logMsg)
		} else {
			// Verifie si un formulaire de modification de catégorie a été envoyé
			categ, isModified, err := AdminIsCatModified(r, categories)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminCategories> Erreur dans la vérification des modifications de la catégorie : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			// Si oui, appelle la fonction de modification de la catégorie
			if isModified {
				err := EditCatHandler(r, categ, currentUser)
				if err != nil {
					logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminCategories> Erreur dans la modification de la catégorie : ", err)
					logs.AddLogsToDatabase(logMsg)
					utils.InternalServError(w)
					return
				}
			}

		}

		// Renvoie la page avec les modifications
		http.Redirect(w, r, "/admin/catlist", http.StatusSeeOther)
	}

	// Charge le template complet de la page
	pageToLoad, err := template.ParseFiles("templates/admin/all-categories.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html")
	if err != nil {
		logMsg := fmt.Sprintf(" ERREUR : <adminhandler.go> Erreur dans la génération du template adminCategories : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Lance la page
	err = pageToLoad.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf(" ERREUR : <adminhandler.go> Erreur dans le chargement du template adminCategories : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
}

// MARK: Liste des utilisateurs
func AdminUsers(users []models.User, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
	data := struct {
		PageName    string
		Users       []models.User
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
		LogList     []models.Log
	}{
		PageName:    "Administrer les utilisateurs",
		Users:       users,
		CurrentUser: currentUser,
		Stats:       stats,
		LogList:     logList,
	}

	// Si un formulaire (modifier, bannir, supprimer) a été envoyé
	if r.Method == "POST" {

		stringID := r.FormValue("userID")
		switch r.FormValue("action") {
		case "edit":
			// Si un compte utilisateur est modifié
			err := UserEditHandler(r, users, currentUser)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans la modification de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

		case "ban":
			// Si on clique pour bannir un utilisateur
			err := BanUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans le bannissement de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := ConvertIDtoUsername(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s (ID : %s) a été banni par %s", username, stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		case "unban":
			// Si on clique pour débannir un utilisateur
			err := UnbanUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans le débannissement de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := ConvertIDtoUsername(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s (ID : %s) a été débanni par %s", username, stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		case "delete":
			// Si on supprime un utilisateur
			err := DeleteUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans la suppression de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := ConvertIDtoUsername(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s (ID : %s) a été supprimé par %s", username, stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		}

		// Redirection avec les données mises à jour
		http.Redirect(w, r, "/admin/userlist", http.StatusSeeOther)
	}

	// Création du template
	pageToLoad, err := template.ParseFiles(
		"templates/admin/all-users.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html")

	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminhandler.go> Erreur dans la génération du template adminUsers : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Lancement de la page
	err = pageToLoad.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la lecture du template adminUsers : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
}

// MARK: Logs
func AdminLogs(r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
	data := struct {
		PageName    string
		Stats       models.Stats
		CurrentUser models.UserLoggedIn
		LogList     []models.Log
	}{
		PageName:    "Panneau d'administration",
		Stats:       stats,
		CurrentUser: currentUser,
		LogList:     logList,
	}

	if r.Method == "POST" {

		switch r.FormValue("action") {
		case "delete":
			stringID := r.FormValue("logID")
			logID, err := strconv.Atoi(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logs.DeleteLog(logID)
		case "refuse":
			// Récupération des données du formulaire
			message := r.FormValue("logmessage")
			receiver := r.FormValue("receiver")
			receiverID, err := strconv.Atoi(receiver)
			if err != nil || receiverID == 0 {
				receiverID, _ = logs.GetIDFromLog(message)
				receiver = strconv.Itoa(receiverID)
			}
			logs.MarkAsHandled(message, receiverID)

			notif := "Votre demande de rejoindre la modération a été refusée."
			logs.AddNotificationToDatabase("ANSWER", receiverID, 0, notif)

			refusedUser, _ := ConvertIDtoUsername(receiver)
			logMsg := fmt.Sprintf("ADMIN : %s a refusé la demande de %s de rejoindre la modération", currentUser.Username, refusedUser)
			logs.AddLogsToDatabase(logMsg)
		case "accept":
			// Récupération des données du formulaire
			message := r.FormValue("logmessage")
			receiver := r.FormValue("receiver")
			receiverID, err := strconv.Atoi(receiver)
			if err != nil || receiverID == 0 {
				receiverID, _ = logs.GetIDFromLog(message)
				receiver = strconv.Itoa(receiverID)
			}
			logs.MarkAsHandled(message, receiverID)

			notif := "Votre demande de rejoindre la modération a été acceptée."
			logs.AddNotificationToDatabase("ANSWER", receiverID, 0, notif)

			err = PromoteToMod(receiverID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans l'ajout de l'utilisateur à la modération : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			promoted, _ := ConvertIDtoUsername(receiver)
			logMsg := fmt.Sprintf("ADMIN : %s a accepté la demande de %s de rejoindre la modération", currentUser.Username, promoted)
			logs.AddLogsToDatabase(logMsg)

		case "answer":
			// Récupération des données du formulaire
			answer := r.FormValue("answer")
			message := r.FormValue("logmessage")
			receiver := r.FormValue("receiver")
			receiverID, _ := strconv.Atoi(receiver)
			date := getdata.FormatDate(r.FormValue("date"))

			// Préparation de la notification
			notif := fmt.Sprintf("Vous avez reçu une réponse à votre signalement du %s : %s.", date, answer)
			logs.AddNotificationToDatabase("ANSWER", receiverID, 0, notif)
			logMsg := fmt.Sprintf("ADMIN : Réponse envoyé au signalement du %s  par %s : %s", date, currentUser.Username, answer)
			logs.AddLogsToDatabase(logMsg)

			err := logs.MarkAsHandled(message, receiverID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
		}

		// Redirection avec les données mises à jour
		http.Redirect(w, r, "/admin/logs", http.StatusSeeOther)
		return
	}

	pageToLoad, err := template.ParseFiles(
		"templates/admin/logs.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html")

	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <adminhandler.go> Erreur dans la génération du template adminLogs : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Lancement de la page
	err = pageToLoad.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la lecture du template adminLogs : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

}
