package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
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
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	getcategoryID := r.URL.Query().Get("category_id")

	getcatID, err := strconv.Atoi(getcategoryID)
	if err != nil {
		// gérer l'erreur si category_id n'est pas un nombre
		http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
		return
	}
	// on charge es categories et l'utilisateur pour construire le header
	categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		log.Printf("<cathandler.go> Erreur dans la construction du header : %v\n", err)
		utils.InternalServError(w)
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
			fmt.Println("Erreur de conversion dans creat-topic-handler catID:", err)
			return
		}
		if topicName == "" || message == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		// --- Récupération deuserID ---
		_, userID, _ := utils.GetUserNameAndIDByCookie(r, db)
		postactions.CreateNewtopic(userID, catID, topicName, message)

		// Redirection vers la page catégorie
		http.Redirect(w, r, fmt.Sprintf("/categorie/%d", catID), http.StatusSeeOther)
		return
	}

	data := struct {
		PageName    string
		Category    models.Category
		CurrentUser models.UserLoggedIn
		Categories  []models.Category
		LoginData   models.LoginData
	}{
		PageName:    "Forum",
		Category:    currentCategory,
		CurrentUser: currentUser,
		Categories:  categories,
		LoginData:   models.LoginData{},
	}

	err = CreatTopicHtml.Execute(w, data)
	if err != nil {
		log.Printf("<create-topic-handler.go> Could not execute template <create-topic.html>: %v\n", err)
		utils.NotFoundHandler(w)

	}
}
