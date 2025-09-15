package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
)

var CatHtml = template.Must(template.New("categorie.html").Funcs(funcMap).ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		utils.NotFoundHandler(w)
		return
	}

	ID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	category, err := utils.GetCatDetails(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		log.Printf("<cathandler.go> Could not operate GetCatDetails: %v\n", err)
		utils.InternalServError(w)
		return
	}

	categories, err := utils.GetCatList()

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

	err = CatHtml.Execute(w, data)
	if err != nil {
		log.Printf("<cathandler.go> Could not execute template <categorie.html> : %v\n", err)
		utils.InternalServError(w)
		return
	}

}
