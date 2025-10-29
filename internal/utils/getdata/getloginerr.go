package getdata

import (
	"github.com/Mathis-Pain/Forum/internal/models"
	"github.com/Mathis-Pain/Forum/internal/sessions"
)

func GetLoginErr(session models.Session) (string, error) {
	var resultErr string
	if sessionErr, ok := session.Data["LoginErr"].(string); ok {
		resultErr = sessionErr
		err := sessions.DeleteSession(session.ID)
		if err != nil {
			return "", err
		}
	} else {
		resultErr = ""
	}

	return resultErr, nil
}
