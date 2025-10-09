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
	_ "github.com/mattn/go-sqlite3"
)

var funcMap = template.FuncMap{
	"preview": getdata.Preview,
}

var HomeHtml = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles(
	"templates/home.html", "templates/login.html", "templates/header.html", "templates/initpage.html",
))

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	// --- Récupération des derniers posts ---
	lastPosts, err := getdata.GetLastPosts()

	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur à l'exécution de GetLastPosts: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// --- Récupération des catégories ---
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur à l'ouverture de la base de données : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// --- Gestion des erreurs de login ---

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur à l'exécution de GetSessionFromRequest: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	var loginErr string
	if session.ID != "" {
		loginErr, err = getdata.GetLoginErr(session)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur à l'exécution de GetLoginErr: %v", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
	}

	// --- Structure de données ---

	data := struct {
		PageName      string
		LoginErr      string
		Posts         []models.LastPost
		Categories    []models.Category
		CurrentUser   models.UserLoggedIn
		Notifications models.Notifications
	}{
		PageName:      "Petites victoires",
		LoginErr:      loginErr,
		Posts:         lastPosts,
		Categories:    categories,
		CurrentUser:   currentUser,
		Notifications: notifications,
	}

	// --- Sinon : Renvoi des données de base au template ---
	err = HomeHtml.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <homehandler.go> Erreur à l'exécution de template <home.html>: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.NotFoundHandler(w)
		return
	}
}
