package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var CatHtml = template.Must(template.New("categorie.html").Funcs(funcMap).ParseFiles(
	"templates/login.html",
	"templates/header.html",
	"templates/categorie.html",
	"templates/initpage.html",
))

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ID := subhandlers.GetPageID(r)
	if ID == 0 {
		utils.NotFoundHandler(w)
		return
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// --- Récupération des catégories ---

	category, err := getdata.GetCatDetails(db, ID)
	if err == sql.ErrNoRows {
		utils.NotFoundHandler(w)
		return
	} else if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la récupération de la catégorie : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// - Récupération des topics supplémentaires -

	stringID := strconv.Itoa(category.ID)

	rows, err := db.Query(`
	SELECT t.id, t.name
	FROM topic t
	WHERE EXISTS (
		SELECT 1
		FROM json_each(t.category_ids)
		WHERE json_each.value = ?
	)
`, stringID)

	if err != nil {
		log.Println("<cathandler.go> Erreur lors de la requête à la db pour les catégories supplémentaires: ", err)
		utils.InternalServError(w)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Topic
		err := rows.Scan(&t.TopicID, &t.Name)
		if err != nil {
			log.Printf("<cathandler.go> Erreur lors de la lecture des topics dans la db (catégories supplémentaires):  %v", err)
			utils.InternalServError(w)
			return
		}
		t.Messages, err = getdata.GetMessageList(db, t.TopicID)
		if err == sql.ErrNoRows {
			t.Messages = []models.Message{}
		} else if err != nil {
			log.Printf("<cathandler.go> Erreur lors de la lecture des messages dans la db (catégories supplémentaires):  %v", err)
			return
		}

		t.LastPost = len(t.Messages) - 1
		if t.LastPost < 0 {
			t.LastPost = 0
		}
		category.Topics = append(category.Topics, t)
		log.Printf("topic ajouté: %v", t)
	}

	// --- Construction du header ---

	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	for i := 0; i < len(category.Topics); i++ {
		category.Topics[i].Messages = getdata.FormatDateAllMessages(category.Topics[i].Messages)
	}

	// --- Gestion des erreurs de login ---

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution de GetSessionFromRequest: %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	var loginErr string
	if session.ID != "" {
		loginErr, err = getdata.GetLoginErr(session)
		if err != nil {
			logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution de GetLoginErr: %v\n", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
			return
		}
	}

	// --- Renvoi des données ---

	data := struct {
		PageName      string
		Category      models.Category
		Categories    []models.Category
		LoginErr      string
		CurrentUser   models.UserLoggedIn
		Notifications models.Notifications
	}{
		PageName:      category.Name,
		Category:      category,
		Categories:    categories,
		LoginErr:      loginErr,
		CurrentUser:   currentUser,
		Notifications: notifications,
	}

	err = CatHtml.Execute(w, data)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'exécution du template <categorie.html> : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

}
