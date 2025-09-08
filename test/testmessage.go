package test

import "github.com/Mathis-Pain/Forum/models"

func TestMessage() models.Message {
	var testmessage models.Message

	testmessage.TopicID = 1
	testmessage.MessageID = 1
	testmessage.Created = "10 / 10 / 2025"
	testmessage.Author = 1
	testmessage.Likes = 0
	testmessage.Dislikes = 0
	testmessage.Content = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	return testmessage
}
