package models

// Erreurs HTML
type HtmlError struct {
	Code      int
	ErrorName string
	Message   string
	PageName  string
}
