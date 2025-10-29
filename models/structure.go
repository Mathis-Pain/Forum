package models

// Catégories, sujets et messages du forum

type Category struct {
	ID          int
	Name        string
	Description string
	Topics      []Topic
}

type Topic struct {
	CatID    int
	AllCatID string
	TopicID  int
	Name     string
	Messages []Message
	LastPost int
	CatName  string
}

type Message struct {
	TopicID   int
	MessageID int
	Created   string
	Author    User
	Likes     int
	Dislikes  int
	Content   string
	Warning   int
}
