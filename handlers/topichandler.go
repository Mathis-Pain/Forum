package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/utils"
)

// var TopicHtml = template.Must(template.ParseFiles("templates/sujet.html"))

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Printf("<homehandler.go> Could not open database : %v\n", err)
		return
	}
	defer db.Close()

	parts := strings.Split(r.URL.Path, "/")
	path := parts[len(parts)-1]

	if !strings.Contains(path, "t") {
		utils.NotFoundHandler(w)
	}

	// topicID := strings.Trim(path, "t")
	// ID, err := strconv.Atoi(topicID)

	// if err != nil {
	// 	utils.InternalServError(w)
	// }

	// topic, err := utils.GetTopicInfo(db, ID)

	// if err == sql.ErrNoRows {
	// 	utils.NotFoundHandler(w)
	// } else if err != nil {
	// 	utils.InternalServError(w)
	// }

	// data := struct {
	// 	Topic models.Topic
	// }{
	// 	Topic: topic,
	// }

	// err = TopicHtml.Execute(w, data)
	// if err != nil {
	// 	log.Printf("Erreur lors de l'exécution du template <sujet.html> : %v\n", err)
	// 	utils.NotFoundHandler(w)
	// }

}
