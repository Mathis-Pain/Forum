package sessions

import (
	"net/http"

	"github.com/Mathis-Pain/Forum/models"
)

func GetSessionFromRequest(r *http.Request) (models.Session, error) {
	var sessionID string
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return models.Session{}, nil
		}
		return models.Session{}, err
	}

	sessionID = cookie.Value
	session, err := GetSession(sessionID)
	if err != nil {
		return models.Session{}, nil
	}

	return session, nil
}
