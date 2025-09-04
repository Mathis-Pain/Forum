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

var CatHtml = template.Must(template.ParseFiles("templates/category.html"))

// Permet au HTMl d'utiliser la fonction preview
var funcMap = template.FuncMap{
	"preview": utils.Preview,
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<homehandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	parts := strings.Split(r.URL.Path, "/")
	path := parts[len(parts)-1]

	if !strings.Contains(path, "c") {
		utils.NotFoundHandler(w)
	}

	catID := strings.Trim(path, "c")
	ID, err := strconv.Atoi(catID)

	if err != nil {
		utils.InternalServError(w)
	}

	category, err := utils.GetCatDetails(db, ID)

	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
	} else if err != nil {
		utils.InternalServError(w)
	}

	data := struct {
		Category models.Category
	}{
		Category: category,
	}

	err = CatHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template <categorie.html> : %v\n", err)
		utils.NotFoundHandler(w)
	}

}
