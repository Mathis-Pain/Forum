package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
	"github.com/Mathis-Pain/Forum/utils/postactions"
)

var AnswerMessage = template.Must(template.New("answermessage.html").Funcs(funcMap).ParseFiles(
	"templates/answermessage.html",
	"templates/create-topic.html",
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func MessageHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messagehandler.go> Erreur dans l'ouverture de la base de données : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()
	// recuperation des messages qui se trouve dans le topic
	ID := r.URL.Query().Get("topic_id")
	intID, err := strconv.Atoi(ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messagehandler.go> Erreur de convertion : ID du sujet invalide (%s)", ID)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	topic, err := getdata.GetTopicInfo(db, intID)
	topic.TopicID = intID
	// verifie si on trouve le topic concerné dans la base de donnée
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messagehandler.go> Erreur dans l'exécution de GetTopicInfo: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	// On charge les categories pour le header
	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messagehandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Empêche les utilisateurs bannis ou non enrigistrés d'accéder à la page
	if currentUser.UserType == 4 || currentUser.UserType == 0 {
		utils.ForbiddenError(w)
		return
	}

	// Récupère les informations du premier et du dernier message du topic pour afficher
	// les références
	topic.Messages = getdata.FormatDateAllMessages(topic.Messages)
	lastMessage := topic.Messages[len(topic.Messages)-1]
	firstMessage := topic.Messages[0]

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			utils.InternalServError(w)
			return
		}
		message := r.FormValue("new-message")
		if message == "" {
			utils.StatusBadRequest(w)
			return
		}

		postactions.NewPost(currentUser.ID, intID, message, "")
		// Redirection vers la page catégorie
		http.Redirect(w, r, fmt.Sprintf("/topic/%d", intID), http.StatusSeeOther)
		return
	}

	data := struct {
		Topic         models.Topic
		PageName      string
		LoginErr      string
		CurrentUser   models.UserLoggedIn
		Categories    []models.Category
		FirstMessage  models.Message
		LastMessage   models.Message
		Notifications models.Notifications
	}{
		Topic:         topic,
		PageName:      "Poster un message",
		LoginErr:      "",
		CurrentUser:   currentUser,
		Categories:    categories,
		FirstMessage:  firstMessage,
		LastMessage:   lastMessage,
		Notifications: notifications,
	}

	err = AnswerMessage.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <messagehandler.go> Could not execute template <answermessage.html>: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.NotFoundHandler(w)

	}
}
