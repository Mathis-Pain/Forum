package admin

import (
	"database/sql"
)

func CheckIfAdmin(username string) (bool, error) {
	var role int
	// ** Récupération du rôle**

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	sql := `SELECT role_id FROM user WHERE username = ?`
	row := db.QueryRow(sql, username)

	err = row.Scan(&role)
	if err != nil {
		return false, err
	}
	if role != 1 {
		return false, nil
	}

	return true, nil
}

func GetUserType(username string) (int, error) {
	var role int
	// ** Récupération du rôle**

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sql := `SELECT role_id FROM user WHERE username = ?`
	row := db.QueryRow(sql, username)

	err = row.Scan(&role)
	if err != nil {
		return 0, err
	}

	return role, nil
}
