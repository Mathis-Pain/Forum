package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var CatHtml = template.Must(template.New("categorie.html").Funcs(funcMap).ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ID := subhandlers.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// --- Récupération des catégories ---

	category, err := getdata.GetCatDetails(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la récupération de la catégorie : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	for i := 0; i < len(category.Topics); i++ {
		category.Topics[i].Messages = getdata.FormatDateAllMessages(category.Topics[i].Messages)
	}

	// --- Gestion des erreurs de login ---

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution de GetSessionFromRequest: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	var loginErr string
	if session.ID != "" {
		loginErr, err = getdata.GetLoginErr(session)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution de GetLoginErr: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
	}

	// --- Renvoi des données ---

	data := struct {
		PageName      string
		Category      models.Category
		Categories    []models.Category
		LoginErr      string
		CurrentUser   models.UserLoggedIn
		Notifications models.Notifications
	}{
		PageName:      category.Name,
		Category:      category,
		Categories:    categories,
		LoginErr:      loginErr,
		CurrentUser:   currentUser,
		Notifications: notifications,
	}

	err = CatHtml.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution du template <categorie.html> : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

}
