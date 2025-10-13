package models

// Données d'un utilisateur
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	ProfilPic string
	Status    string
}

type UserLoggedIn struct {
	ID        int
	Username  string
	ProfilPic string
	LogStatus bool
	UserType  int
}
