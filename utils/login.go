package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUserInfo(db *sql.DB, input string) (models.User, error) {
	sql := `SELECT ID, password FROM users WHERE email = ? OR username = ?`
	row := db.QueryRow(sql, input, input)

	var user models.User
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func Authentification(db *sql.DB, login string, password string) error {
	user, err := GetUserInfo(db, login)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("<login.go l.35 : login failed, User %v not found\n", login)
		return fmt.Errorf("incorrect password or username.")
	} else if err != nil {
		mylog := fmt.Errorf("could not recover user infos from the database %v", err)
		log.Println("ERROR <login.go> l39 :", mylog)
		return mylog
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("<login.go> l.45 : password does not match")
		return fmt.Errorf("incorrect password or username.")
	}

	return nil
}
