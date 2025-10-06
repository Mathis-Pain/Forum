package sessions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Mathis-Pain/Forum/models"
)

// GetSession récupère une session depuis la DB
func GetSession(sessionID string) (models.Session, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return models.Session{}, err
	}
	defer db.Close()

	var session models.Session
	var dataJSON string

	err = db.QueryRow(`
		SELECT id, user_id, data, expires_at, created_at
		FROM sessions
		WHERE id = ?
	`, sessionID).Scan(&session.ID, &session.UserID, &dataJSON, &session.ExpiresAt, &session.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			//log.Print("<get-session.go> Erreur dans la récupération de session, aucune session trouvée :", err)
			return models.Session{}, err
		}
		return models.Session{}, err
	}

	if time.Now().After(session.ExpiresAt) {
		return models.Session{}, errors.New("session expired")
	}

	if err := json.Unmarshal([]byte(dataJSON), &session.Data); err != nil {
		log.Print("<get-session.go> Erreur dans la récupération de session :", err)
		return models.Session{}, err
	}

	return session, nil
}

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
