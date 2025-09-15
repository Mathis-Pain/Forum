package test

import "github.com/Mathis-Pain/Forum/models"

func TestTopic() models.Topic {
	var testtopic models.Topic

	testtopic.CatID = 1
	testtopic.TopicID = 1
	testtopic.Name = "Test"

	testtopic.Messages = append(testtopic.Messages, TestMessage())
	testtopic.Messages = append(testtopic.Messages, TestMessage2())

	testtopic.LastPost = len(testtopic.Messages) - 1

	return testtopic
}
