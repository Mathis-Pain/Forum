package models

// Erreurs dans le formulaire d'inscription
type RegisterDataError struct {
	NameError  string
	EmailError string
	PassError  string
}
