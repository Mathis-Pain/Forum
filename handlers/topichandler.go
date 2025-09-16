package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
)

var TopicHtml = template.Must(template.New("topic.html").ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/topic.html",
	"templates/initpage.html",
	"templates/reponsebox.html",
))

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		utils.NotFoundHandler(w)
		return
	}

	ID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<topichandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	topic, err := utils.GetTopicInfo(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		log.Printf("<topichandler.go> Could not operate GetTopicInfo: %v\n", err)
		utils.InternalServError(w)
		return
	}

	categories, err := utils.GetCatList()

	if err != nil {
		log.Printf("<cathandler.go> Could not operate GetCatList: %v\n", err)
		utils.InternalServError(w)
		return
	}

	data := struct {
		Topic      models.Topic
		Categories []models.Category
		LoginData  models.LoginData
	}{
		Topic:      topic,
		Categories: categories,
		LoginData:  models.LoginData{},
	}

	err = TopicHtml.Execute(w, data)
	if err != nil {
		log.Printf("<topichandler.go> Could not execute template <topic.html> : %v\n", err)
		utils.InternalServError(w)
		return
	}
}
