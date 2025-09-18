package sessions

import (
	"database/sql"
	"encoding/json"

	"github.com/Mathis-Pain/Forum/models"
)

// chemin d'acces a la db
const dbPath = "/data/forum.db"

// SaveSession sauvegarde ou met à jour une session
// Transforme les données Go en text Json pour pouvoir les sotcker dans la db
func SaveSession(db *sql.DB, session models.Session) error {
	dataJSON, err := json.Marshal(session.Data)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR REPLACE INTO sessions
		(id, user_id, data, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, session.ID, session.UserID, string(dataJSON), session.ExpiresAt, session.CreatedAt)

	return err
}
