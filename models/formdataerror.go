package models

// Struct pour transmettre les erreurs d'inscription au template
type FormDataError struct {
	NameError  string
	EmailError string
	PassError  string
}
