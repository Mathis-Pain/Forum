package utils

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func GetUserInfoFromLogin(db *sql.DB, login string) (models.User, error) {
	sql := `SELECT ID, password FROM users WHERE email = ? OR username = ?`
	row := db.QueryRow(sql, login, login)

	var user models.User
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
