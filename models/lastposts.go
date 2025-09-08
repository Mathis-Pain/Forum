package models

// Affichage des derniers messages sur la page d'accueil
// Struct Message avec un champ TopicName en plus
type LastPost struct {
	Message
	TopicName string
}
