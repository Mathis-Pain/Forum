package handlers

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	categories := []models.Category{
		{ID: 1, Name: "Informatique"},
		{ID: 2, Name: "Cuisine"},
	}

	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}

	tmpl := template.Must(template.ParseFiles("templates/test.html"))
	tmpl.Execute(w, data)
}
