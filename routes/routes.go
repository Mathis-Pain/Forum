package routes

import (
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers"
	"github.com/Mathis-Pain/Forum/middleware"

	"github.com/Mathis-Pain/Forum/utils"
)

func InitRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	// Route Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			utils.NotFoundHandler(w)
			return
		}
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/registration", handlers.SignUpSubmitHandler)
	mux.HandleFunc("/profil", middleware.AuthMiddleware(handlers.ProfilHandler))
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/categorie", handlers.CategoriesHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
