package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
	"github.com/Mathis-Pain/Forum/utils/postactions"
)

var CreatTopicHtml = template.Must(template.New("create-topic.html").Funcs(funcMap).ParseFiles(
	"templates/create-topic.html",
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Could not open database : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	getcategoryID := r.URL.Query().Get("category_id")

	getcatID, err := strconv.Atoi(getcategoryID)
	if err != nil {
		utils.StatusBadRequest(w)
		return
	}
	// on charge les categories et l'utilisateur pour construire le header
	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Empêche les utilisateurs bannis ou non enrigistrés d'accéder à la page
	if currentUser.UserType == 4 || currentUser.UserType == 0 {
		utils.ForbiddenError(w)
		return
	}

	// on prend getcatID pour chercher la categories qui correspond dans la bdd et la donner au template
	var currentCategory models.Category
	for _, cat := range categories {
		if cat.ID == getcatID {
			currentCategory = cat
		}
	}

	// --- Création du topic
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			utils.InternalServError(w)
			return
		}
		// Verification username et password non nul
		topicName := r.FormValue("title")
		message := r.FormValue("message")
		stringcatID := r.FormValue("category_id")
		catID, err := strconv.Atoi(stringcatID)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <createtopichandler.go> L'ID de la catégorie n'est pas valide :", err)
			logs.AddLogsToDatabase(logMsg)
			utils.StatusBadRequest(w)
			return
		}
		if topicName == "" || message == "" {
			utils.StatusBadRequest(w)
			return
		}

		// --- Récupération deuserID ---
		username, userID, _ := utils.GetUserNameAndIDByCookie(r, db)
		postactions.CreateNewtopic(userID, catID, topicName, message)

		categ, _ := getdata.GetCatDetails(db, catID)

		logMsg := fmt.Sprintf("USER : Nouveau sujet ouvert dans la catégorie \"%s\" par %s : \"%s\"", categ.Name, username, topicName)
		logs.AddLogsToDatabase(logMsg)

		// Redirection vers la page de la catégorie
		http.Redirect(w, r, fmt.Sprintf("/categorie/%d", catID), http.StatusSeeOther)
		return
	}

	pagename := "Ouvrir un nouveau sujet - " + currentCategory.Name

	data := struct {
		PageName      string
		Category      models.Category
		CurrentUser   models.UserLoggedIn
		Categories    []models.Category
		LoginErr      string
		Notifications models.Notifications
	}{
		PageName:      pagename,
		Category:      currentCategory,
		CurrentUser:   currentUser,
		Categories:    categories,
		LoginErr:      "",
		Notifications: notifications,
	}

	err = CreatTopicHtml.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <create-topic-handler.go> Erreur à l'exécution du template <create-topic.html>: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.NotFoundHandler(w)

	}
}
