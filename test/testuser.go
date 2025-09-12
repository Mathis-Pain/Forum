package test

import "github.com/Mathis-Pain/Forum/models"

func TestUser() models.User {
	var testuser models.User

	testuser.Username = "Moi"
	testuser.ID = 1
	testuser.Status = "ADMIN"
	testuser.ProfilPic = "https://i.imgur.com/gvmTSWn.png"

	return testuser
}
