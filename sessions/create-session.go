package sessions

import (
	"database/sql"
	"time"

	"github.com/Mathis-Pain/Forum/models"
)

// CreateSession crée une session pour un utilisateur et la sauvegarde
func CreateSession(userID int) (models.Session, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return models.Session{}, err
	}
	defer db.Close()

	sessionID, err := GenerateSessionID()
	if err != nil {
		return models.Session{}, err
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := SaveSession(db, session); err != nil {
		return models.Session{}, err
	}

	return session, nil
}
