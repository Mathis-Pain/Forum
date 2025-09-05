package routes

import (
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers"
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

	mux.HandleFunc("/signup", handlers.SignUpFormHandler)
	mux.HandleFunc("/registration", handlers.SignUpSubmitHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
