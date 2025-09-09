package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	// Permet au HTMl d'utiliser la fonction preview
	var funcMap = template.FuncMap{
		"preview": utils.Preview,
	}

	var CatHtml = template.Must(template.New("categorie.html").Funcs(funcMap).ParseFiles(
		"templates/login.html",
		"templates/header.html",
		"templates/categorie.html",
	))

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<cathandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	/*parts := strings.Split(r.URL.Path, "/")
	path := parts[len(parts)-1]

	fmt.Println("path : ", path)

	if !strings.Contains(path, "c") {
		utils.NotFoundHandler(w)
	}

	catID := strings.Trim(path, "c")
	ID, err := strconv.Atoi(catID)

	if err != nil {
		utils.InternalServError(w)
	}*/

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
		// Appeler la fonction avec l’ID
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

	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}

	err = CatHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template <categorie.html> : %v\n", err)
	}

}
