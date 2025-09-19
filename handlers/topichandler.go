package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
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
	ID := utils.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<topichandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	topic, err := getdata.GetTopicInfo(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		log.Printf("<topichandler.go> Could not operate GetTopicInfo: %v\n", err)
		utils.InternalServError(w)
		return
	}

	categories, err := getdata.GetCatList()

	if err != nil {
		log.Printf("<cathandler.go> Could not operate GetCatList: %v\n", err)
		utils.InternalServError(w)
		return
	}

	userLoggedIn := false
	_, err = r.Cookie("session_id")
	if err == nil {
		userLoggedIn = true
	}

	data := struct {
		Topic      models.Topic
		Categories []models.Category
		LoginData  models.LoginData
		LogStatus  bool
	}{
		Topic:      topic,
		Categories: categories,
		LoginData:  models.LoginData{},
		LogStatus:  userLoggedIn,
	}

	err = TopicHtml.Execute(w, data)
	if err != nil {
		log.Printf("<topichandler.go> Could not execute template <topic.html> : %v\n", err)
		utils.InternalServError(w)
		return
	}
}
