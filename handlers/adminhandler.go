package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/Mathis-Pain/Forum/handlers/adminhandlers"
	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// Gestion des pages du panneau d'administration
// Gère toutes les pages dont l'url contient "/admin/" puis redirige vers le handler correspondant dans <adminpages.go>
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

	// Range les sujets par ordre de publication
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].TopicID > topics[j].TopicID
	})

	// Récupère la liste des logs dans la base de données
	logList, err := logs.DisplayLogs()
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminhandler.go> Erreur dans la récupération des logs : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère les statistiques du forum (dernier utilisateur, derniuère catégorie, etc)
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
		adminhandlers.AdminHome(categories, topics, stats, users, w, currentUser, lastmonthpost, logList)
	} else {
		switch parts[2] {
		case "userlist":
			// Affiche la liste des utilisateurs
			adminhandlers.AdminUsers(users, r, w, currentUser, stats, logList)
		case "catlist":
			// Affiche la liste des catégories
			adminhandlers.AdminCategories(categories, r, w, currentUser, stats, logList)
		case "topiclist":
			// Affiche la liste des sujets
			adminhandlers.AdminTopics(topics, categories, r, w, currentUser, stats, logList)
		case "logs":
			// Affiche la liste des logs
			adminhandlers.AdminLogs(r, w, currentUser, stats, logList)
		default:
			utils.NotFoundHandler(w)
			return
		}
	}
}
