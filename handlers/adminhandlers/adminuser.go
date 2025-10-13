package adminhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
	"github.com/Mathis-Pain/Forum/utils/logs"
)

// MARK: Modifier un utilisateur
func UserEditHandler(r *http.Request, users []models.User, currentUser models.UserLoggedIn) error {
	// Récupère l'ID de l'utilisateur dans le formulaire
	stringID := r.FormValue("userID")
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return err
	}

	// Repère l'utilisateur à modifier via son ID
	var user models.User
	for _, current := range users {
		if current.ID == ID {
			user = current
			break
		}
	}

	// Récupère le nom et le statut dans le formulaire
	username := r.FormValue("username")
	status := r.FormValue("status")
	previousName := r.FormValue("previous")
	previousStatus := r.FormValue("current")

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminuser.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	// Vérifie si le statut a été modifié
	nameChange := false
	statusChange := false
	if status != previousStatus && status != "" {
		statusChange = true
	}

	// Vérifie si le nom a été modifié et si le nouveau nom n'est pas déjà pris
	uniqueName := CheckUniqueName(username, db)
	logName := username
	if username != previousName && uniqueName {
		logName := username
		logName += " (anciennement " + previousName + ")"
		nameChange = true
	}

	if !nameChange && !statusChange {
		return nil
	}

	// Ajout des logs et des notifications
	notifMsg := fmt.Sprintf("Votre compte a été modifié par un administrateur (%s).", currentUser.Username)
	logMsg := fmt.Sprintf("ADMIN : L'utilisateur %s a été modifié par %s.", logName, currentUser.Username)

	if nameChange {
		// Si le pseudo a été modifié
		if username != previousName {
			user.Username = username
			notifMsg += fmt.Sprintf(" Votre nom d'utilisateur a été changé en %s.", username)
			logMsg += fmt.Sprintf(" Son nom d'utilisateur a été changé en %s.", username)
		}
	}

	// Si le statut a été modifié
	if statusChange {
		user.Status = status
		if status == "Membre " {
			status = "Membre"
		}
		notifMsg += fmt.Sprintf(" Vous avez changé de statut et êtes maintenant %s.", status)
		logMsg += fmt.Sprintf(" Son statut a été modifié en %s.", status)
	}

	// Met à jour l'utilisateur et son rôle à partir de son ID
	sqlUpdate := `UPDATE user SET username = ?, role_id = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, getdata.CodeUserStatus(user.Status), user.ID)
	if err != nil {
		return err
	}

	logs.AddNotificationToDatabase("ADMIN", ID, 0, notifMsg)
	logs.AddLogsToDatabase(logMsg)

	return nil
}

// MARK: Bannir un utilisateur
func BanUserHandler(stringID string) error {
	// Récupère l'ID de l'utilisateur à bannir
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminuser.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	notificationMessage := "Votre compte a été banni par un administrateur. Vous ne pouvez plus poster ni répondre à des messages."
	logs.AddNotificationToDatabase("ADMIN", ID, 0, notificationMessage)

	// Met à jour l'utilisateur avec le statut BANNI (4)
	sqlUpdate := `UPDATE user SET role_id = 4 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}

// MARK: Débannir un utilisateur
func UnbanUserHandler(stringID string) error {
	// Récupération de l'ID
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminuser.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	notificationMessage := "Votre compte a été débanni par un administrateur. Vous pouvez à nouveau poster ou répondre à des messages."
	logs.AddNotificationToDatabase("ADMIN", ID, 0, notificationMessage)

	// Met à jour l'utilisateur avec le statut MEMBRE (3)
	sqlUpdate := `UPDATE user SET role_id = 3 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}

// MARK: Supprimer un utilisateur
func DeleteUserHandler(stringID string) error {
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminuser.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	// Supprime l'utilisateur de la base de données
	sqlUpdate := `DELETE FROM user WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		return err
	}

	return nil
}

// MARK: Promouvoir modérateur
func PromoteToMod(userID int) error {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <adminuser.go> Erreur à l'ouverture de la base de données :", err)
		logs.AddLogsToDatabase(logMsg)
		return err
	}
	defer db.Close()

	sqlUpdate := `UPDATE user SET role_id = 2 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func CheckUniqueName(username string, db *sql.DB) bool {
	sqlQuery := `SELECT id FROM user WHERE username = ?`
	row := db.QueryRow(sqlQuery, username)

	err := row.Scan(&username)
	if err == sql.ErrNoRows {
		return true
	} else {
		return false
	}

}
