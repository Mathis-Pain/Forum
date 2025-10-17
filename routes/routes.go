package routes

import (
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers"
	"github.com/Mathis-Pain/Forum/handlers/authhandlers"
	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/middleware"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/test"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/external"
)

func InitRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	// Route Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/home" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		if r.URL.Path != "/" {
			utils.NotFoundHandler(w)
			return
		}
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/test", test.TestHandler)

	//******** Pages secondaires (connexion, likes, actions, etc) avec redirection automatique
	mux.HandleFunc("/login", authhandlers.LoginHandler)
	mux.HandleFunc("/logout", authhandlers.LogOutHandler)
	mux.HandleFunc("/like", subhandlers.LikePostHandler)
	mux.HandleFunc("/dislike", subhandlers.DislikePostHandler)
	mux.HandleFunc("/messageactions", subhandlers.MessageActionsHandler)

	//******** Pages principales (accessibles)
	// Pages pour poster des messages
	mux.HandleFunc("/new-topic", handlers.CreateTopicHandler)
	mux.HandleFunc("/answermessage", handlers.MessageHandler)
	mux.HandleFunc("/registration", authhandlers.SignUpSubmitHandler)

	// Pages affichant des informations
	mux.Handle("/profil", middleware.AuthMiddleware(http.HandlerFunc(handlers.ProfilHandler)))
	mux.HandleFunc("/categorie/", handlers.CategoriesHandler)
	mux.HandleFunc("/admin/", handlers.AdminHandler)
	mux.HandleFunc("/topic/", handlers.TopicHandler)

	// Authentification par google ou github
	mux.HandleFunc("/auth/google/login", external.HandleGoogleLogin)
	mux.HandleFunc("/auth/google/callback", external.HandleGoogleCallback)
	mux.HandleFunc("/auth/github/login", external.HandleGitHubLogin)
	mux.HandleFunc("/auth/github/callback", external.HandleGitHubCallback)
	mux.HandleFunc("/auth/discord/login", external.HandleDiscordLogin)
	mux.HandleFunc("/auth/discord/callback", external.HandleDiscordCallback)

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
