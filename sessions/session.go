package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

var (
	sessions = make(map[string]Session)
	mu       sync.Mutex
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(data map[string]interface{}) (string, error) {
	id, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	mu.Lock()
	defer mu.Unlock()

	sessions[id] = Session{
		ID:        id,
		Data:      data,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expiration après 24h
	}
	return id, nil
}

func GetSession(id string) (Session, error) {
	mu.Lock()
	defer mu.Unlock()

	session, exists := sessions[id]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessions, id)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func DeleteSession(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, id)
}

// Nettoyage périodique des sessions expirées (à appeler dans main.go)
func CleanupExpiredSessions() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for id, session := range sessions {
		if now.After(session.ExpiresAt) {
			delete(sessions, id)
		}
	}
}
