package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

var CatHtml = template.Must(template.New("categorie.html").Funcs(funcMap).ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ID := utils.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "data/forum.db")
	if err != nil {
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	category, err := getdata.GetCatDetails(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		log.Printf("<cathandler.go> Could not operate GetCatDetails: %v\n", err)
		utils.InternalServError(w)
		return
	}

	categories, err := getdata.GetCatList()

	if err != nil {
		log.Printf("<cathandler.go> Could not operate GetCatList: %v\n", err)
		utils.InternalServError(w)
		return
	}

	data := struct {
		Category   models.Category
		Categories []models.Category
		LoginData  models.LoginData
	}{
		Category:   category,
		Categories: categories,
		LoginData:  models.LoginData{},
	}

	// --- Si POST, on remplit LoginData ---

	if r.Method == "POST" {
		loginData, err := utils.LoginPopUp(r, w)
		if err == nil {
			data.LoginData = loginData
		}

		// Connexion réussie (ouverture de session, accès aux boutons, etc, à ajouter ici)

	}

	err = CatHtml.Execute(w, data)
	if err != nil {
		log.Printf("<cathandler.go> Could not execute template <categorie.html> : %v\n", err)
		utils.InternalServError(w)
		return
	}

}
