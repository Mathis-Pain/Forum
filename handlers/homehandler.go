package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	_ "github.com/mattn/go-sqlite3"
)

var funcMap = template.FuncMap{
	"preview": utils.Preview,
}

var HomeHtml = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles(
	"templates/home.html", "templates/login.html", "templates/header.html", "templates/initpage.html",
))

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	// --- Récupération des derniers posts ---
	lastPosts, err := getdata.GetLastPosts()

	if err != nil {
		log.Printf("<homehandler.go> Could not oprate GetLastPosts: %v\n", err)
		utils.InternalServError(w)
		return
	}

	// --- Récupération des catégories ---
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		log.Printf("<cathandler.go> Erreur dans la construction du header : %v\n", err)
		utils.InternalServError(w)
		return
	}

	// --- Structure de données ---

	data := struct {
		LoginData   models.LoginData
		Posts       []models.LastPost
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
	}{
		LoginData:   models.LoginData{},
		Posts:       lastPosts,
		Categories:  categories,
		CurrentUser: currentUser,
	}

	// --- Si POST, on remplit LoginData ---

	if r.Method == "POST" {
		loginData, err := utils.LoginPopUp(r, w)
		if err == nil {
			data.LoginData = loginData
		}

		// Connexion réussie (ouverture de session, accès aux boutons, etc, à ajouter ici)

	}

	// --- Sinon : Renvoi des données de base au template ---
	err = HomeHtml.Execute(w, data)
	if err != nil {
		log.Printf("<homehandler.go> Could not execute template <home.html>: %v\n", err)
		utils.NotFoundHandler(w)

	}
}
