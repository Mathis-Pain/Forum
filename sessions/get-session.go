package sessions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
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
			log.Print("<get-session.go> Erreur dans la récupération de session, aucune session trouvée :", err)
			return models.Session{}, errors.New("session not found")
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
