package handlers

import (
	"database/sql"
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
	// --- Ouverture de la db ---
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	// --- Récupération des données communes ---

	// - Récupération des derniers posts -
	lastPosts, err := utils.GetLastPosts()

	if err == sql.ErrNoRows {
		//Si il n'y a pas de posts on renvoie une slice vide
		lastPosts = []models.LastPost{}
	} else if err != nil {
		utils.InternalServError(w)
	}

	// - Récupération des catégories -

	var category models.Category
	var categories []models.Category

	// Préparer la requête
	rows, err := db.Query("SELECT id FROM category")
	if err != nil {
		log.Println("Erreur lors de la requête à la db table category : ", err)
		return
	}
	defer rows.Close()

	// Parcourir les résultats
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println("Erreur récupération des données category : ", err)
			return
		}
		category, err = utils.GetCatDetails(db, id)
		categories = append(categories, category)
	}

	// Vérifier les erreurs après la boucle

	if err == sql.ErrNoRows {
		log.Printf("Lecture depuis bd (category) : aucun donnée correpondante, %v", err)
		utils.NotFoundHandler(w)
	} else if err != nil {
		log.Printf("Erreur à l'appel de GetCatDetails : %v", err)
		utils.InternalServError(w)
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
		log.Printf("Erreur lors de l'exécution du template HomeHtml: %v\n", err)
		utils.NotFoundHandler(w)

	}
}
