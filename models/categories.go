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
	TopicID  int
	Name     string
	Messages []Message
}

type Message struct {
	TopicID   int
	MessageID int
	Created   string
	Author    int
	Likes     int
	Dislikes  int
	Content   string
}
