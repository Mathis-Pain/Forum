package utils

import (
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
)

var (
	ErrorHtml = template.Must(template.ParseFiles("templates/error/error.html"))
)

// Erreur 404 - Not Found
func NotFoundHandler(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      404,
		ErrorName: "Page introuvable",
		Message:   "Désolé, la page que vous recherchez n'existe pas.",
		ErrorFull: "404 - Not Found",
	}
	w.WriteHeader(http.StatusNotFound)
	ErrorHtml.Execute(w, data)
}

// Erreur 400 - Bad Request
func StatusBadRequest(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      400,
		ErrorName: "Requête non prise en charge",
		Message:   "L'action que vous avez tenté d'effectuer n'est pas prise en charge.",
		ErrorFull: "400 - Bad Request",
	}
	w.WriteHeader(http.StatusBadRequest)
	ErrorHtml.Execute(w, data)
}

// Erreur 500 - Erreur Serveur
func InternalServError(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      500,
		ErrorName: "Erreur interne au serveur",
		Message:   "Le serveur a rencontré une erreur. Veuillez réessayer.",
		ErrorFull: "500 - Internal Servor Error",
	}
	w.WriteHeader(http.StatusInternalServerError)
	ErrorHtml.Execute(w, data)
}

// Erreur 405 - Méthode invalide
func MethodNotAllowedError(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      405,
		ErrorName: "Méthode non autorisée",
		Message:   "L'accès à cette page n'est pas autorisé avec cette méthode HTML.",
		ErrorFull: "405 - Method Not Allowed",
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	ErrorHtml.Execute(w, data)
}

// Erreur 401 - Connexion nécessaire
func UnauthorizedError(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      401,
		ErrorName: "Connexion nécessaire",
		Message:   "Vous n'êtes pas autorisé.e à accéder à cette page. Veuillez vous connecter et réessayer.",
		ErrorFull: "401 - Unauthorized",
	}
	w.WriteHeader(http.StatusUnauthorized)
	ErrorHtml.Execute(w, data)
}

// Erreur 403 - Non autorisé
func ForbiddenError(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      403,
		ErrorName: "Accès interdit",
		Message:   "Vous n'êtes pas autorisé.e à accéder à cette page.",
		ErrorFull: "403 - Forbidden",
	}
	w.WriteHeader(http.StatusForbidden)
	ErrorHtml.Execute(w, data)
}

// Erreur 408 - Time Out
func TimeOutError(w http.ResponseWriter) {
	data := models.HtmlError{
		Code:      408,
		ErrorName: "Communication trop lente",
		Message:   "Le serveur a mis trop de temps à répondre à la requête.",
		ErrorFull: "408 - Request Time Out",
	}
	w.WriteHeader(http.StatusRequestTimeout)
	ErrorHtml.Execute(w, data)
}
