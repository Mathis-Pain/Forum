package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Authentification(db *sql.DB, login string, password string) error {
	user, err := GetUserInfoFromLogin(db, login)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("<login.go> : login failed, User %v not found\n", login)
		return fmt.Errorf("incorrect password or username.")
	} else if err != nil {
		mylog := fmt.Errorf("(db) could not retrieve user infos from the database %w", err)
		log.Println("ERROR <login.go> l39 :", mylog)
		return mylog
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("<login.go> : password does not match")
		return fmt.Errorf("incorrect password or username.")
	}

	return nil
}
