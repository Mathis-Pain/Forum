package test

import "github.com/Mathis-Pain/Forum/models"

func TestLastPost() models.LastPost {
	var testlastpost models.LastPost

	testlastpost.TopicName = "Test"
	testlastpost.TopicID = 1
	testlastpost.MessageID = 1
	testlastpost.Created = "10 / 10 / 2025"
	testlastpost.Author = 1
	testlastpost.Likes = 0
	testlastpost.Dislikes = 0
	testlastpost.Content = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	return testlastpost
}
