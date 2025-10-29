package models

// Erreurs HTML
type HtmlError struct {
	Code      int
	ErrorName string
	Message   string
	PageName  string
}

// Erreurs dans le formulaire d'inscription
type RegisterDataError struct {
	NameError  string
	EmailError string
	PassError  string
}
