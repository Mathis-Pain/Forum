package utils

import (
	"html/template"
	"net/http"
)

var (
	notFoundhtml          = template.Must(template.ParseFiles("templates/error/404.html"))
	statBadRequesthtml    = template.Must(template.ParseFiles("templates/error/400.html"))
	internalServErrorhtml = template.Must(template.ParseFiles("templates/error/500.html"))
)

// execute notFoundhtml error 404
func NotFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	notFoundhtml.Execute(w, nil)
}

// execute statBasRequesthtml error 400
func StatusBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	statBadRequesthtml.Execute(w, nil)
}

// execute internalServErrorhtml error 500
func InternalServError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	internalServErrorhtml.Execute(w, nil)
}
