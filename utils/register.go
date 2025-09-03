package utils

import (
	"fmt"
	"log"
	"net/mail"
	"unicode"
)

func ValidPasswd(password string, confirmPassword string) error {
	if password != confirmPassword {
		mylog := fmt.Errorf("le mot de passe saisi est différent, merci d'entrer un mot de passe identique")
		log.Print(mylog)
		return mylog
	}
	if len(password) < 6 || len(password) >= 40 {
		mylog := fmt.Errorf("la longueur du mot de passe doit être comprise entre 6 et 40 caractères")
		log.Print(mylog)
		return mylog
	}

	nb := false
	maj := false

	for _, char := range password {
		if char >= '0' && char <= '9' {
			nb = true
		}
		if char >= 'A' && char <= 'Z' {
			maj = true
		}
	}

	if !maj {
		mylog := fmt.Errorf("le mot de passe doit comporter au moins une majuscule")
		log.Print(mylog)
		return mylog
	}

	if !nb {
		mylog := fmt.Errorf("le mot de passe doit comporter au moins un chiffre")
		log.Print(mylog)
		return mylog
	}

	for _, char := range password {
		if !unicode.IsPrint(char) {
			mylog := fmt.Errorf("ce caractère est invalide : %v : merci de le supprimer ou de le remplacer", char)
			log.Print(mylog)
			return mylog
		}
	}

	return nil
}

func ValidName(name string) error {
	if len(name) < 3 {
		mylog := fmt.Errorf("le nom d'utilisateur doit comporter au moins trois caractères")
		log.Print(mylog)
		return mylog
	}
	if len(name) >= 20 {
		mylog := fmt.Errorf("le nom d'utilisateur doit comporter moins de vingt caractères")
		log.Print(mylog)
		return mylog
	}
	for _, char := range name {
		if !unicode.IsPrint(char) {
			mylog := fmt.Errorf("ce caractère est invalide : %v : merci de le supprimer ou de le remplacer", char)
			log.Print(mylog)
			return mylog
		}
	}
	return nil
}

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		mylog := fmt.Errorf("l'adresse e-mail est invalide : merci de rentrer une adresse e-mail valide")
		log.Print(mylog)
		return mylog
	}
	return nil
}
