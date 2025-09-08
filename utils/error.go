package utils

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
)

var (
	ErrorHtml = template.Must(template.ParseFiles("templates/error/error.html"))
)

// sendError est une fonction utilitaire pour rendre un template d'erreur
func sendError(w http.ResponseWriter, code int, name, message, full string) {
	data := models.HtmlError{
		Code:      code,
		ErrorName: name,
		Message:   message,
		ErrorFull: full,
	}

	// On écrit dans un buffer avant d'envoyer au client
	var buf bytes.Buffer
	if err := ErrorHtml.Execute(&buf, data); err != nil {
		// Si le template échoue, on renvoie un 500
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	// Envoie le header et le contenu
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}

// NotFoundHandler renvoie une erreur 404
func NotFoundHandler(w http.ResponseWriter) {
	sendError(w, http.StatusNotFound,
		"Page introuvable",
		"Désolé, la page que vous recherchez n'existe pas.",
		"404 - Not Found",
	)
}

// StatusBadRequest renvoie une erreur 400
func StatusBadRequest(w http.ResponseWriter) {
	sendError(w, http.StatusBadRequest,
		"Requête non prise en charge",
		"L'action que vous avez tenté d'effectuer n'est pas prise en charge.",
		"400 - Bad Request",
	)
}

// InternalServError renvoie une erreur 500
func InternalServError(w http.ResponseWriter) {
	sendError(w, http.StatusInternalServerError,
		"Erreur interne au serveur",
		"Le serveur a rencontré une erreur. Veuillez réessayer.",
		"500 - Internal Server Error",
	)
}

// MethodNotAllowedError renvoie une erreur 405
func MethodNotAllowedError(w http.ResponseWriter) {
	sendError(w, http.StatusMethodNotAllowed,
		"Méthode non autorisée",
		"L'accès à cette page n'est pas autorisé avec cette méthode HTML.",
		"405 - Method Not Allowed",
	)
}

// UnauthorizedError renvoie une erreur 401
func UnauthorizedError(w http.ResponseWriter) {
	sendError(w, http.StatusUnauthorized,
		"Connexion nécessaire",
		"Vous n'êtes pas autorisé.e à accéder à cette page. Veuillez vous connecter et réessayer.",
		"401 - Unauthorized",
	)
}

// ForbiddenError renvoie une erreur 403
func ForbiddenError(w http.ResponseWriter) {
	sendError(w, http.StatusForbidden,
		"Accès interdit",
		"Vous n'êtes pas autorisé.e à accéder à cette page.",
		"403 - Forbidden",
	)
}

// TimeOutError renvoie une erreur 408
func TimeOutError(w http.ResponseWriter) {
	sendError(w, http.StatusRequestTimeout,
		"Communication trop lente",
		"Le serveur a mis trop de temps à répondre à la requête.",
		"408 - Request Time Out",
	)
}
