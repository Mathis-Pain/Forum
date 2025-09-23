package models

// Affichage des derniers messages sur la page d'accueil
// Struct Message avec un champ TopicName en plus
type LastPost struct {
	Message
	TopicName string
}

type Stats struct {
	TotalTopics   int
	LastMonthPost int
	LastUser      string
	TotalUsers    int
	Reported      int
}

type Likes struct {
	UserID     int
	LikedPosts []int
}
