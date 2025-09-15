package test

import "github.com/Mathis-Pain/Forum/models"

func TestMessage2() models.Message {
	var testmessage models.Message

	testmessage.TopicID = 1
	testmessage.MessageID = 1
	testmessage.Created = "10 / 10 / 2025"
	testmessage.Author = TestUser()
	testmessage.Likes = 0
	testmessage.Dislikes = 0
	testmessage.Content = "Ceci est le dernier message de ce topic"

	return testmessage
}
