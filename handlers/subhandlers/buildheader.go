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
	// MARK: Liste des catégories
	categories, err := CategoriesDropDownMenu()
	if err != nil && err != sql.ErrNoRows {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération de la liste des catégories :", err)
		logs.AddLogsToDatabase(logMsg)
		return models.Notifications{}, nil, models.UserLoggedIn{}, err
	}

	// MARK: Utilisateur en ligne
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

	// MARK: Liste notifications
	notifications, err := logs.DisplayNotifications(currentUser.ID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <buildheader.go> Erreur dans la récupération des notifications :", err)
		logs.AddLogsToDatabase(logMsg)
		return models.Notifications{}, categories, currentUser, err
	}

	// MARK: Actions notifications
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

// Fabrique la liste des catégories pour le menu déroulant
func CategoriesDropDownMenu() ([]models.Category, error) {
	categories, err := getdata.GetCatList()
	if err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}
