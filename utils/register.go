package utils

import (
	"fmt"
	"net/mail"
	"unicode"
)

func ValidPasswd(password string, confirmPassword string) string {
	if password != confirmPassword {
		mylog := "Le mot de passe et sa confirmation sont différents. Merci d'entrer des mots de passe identiques"
		return mylog
	}
	if len(password) < 6 || len(password) >= 40 {
		mylog := "La longueur du mot de passe doit être comprise entre 6 et 40 caractères"
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
		mylog := "Le mot de passe doit comporter au moins une majuscule"
		return mylog
	}

	if !nb {
		mylog := "Le mot de passe doit comporter au moins un chiffre"
		return mylog
	}

	for _, char := range password {
		if !unicode.IsPrint(char) {
			mylog := fmt.Sprintf("Ce caractère est invalide : %v, merci de le supprimer ou de le remplacer", char)
			return mylog
		}
	}

	return ""
}

func ValidName(name string) string {
	if len(name) < 3 {
		mylog := "Le nom d'utilisateur doit comporter au moins trois caractères"
		return mylog
	}
	if len(name) >= 20 {
		mylog := "Le nom d'utilisateur doit comporter moins de vingt caractères"
		return mylog
	}
	for _, char := range name {
		if !unicode.IsPrint(char) {
			mylog := fmt.Sprintf("Ce caractère est invalide : %v : merci de le supprimer ou de le remplacer", char)
			return mylog
		}
	}
	return ""
}

func ValidEmail(email string) string {
	_, err := mail.ParseAddress(email)
	if err != nil {
		mylog := "L'adresse e-mail est invalide : merci de rentrer une adresse e-mail valide"

		return mylog
	}
	return ""
}
