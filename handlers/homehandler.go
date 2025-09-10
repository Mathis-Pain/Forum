package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	_ "github.com/mattn/go-sqlite3"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var funcMap = template.FuncMap{
		"preview": utils.Preview,
	}

	var HomeHtml = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles(
		"templates/home.html",
		"templates/login.html",
		"templates/header.html",
	))

	// --- Récupération des derniers posts ---
	lastPosts, err := utils.GetLastPosts()

	if err != nil {
		log.Printf("<homehandler.go> Could not oprate GetLastPosts: %v\n", err)
		utils.InternalServError(w)
		return
	}

	// --- Récupération des catégories ---

	categories, err := utils.GetCatList()

	if err != nil {
		log.Printf("<homehandler.go> Could not operate GetCatList: %v\n", err)
		utils.InternalServError(w)
		return
	}

	// --- Structure de données ---

	data := struct {
		LoginData  models.LoginData
		Posts      []models.LastPost
		Categories []models.Category
	}{
		LoginData:  models.LoginData{},
		Posts:      lastPosts,
		Categories: categories,
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
