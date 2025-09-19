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
)

var ProfilHtml = template.Must(template.New("profil.html").ParseFiles(
	"templates/profil.html",
	"templates/login.html",
	"templates/header.html",
	"templates/initpage.html"))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<profilhandler.go> Erreur à l'ouverture de la base de données :", err)
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Création du header
	categories, currentUser, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		log.Printf("<cathandler.go> Erreur dans la construction du header : %v\n", err)
		utils.InternalServError(w)
		return
	}

	// ** Récupération des infos de l'utilisateur **
	user, err := getUserProfile(currentUser.Username, db)
	if err != nil {
		log.Println("<profilhandler.go> Erreur dans la récupération des données utilisateur :", err)
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
	pageName := fmt.Sprintf("Voir mon profil : %s", user.Username)

	data := struct {
		PageName    string
		User        models.User
		Posts       []models.LastPost
		LoginData   models.LoginData
		Categories  []models.Category
		CurrentUser models.UserLoggedIn
	}{
		PageName:    pageName,
		User:        user,
		Posts:       userPosts,
		LoginData:   models.LoginData{},
		Categories:  categories,
		CurrentUser: currentUser,
	}

	err = ProfilHtml.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		utils.InternalServError(w)
	}
}

func getUserProfile(username string, db *sql.DB) (models.User, error) {
	var user models.User

	sql := `SELECT id, username, email, profilpic FROM user WHERE username = ?`
	row := db.QueryRow(sql, username)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilPic)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
