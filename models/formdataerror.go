package models

// Struct pour transmettre les erreurs au template
type FormDataError struct {
	NameError  string
	EmailError string
	PassError  string
}
