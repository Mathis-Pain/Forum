package test

import "github.com/Mathis-Pain/Forum/models"

func TestTopic() models.Topic {
	var testtopic models.Topic

	testtopic.CatID = 1
	testtopic.TopicID = 1
	testtopic.Name = "Test"

	testtopic.Messages = append(testtopic.Messages, TestMessage())

	return testtopic
}
