package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	admin "github.com/Mathis-Pain/Forum/utils/adminfuncs"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

func BuildHeader(r *http.Request, w http.ResponseWriter, db *sql.DB) (models.Notifications, []models.Category, models.UserLoggedIn, error) {
	categories, err := CategoriesDropDownMenu()
	if err != nil && err != sql.ErrNoRows {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération de la liste des catégories :", err)
		logs.AddLogsToDatabase(logMsg)
		return models.Notifications{}, nil, models.UserLoggedIn{}, err
	}

	if r.Method == "POST" && r.FormValue("notif-action") != "" {
		stringID := r.FormValue("notifID")
		notifID, _ := strconv.Atoi(stringID)
		switch r.FormValue("notif-action") {
		case "markread":
			err := logs.MarkAsRead(notifID)
			if err != nil {
				utils.InternalServError(w)
				return models.Notifications{}, nil, models.UserLoggedIn{}, err
			}
		case "delete":
			err := logs.DeleteNotif(notifID)
			if err != nil {
				utils.InternalServError(w)
				return models.Notifications{}, nil, models.UserLoggedIn{}, err
			}
		default:
			utils.StatusBadRequest(w)
			return models.Notifications{}, nil, models.UserLoggedIn{}, err
		}
	}

	var currentUser models.UserLoggedIn

	currentUser.LogStatus = CheckLogStatus(r)

	// Si un utilisateur est en ligne, récupère son nom pour l'afficher à droite + son ID pour le profil
	if currentUser.LogStatus {
		// Récupère le pseudo et l'ID de l'utilisateur si un utilisateur est en ligne
		currentUser.Username, currentUser.ID, err = utils.GetUserNameAndIDByCookie(r, db)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération des données utilisateur :", err)
			logs.AddLogsToDatabase(logMsg)
			return models.Notifications{}, categories, currentUser, err
		}
		currentUser.UserType, err = admin.GetUserType(currentUser.Username)
		if err != nil {
			logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération des données utilisateur :", err)
			logs.AddLogsToDatabase(logMsg)
			return models.Notifications{}, categories, currentUser, err
		}
	}

	notifications, err := logs.DisplayNotifications(currentUser.ID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération des notifications :", err)
		logs.AddLogsToDatabase(logMsg)
		return models.Notifications{}, categories, currentUser, err
	}

	return notifications, categories, currentUser, nil

}

// Vérifie si un utilisateur est connecté
func CheckLogStatus(r *http.Request) bool {
	userLoggedIn := false
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <buildheader.go> Erreur dans l'exécution de GetSessionFromRequest: %v", err)
		logs.AddLogsToDatabase(logMsg)
		return false
	}
	if session.UserID != 0 {
		userLoggedIn = true
	}
	return userLoggedIn
}

// // Récupère le pseudo et l'ID de l'utilisateur si un utilisateur est en ligne
// func getUserNameAndID(r *http.Request, db *sql.DB) (string, int, error) {
// 	// Récupère l'ID de l'utilisateur via sa session
// 	cookie, err := r.Cookie("session_id")
// 	if err != nil {
// 		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération du cookie : ", err)
// 		return "", 0, err
// 	}
// 	session, err := sessions.GetSession(cookie.Value)
// 	if err != nil && err != sql.ErrNoRows {
// 		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération de session : ", err)
// 		return "", 0, err
// 	}

// 	// Récupère le pseudo de l'utilisateur
// 	sqlQuery := `SELECT username FROM user WHERE id = ?`
// 	row := db.QueryRow(sqlQuery, session.UserID)

// 	var username string

// 	err = row.Scan(&username)
// 	if err != nil && err != sql.ErrNoRows {
// 		return "", 0, err
// 	}

// 	return username, session.UserID, nil
// }

// Fabrique la liste des catégories pour le menu déroulant
func CategoriesDropDownMenu() ([]models.Category, error) {
	categories, err := getdata.GetCatList()
	if err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}
