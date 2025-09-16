package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
)

var ProfilHtml = template.Must(template.New("profil.html").ParseFiles("templates/profil.html", "templates/login.html", "templates/header.html", "templates/initpage.html"))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	// ** Récupération des catégories pour le header **
	categories, err := utils.GetCatList()

	if err != nil {
		log.Printf("<profilhandler.go> Could not operate GetCatList: %v\n", err)
		utils.InternalServError(w)
		return
	}

	// ** Récupération des infos user à partir du cookie **
	cookies := r.Cookies()
	var sessCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessCookie = cookie
		}
	}

	user, err := utils.GetUserInfoFromSess(sessCookie.Value)
	if err != nil {
		log.Printf("<profilhandler.go> Could not operate GetUserInfoFromSess: %v\n", err)
		utils.InternalServError(w)
		return
	}
	userPosts, err := utils.GetUserPosts(user.ID)
	if err != nil {
		log.Printf("<profilhandler.go> Could not operate GetUserPosts: %v\n", err)
		utils.InternalServError(w)
		return
	}

	// ** Renvoi des données dans le template **

	data := struct {
		Categories []models.Category
		User       models.User
		Posts      []models.LastPost
	}{
		Categories: categories,
		User:       user,
		Posts:      userPosts,
	}

	err = ProfilHtml.Execute(w, data)
	if err != nil {
		utils.InternalServError(w)
	}
}
