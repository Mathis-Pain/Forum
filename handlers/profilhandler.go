package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
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

	user, userPosts, likedPosts, dislikedPosts, myTopics, err := GetProfileInfo(currentUser, db)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	if r.Method == "POST" {
		action := r.FormValue("action")

		switch action {
		case "editprofil":
			UpdateUserProfil(r, user, db)
		case "requestmod":
			subhandlers.RequestMod(db, user)
		}

		url := "/profil"
		http.Redirect(w, r, url, http.StatusSeeOther)

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

func UpdateUserProfil(r *http.Request, user models.User, db *sql.DB) error {
	newname := r.FormValue("newname")
	newimg := r.FormValue("newimg")
	newmail := r.FormValue("newmail")

	nameChanged, imgChanged, mailChanged := false, false, false

	if user.Username != newname {
		// log.Print("Le nom d'utilisateur a été modifié. nameChanged = true")
		nameChanged = true
	}

	if user.ProfilPic != newimg {
		// log.Print("L'image de profil a été modifiée. imgChanged = true")
		imgChanged = true
	}

	if isExt, err := ExternalUser(user.ID, db); err == nil {
		if !isExt && newmail != user.Email {
			// log.Print("L'adresse mail a été modifiée. mailChanged = true")
			mailChanged = true
		}
	}

	logMsg := ""

	if !nameChanged && !imgChanged && !mailChanged {
		// log.Print("Aucune modification effectuée.")
		return nil
	} else {
		// log.Print("Utilisateur modifié.")
		logMsg = "USER : L'utilisateur "
	}

	if nameChanged {
		logMsg += fmt.Sprintf("%s (anciennement %s) a modifié son profil.", newname, user.Username)
		user.Username = newname
	} else {
		logMsg += fmt.Sprintf("%s a modifié son profil.", user.Username)
	}

	if mailChanged {
		logMsg += " Son adresse email a été modifiée."
		user.Email = newmail
	}

	if imgChanged {
		logMsg += "Sa photo de profil a été modifiée."
		user.ProfilPic = newimg
	}

	err := UpdateUserDatabase(user, db)
	if err != nil {
		log.Print("Erreur :", err)
		return err
	}

	logs.AddLogsToDatabase(logMsg)
	return nil
}

func UpdateUserDatabase(user models.User, db *sql.DB) error {
	sqlUpdate := `UPDATE user SET username = ?, email = ?, profilpic = ? WHERE id = ?`
	_, err := db.Exec(sqlUpdate, user.Username, user.Email, user.ProfilPic, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func ExternalUser(userID int, db *sql.DB) (bool, error) {
	var googleID string
	sqlQuery := `SELECT google_id FROM user WHERE id = ?`
	row := db.QueryRow(sqlQuery, userID)

	err := row.Scan(&googleID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func GetProfileInfo(currentUser models.UserLoggedIn, db *sql.DB) (models.User, []models.LastPost, []models.LastPost, []models.LastPost, []models.LastPost, error) {
	// Récupération des infos de l'utilisateur **
	user, err := getUserProfile(currentUser.Username, db)
	if err != nil {
		logMsg := fmt.Sprintln("ERREUR : <profilhandler.go> Erreur dans la récupération des données utilisateur :", err)
		logs.AddLogsToDatabase(logMsg)
		return models.User{}, nil, nil, nil, nil, err
	}

	// Récupère la liste complète des messages postés par l'utilisateur
	userPosts, err := utils.GetUserPosts(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserPosts: %v", err)
		logs.AddLogsToDatabase(logMsg)
		return user, nil, nil, nil, nil, err
	}

	// Récupère la liste des sujets likés et dislikés par l'utilisateur
	likedPosts, err := utils.GetUserLikes(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserLikes : %v", err)
		logs.AddLogsToDatabase(logMsg)
		return user, userPosts, nil, nil, nil, err
	}

	dislikedPosts, err := utils.GetUserDislikes(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserDislikes : %v", err)
		logs.AddLogsToDatabase(logMsg)
		return user, userPosts, likedPosts, nil, nil, err
	}

	// Affiche la liste des sujets ouvert par l'utilisateur
	myTopics, err := utils.GetUserTopics(user.ID)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <profilhandler.go> Erreur à l'exécution de GetUserTopics : %v", err)
		logs.AddLogsToDatabase(logMsg)
		return user, userPosts, likedPosts, dislikedPosts, nil, err
	}

	// Formatage de la date
	var currentTopic []models.Message
	for i := 0; i < len(myTopics); i++ {
		currentTopic = append(currentTopic, models.Message{})
		currentTopic[i].Created = myTopics[i].Created
		currentTopic = getdata.FormatDateAllMessages(currentTopic)
		myTopics[i].Created = currentTopic[i].Created
	}

	return user, userPosts, likedPosts, dislikedPosts, myTopics, nil
}
