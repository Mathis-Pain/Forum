package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var funcShort = template.FuncMap{
	"preview": getdata.Preview,
}

// Gestion des pages du panneau d'administration
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur à l'ouverture de la base de données : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Récupère la liste des catégories et l'utilisateur connecté
	_, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf(" ERREUR : <adminhandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Vérifie si l'utilisateur en ligne a bien un statut ADMIN
	isAdmin, err := admin.CheckIfAdmin(currentUser.Username)
	if !isAdmin && err == nil {
		// Si l'utilisateur n'est pas un admin, envoie une erreur "Unauthorized"
		utils.UnauthorizedError(w)
		return
	} else if err != nil {
		if err == sql.ErrNoRows {
			utils.UnauthorizedError(w)
			return
		}
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la vérification des accréditations :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère l'url pour savoir quelle page du panneau admin ouvrir
	parts := strings.Split(r.URL.Path, "/")

	// Met à jour la liste des catégories avec la liste complète de tous les sujets
	categories, topics, err := admin.GetAllTopics(categories, db)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la récupération des sujets : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].TopicID > topics[j].TopicID
	})

	logList, err := logs.DisplayLogs()
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la récupération des logs : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère les statistiques  du forum
	lastmonthpost, stats, users, err := admin.GetStats(topics)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la récupération des statistiques : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	stats.TotalCats = len(categories)

	index := len(categories) - 1
	if index < 0 {
		index = 0
	}
	stats.LastCat = categories[index].Name

	// Analyse l'url pour choisir quelle page afficher
	if len(parts) == 3 && parts[2] == "" {
		// S'il n'y a aucune page spécifique demandée, affiche l'accueil du panneau d'aministration
		adminHome(categories, topics, stats, users, w, currentUser, lastmonthpost, logList)
	} else {
		switch parts[2] {
		case "userlist":
			// Affiche la liste des utilisateurs
			adminUsers(users, r, w, currentUser, stats, logList)
		case "catlist":
			// Affiche la liste des catégories
			adminCategories(categories, r, w, currentUser, stats, logList)
		case "topiclist":
			// Affiche la liste des sujets
			adminTopics(topics, categories, r, w, currentUser, stats, logList)
		case "logs":
			// Affiche la liste des logs
			adminLogs(r, w, currentUser, stats, logList)
		default:
			utils.NotFoundHandler(w)
			return
		}
	}
}

// Page Liste des sujets du panneau d'administration
func adminTopics(topics []models.Topic, categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
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
			err := subhandlers.EditTopicHandler(r, topics, currentUser.Username)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminTopics> Erreur dans la modification du sujet : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

		} else if stringID := r.FormValue("topicToDelete"); stringID != "" {
			// Si on clique pour supprimer un sujet
			err := subhandlers.DeleteTopicHandler(stringID)
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

// Page la liste des catégories du panneau d'administration
func adminCategories(categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
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
			err := subhandlers.DeleteCatHandler(stringID)
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
			err := subhandlers.AddCatHandler(r)
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
			categ, isModified, err := subhandlers.AdminIsCatModified(r, categories)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminCategories> Erreur dans la vérification des modifications de la catégorie : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			// Si oui, appelle la fonction de modification de la catégorie
			if isModified {
				err := subhandlers.EditCatHandler(r, categ, currentUser)
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

// Page la liste des utilisateurs du panneau d'administration
func adminUsers(users []models.User, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
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
			err := subhandlers.UserEditHandler(r, users, currentUser)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans la modification de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

		case "ban":
			// Si on clique pour bannir un utilisateur
			err := subhandlers.BanUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans le bannissement de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := convertIDtoUsername(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s (ID : %s) a été banni par %s", username, stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		case "unban":
			// Si on clique pour débannir un utilisateur
			err := subhandlers.UnbanUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans le débannissement de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := convertIDtoUsername(stringID)
			if err != nil {
				utils.InternalServError(w)
				return
			}
			logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s (ID : %s) a été débanni par %s", username, stringID, currentUser.Username)
			logs.AddLogsToDatabase(logMsg)
		case "delete":
			// Si on supprime un utilisateur
			err := subhandlers.DeleteUserHandler(stringID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go adminUsers> Erreur dans la suppression de l'utilisateur : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			username, err := convertIDtoUsername(stringID)
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

// Page principale du panneau d'administration
func adminHome(categories []models.Category, topics []models.Topic, stats models.Stats, users []models.User, w http.ResponseWriter, currentUser models.UserLoggedIn, postList []models.LastPost, logList []models.Log) {
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

// Page logs du panneau d'administration
func adminLogs(r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats, logList []models.Log) {
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

			refusedUser, _ := convertIDtoUsername(receiver)
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

			err = subhandlers.PromoteToMod(receiverID)
			if err != nil {
				logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans l'ajout de l'utilisateur à la modération : ", err)
				logs.AddLogsToDatabase(logMsg)
				utils.InternalServError(w)
				return
			}

			promoted, _ := convertIDtoUsername(receiver)
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

// Fonction pour récupérer le pseudo à partir de l'ID
func convertIDtoUsername(stringID string) (string, error) {
	ID, err := strconv.Atoi(stringID)

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	user, err := getdata.GetUserInfoFromID(db, ID)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
