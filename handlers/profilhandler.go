package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

var ProfilHtml = template.Must(template.New("profil.html").ParseFiles(
	"templates/profil.html",
	"templates/login.html",
	"templates/header.html",
	"templates/initpage.html"))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <profilhandler.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Création du header
	notifications, categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur dans la construction du header : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ** Récupération des infos de l'utilisateur **
	user, err := getUserProfile(currentUser.Username, db)
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <profilhandler.go> Erreur dans la récupération des données utilisateur :", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère la liste complète des messages postés par l'utilisateur
	userPosts, err := utils.GetUserPosts(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserPosts: %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Récupère la liste des sujets likés et dislikés par l'utilisateur
	likedPosts, err := utils.GetUserLikes(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserLikes : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	dislikedPosts, err := utils.GetUserDislikes(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserDislikes : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Affiche la liste des sujets ouvert par l'utilisateur
	myTopics, err := utils.GetUserTopics(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserTopics : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// Formatage de la date
	var currentTopic []models.Message
	for i := 0; i < len(myTopics); i++ {
		currentTopic = append(currentTopic, models.Message{})
		currentTopic[i].Created = myTopics[i].Created
		currentTopic = getdata.FormatDateAllMessages(currentTopic)
		myTopics[i].Created = currentTopic[i].Created
	}

	// ** Renvoi des données dans le template **
	pageName := fmt.Sprintf("Voir mon profil : %s", user.Username)

	data := struct {
		PageName      string
		User          models.User
		MyPosts       []models.LastPost
		LikedPosts    []models.LastPost
		DislikedPosts []models.LastPost
		Topics        []models.LastPost
		LoginErr      string
		Categories    []models.Category
		CurrentUser   models.UserLoggedIn
		Notifications models.Notifications
	}{
		PageName:      pageName,
		User:          user,
		MyPosts:       userPosts,
		LikedPosts:    likedPosts,
		DislikedPosts: dislikedPosts,
		Topics:        myTopics,
		LoginErr:      "",
		Categories:    categories,
		CurrentUser:   currentUser,
		Notifications: notifications,
	}

	err = ProfilHtml.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		utils.InternalServError(w)
	}
}

func getUserProfile(username string, db *sql.DB) (models.User, error) {
	var user models.User

	sql := `SELECT id, username, email, profilpic, role_id FROM user WHERE username = ?`
	row := db.QueryRow(sql, username)

	var role string

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilPic, &role)
	if err != nil {
		return models.User{}, err
	}

	user.Status = getdata.SetUserStatus(role)

	return user, nil
}
