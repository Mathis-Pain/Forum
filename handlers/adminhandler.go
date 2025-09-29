package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

var funcShort = template.FuncMap{
	"preview": getdata.Preview,
}

// Gestion des pages du panneau d'administration
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<profilhandler.go> Erreur à l'ouverture de la base de données :", err)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Récupère la liste des catégories et l'utilisateur connecté
	categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		log.Printf("<cathandler.go> Erreur dans la construction du header : %v\n", err)
		utils.InternalServError(w)
		return
	}

	// Vérifie si l'utilisateur en ligne a bien un statut ADMIN
	isAdmin, err := admin.CheckIfAdmin(currentUser.Username)
	if !isAdmin && err == nil {
		// Si l'utilisateur n'est pas un admin, envoie une erreur "Unauthorized"
		log.Print("Tentative d'accès non autorisé au panneau d'administration.")
		utils.UnauthorizedError(w)
		return
	} else if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Tentative d'accès non autorisé au panneau d'administration.")
			utils.UnauthorizedError(w)
			return
		}
		log.Print("<adminhandler.go> Erreur dans la vérification des accréditations :", err)
		utils.InternalServError(w)
		return
	}

	// Récupère l'url pour savoir quelle page du panneau admin ouvrir
	parts := strings.Split(r.URL.Path, "/")

	// Met à jour la liste des catégories avec la liste complète de tous les sujets
	categories, topics, err := admin.GetAllTopics(categories, db)
	if err != nil {
		log.Print("Erreur dans la récupération des sujets : ", err)
		utils.InternalServError(w)
		return
	}

	// Récupère les statistiques  du forum
	lastmonthpost, stats, users, err := admin.GetStats(topics)
	if err != nil {
		log.Print("Erreur dans la récupération des statistiques : ", err)
		utils.InternalServError(w)
		return
	}
	stats.TotalCats = len(categories)
	stats.LastCat = categories[0].Name

	// Analyse l'url pour choisir quelle page afficher
	if len(parts) == 3 && parts[2] == "" {
		// S'il n'y a aucune page spécifique demandée, affiche l'accueil du panneau d'aministration
		adminHome(categories, topics, stats, users, w, currentUser, lastmonthpost)
	} else {
		switch parts[2] {
		case "userlist":
			// Affiche la liste des utilisateurs
			adminUsers(users, r, w, currentUser, stats)
		case "catlist":
			// Affiche la liste des catégories
			adminCategories(categories, r, w, currentUser, stats)
		case "topiclist":
			// Affiche la liste des sujets
			adminTopics(topics, categories, r, w, currentUser, stats)
		case "seeposts":
			// Affiche les messages du dernier moi
			adminPost(lastmonthpost, r, w, currentUser, stats)
		}
	}
}

// Fonction pour affiche la liste des messages (non opérationnelle pour l'instant)
func adminPost(lastmonthpost []models.LastPost, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats) {
	data := struct {
		PageName    string
		LastMonth   []models.LastPost
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
	}{
		PageName:    "Messages du dernier mois",
		LastMonth:   lastmonthpost,
		CurrentUser: currentUser,
		Stats:       stats,
	}

	pageToLoad, err := template.ParseFiles("templates/all-posts.html", "templates/header.html", "templates/initpage.html")
	if err != nil {
		log.Printf("<adminhandler.go> Erreur dans la génération du template adminPost : %v", err)
		utils.InternalServError(w)
		return
	}

	err = pageToLoad.Execute(w, data)
	if err != nil {
		utils.InternalServError(w)
		return
	}
}

// Page admin de la liste des sujets.
func adminTopics(topics []models.Topic, categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats) {
	// Données à renvoyer à la page
	data := struct {
		PageName    string
		Topics      []models.Topic
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
	}{
		PageName:    "Administration des sujets",
		Topics:      topics,
		Categories:  categories,
		CurrentUser: currentUser,
		Stats:       stats,
	}

	// Si un formulaire (modification ou suppression de post) a été utilisé
	if r.Method == "POST" {
		if name := r.FormValue("topicname"); name != "" {
			// Si un sujet a été modifié
			err := subhandlers.EditTopicHandler(r, topics)
			if err != nil {
				log.Print("<adminhandler.go adminTopics> Erreur dans la modification du sujet : ", err)
				utils.InternalServError(w)
				return
			}

		} else if stringID := r.FormValue("topicToDelete"); stringID != "" {
			// Si on clique pour supprimer un sujet
			err := subhandlers.DeleteTopicHandler(stringID)
			if err != nil {
				log.Print("<adminhandler.go adminTopics> Erreur dans la suppression du sujet : ", err)
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
		log.Printf("<adminhandler.go> Erreur dans la génération du template adminTopics : %v", err)
		utils.InternalServError(w)
		return
	}

	// Charge la page en fonction des informations récupérées
	err = pageToLoad.Execute(w, data)
	if err != nil {
		log.Print("Erreur à l'ouverture de la page adminTopic :", err)
		utils.InternalServError(w)
		return
	}
}

// Page admin de la liste des catégories
func adminCategories(categories []models.Category, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats) {
	data := struct {
		PageName    string
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
	}{
		PageName:    "Administration des catégories",
		Categories:  categories,
		CurrentUser: currentUser,
		Stats:       stats,
	}

	// Si un formulaire (créer, modifier ou supprimer) a été utilisé
	if r.Method == "POST" {
		if stringID := r.FormValue("catToDelete"); stringID != "" {
			// Si on clique pour supprimer une catégorie
			err := subhandlers.DeleteCatHandler(stringID)
			if err != nil {
				log.Print("<adminhandler.go adminCategories> Erreur dans la suppression de la catégorie : ", err)
				utils.InternalServError(w)
				return
			}
		} else if newcat := r.FormValue("newcatname"); newcat != "" {
			// Si on crée une nouvelle catégorie
			err := subhandlers.AddCatHandler(r)
			if err != nil {
				log.Print("<adminhandler.go adminCategories> Erreur dans la création de la catégorie : ", err)
				utils.InternalServError(w)
				return
			}
		} else {
			// Verifie si un formulaire de modification de catégorie a été envoyé
			categ, isModified, err := subhandlers.AdminIsCatModified(r, categories)
			if err != nil {
				log.Print("<adminhandler.go adminCategories> Erreur dans la modification de la catégorie : ", err)
				utils.InternalServError(w)
				return
			}

			// Si oui, appelle la fonction de modification de la catégorie
			if isModified {
				err := subhandlers.EditCatHandler(r, categ)
				if err != nil {
					log.Print("<adminhandler.go adminCategories> Erreur dans la modification de la catégorie : ", err)
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
		log.Printf("<adminhandler.go> Erreur dans la génération du template adminCategories : %v", err)
		utils.InternalServError(w)
		return
	}

	// Lance la page
	err = pageToLoad.Execute(w, data)
	if err != nil {
		log.Printf("<adminhandler.go> Erreur dans le chargement du template adminCategories : %v", err)
		utils.InternalServError(w)
		return
	}
}

// Page admin de la liste des utilisateurs
func adminUsers(users []models.User, r *http.Request, w http.ResponseWriter, currentUser models.UserLoggedIn, stats models.Stats) {
	data := struct {
		PageName    string
		Users       []models.User
		CurrentUser models.UserLoggedIn
		Stats       models.Stats
	}{
		PageName:    "Administrer les utilisateurs",
		Users:       users,
		CurrentUser: currentUser,
		Stats:       stats,
	}

	// Si un formulaire (modifier, bannir, supprimer) a été envoyé
	if r.Method == "POST" {
		// Si un compte utilisateur est modifié
		if username := r.FormValue("username"); username != "" {
			err := subhandlers.UserEditHandler(r, users)
			if err != nil {
				log.Print("<adminhandler.go adminUsers> Erreur dans la modification de l'utilisateur : ", err)
				utils.InternalServError(w)
				return
			}
		} else if stringID := r.FormValue("userToBan"); stringID != "" {
			// Si on clique pour bannir un utilisateur
			err := subhandlers.BanUserHandler(stringID)
			if err != nil {
				log.Print("<adminhandler.go adminUsers> Erreur dans le bannissement de l'utilisateur : ", err)
				utils.InternalServError(w)
				return
			}
		} else if stringID := r.FormValue("userToFree"); stringID != "" {
			// Si on clique pour débannir un utilisateur
			err := subhandlers.UnbanUserHandler(stringID)
			if err != nil {
				log.Print("<adminhandler.go adminUsers> Erreur dans le débannissement de l'utilisateur : ", err)
				utils.InternalServError(w)
				return
			}
		} else if stringID := r.FormValue("userToDelete"); stringID != "" {
			// Si on supprime un utilisateur
			err := subhandlers.DeleteUserHandler(stringID)
			if err != nil {
				log.Print("<adminhandler.go adminUsers> Erreur dans la suppression de l'utilisateur : ", err)
				utils.InternalServError(w)
				return
			}
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
		log.Printf("<adminhandler.go> Erreur dans la génération du template adminUsers : %v", err)
		utils.InternalServError(w)
		return
	}

	// Lancement de la page
	err = pageToLoad.Execute(w, data)
	if err != nil {
		log.Print("<adminhandler.go> Erreur dans la lecture du template adminUsers : ", err)
		utils.InternalServError(w)
		return
	}
}

func adminHome(categories []models.Category, topics []models.Topic, stats models.Stats, users []models.User, w http.ResponseWriter, currentUser models.UserLoggedIn, postList []models.LastPost) {
	data := struct {
		PageName    string
		Categories  []models.Category
		Topics      []models.Topic
		Users       []models.User
		Stats       models.Stats
		PostList    []models.LastPost
		CurrentUser models.UserLoggedIn
	}{
		PageName:    "Panneau d'administration",
		Categories:  categories,
		Topics:      topics,
		Users:       users,
		Stats:       stats,
		PostList:    postList,
		CurrentUser: currentUser,
	}

	pageToLoad := template.Must(template.New("admin.html").Funcs(funcShort).ParseFiles("templates/admin/admin.html",
		"templates/admin/adminheader.html",
		"templates/admin/adminsidebar.html",
		"templates/initpage.html"))

	err := pageToLoad.Execute(w, data)
	if err != nil {
		log.Print("<adminhandler.go> Erreur dans la lecture du template adminHome : ", err)
		utils.InternalServError(w)
		return
	}
}
