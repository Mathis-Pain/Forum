package sessions

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSessionID génère un ID de session aléatoire
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
