package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

var TopicHtml = template.Must(template.New("topic.html").ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/topic.html",
	"templates/initpage.html",
	"templates/reponsebox.html",
))

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	ID := subhandlers.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("ERREUR : <topichandler.go> Erreur dans l'ouverture de la base de données : %v\n", err)
		return
	}
	defer db.Close()

	topic, err := getdata.GetTopicInfo(db, ID)

	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		log.Printf("ERREUR : <topichandler.go> Erreur dans l'exécution de GetTopicInfo: %v\n", err)
		utils.InternalServError(w)
		return
	}

	topic.TopicID = ID

	// Supprime le sujet et redirige vers la page d'accueil s'il ne contient aucun message (sécurité anti bug de la BDD)
	if len(topic.Messages) == 0 {
		ID := strconv.Itoa(topic.TopicID)
		subhandlers.DeleteTopicHandler(ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		log.Printf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v\n", err)
		utils.InternalServError(w)
		return
	}

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Printf("ERREUR : <topichandler.go> Erreur dans l'exécution de GetSessionFromRequest: %v\n", err)
		utils.InternalServError(w)
		return
	}
	var loginErr string
	if session.ID != "" {
		loginErr, err = getdata.GetLoginErr(session)
		if err != nil {
			log.Printf("ERREUR : <topichandler.go> Erreur dans l'exécution de GetLoginErr: %v\n", err)
			utils.InternalServError(w)
			return
		}
	}

	_, allTopics, err := admin.GetAllTopics(categories, db)
	if err != nil {
		log.Print("ERREUR : <topichandler.go> Erreur dans la récupération de la liste des sujets : ", err)
		utils.InternalServError(w)
		return
	}

	data := struct {
		PageName    string
		AllTopics   []models.Topic
		Topic       models.Topic
		Categories  []models.Category
		LoginErr    string
		CurrentUser models.UserLoggedIn
	}{
		PageName:    topic.Name,
		AllTopics:   allTopics,
		Topic:       topic,
		Categories:  categories,
		LoginErr:    loginErr,
		CurrentUser: currentUser,
	}

	err = TopicHtml.Execute(w, data)
	if err != nil {
		log.Printf("ERREUR : <topichandler.go> Erreur dans l'exécution de template <topic.html> : %v\n", err)
		utils.InternalServError(w)
		return
	}
}
