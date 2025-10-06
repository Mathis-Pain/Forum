package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/Mathis-Pain/Forum/models"
)

// GenerateSessionID génère un ID de session aléatoire
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

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
		ExpiresAt: time.Now().Add(4 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := SaveSession(db, session); err != nil {
		return models.Session{}, err
	}

	return session, nil
}
