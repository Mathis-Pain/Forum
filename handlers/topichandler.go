package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var TopicHtml = template.Must(template.New("topic.html").ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/topic.html",
	"templates/initpage.html",
))

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	ID := subhandlers.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <topichandler.go> Erreur à l'ouverture de la base de données : %v", err)
		logs.AddLogsToDatabase(logMsg)
		return
	}
	defer db.Close()

	topic, err := getdata.GetTopicInfo(db, ID)

	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <topichandler.go> Erreur à l'exécution de GetTopicInfo: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	topic.TopicID = ID

	topic.Messages = getdata.FormatDateAllMessages(topic.Messages)

	categ, _ := getdata.GetCatDetails(db, topic.CatID)

	// Supprime le sujet et redirige vers la page d'accueil s'il ne contient aucun message (sécurité anti bug de la BDD)
	if len(topic.Messages) == 0 {
		ID := strconv.Itoa(topic.TopicID)
		subhandlers.DeleteTopicHandler(ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <topichandler.go> Erreur à l'exécution de GetSessionFromRequest: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	var loginErr string
	if session.ID != "" {
		loginErr, err = getdata.GetLoginErr(session)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <topichandler.go> Erreur à l'exécution de GetLoginErr: %v", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
	}

	_, allTopics, err := admin.GetAllTopics(categories, db)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <topichandler.go> Erreur dans la récupération de la liste des sujets : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	data := struct {
		PageName      string
		AllTopics     []models.Topic
		Topic         models.Topic
		CatName       string
		Categories    []models.Category
		LoginErr      string
		CurrentUser   models.UserLoggedIn
		Notifications models.Notifications
	}{
		PageName:      topic.Name,
		AllTopics:     allTopics,
		Topic:         topic,
		CatName:       categ.Name,
		Categories:    categories,
		LoginErr:      loginErr,
		CurrentUser:   currentUser,
		Notifications: notifications,
	}

	err = TopicHtml.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <topichandler.go> Erreur à l'exécution de template <topic.html> : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
}
