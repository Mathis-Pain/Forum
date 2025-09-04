package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetUserInfoFromLogin(db *sql.DB, login string) (models.User, error) {
	sql := `SELECT id, password FROM user WHERE email = ? OR username = ?`
	row := db.QueryRow(sql, login, login)

	var user models.User
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserInfoFromID(db *sql.DB, ID int) (models.User, error) {
	sql := `SELECT username, profilpic FROM users WHERE id = ?`
	row := db.QueryRow(sql, ID)

	var user models.User
	err := row.Scan(&user.Username, &user.ProfilPic)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
