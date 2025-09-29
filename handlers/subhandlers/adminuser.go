package subhandlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// Fonction pour modifier un utilisateur (nom et statut)
func UserEditHandler(r *http.Request, users []models.User) error {
	// Récupère l'ID de l'utilisateur dans le formulaire
	stringID := r.FormValue("userID")
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<adminuser.go> Erreur dans la récupération de l'ID utilisateur : ", err)
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

	if username != "" {
		user.Username = username
	}
	if status != "" {
		user.Status = status
	}

	// Ouverture de la base de données
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<adminuser.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Met à jour l'utilisateur et son rôle à partir de son ID
	sqlUpdate := `UPDATE user SET username = ?, role_id = ? WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, getdata.CodeUserStatus(user.Status), user.ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// Fonction pour bannir un utilisateur
func BanUserHandler(stringID string) error {
	// Récupère l'ID de l'utilisateur à bannir
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<adminuser.go> Erreur dans la récupération de l'utilisateur à bannir")
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<adminuser.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Met à jour l'utilisateur avec le statut BANNI (4)
	sqlUpdate := `UPDATE user SET role_id = 4 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// Fonction pour "libérer" un utilisateur
func UnbanUserHandler(stringID string) error {
	// Récupération de l'ID
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<adminuser.go> Erreur dans la récupération de l'utilisateur à débannir")
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<adminuser.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Met à jour l'utilisateur avec le statut MEMBRE (3)
	sqlUpdate := `UPDATE user SET role_id = 3 WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

// Fonction pour supprimer un utilisateur
func DeleteUserHandler(stringID string) error {
	ID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print("<adminuser.go> Erreur dans la récupération de l'utilisateur à supprimer", err)
		return err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<adminuser.go> Erreur à l'ouverture de la base de données :", err)
		return err
	}
	defer db.Close()

	// Supprime l'utilisateur de la base de données
	sqlUpdate := `DELETE FROM user WHERE id = ?`
	stmt, err := db.Prepare(sqlUpdate)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
