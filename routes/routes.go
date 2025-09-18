package routes

import (
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers"
	"github.com/Mathis-Pain/Forum/middleware"
	"github.com/Mathis-Pain/Forum/sessions"

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
	mux.Handle("/profil", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfilHandler)))
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/categorie/", handlers.CategoriesHandler)
	mux.HandleFunc("/topic/", handlers.TopicHandler)
	mux.HandleFunc("/test", handlers.TestHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("session_id")
		fmt.Fprintf(w, "Cookie: %+v\n", cookie)
		if cookie != nil {
			session, err := sessions.GetSession(cookie.Value)
			fmt.Fprintf(w, "Session: %+v, err=%v\n", session, err)
		}
	})

	return mux
}
