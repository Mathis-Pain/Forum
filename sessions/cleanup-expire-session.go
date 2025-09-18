package sessions

import "database/sql"

// CleanupExpiredSessions supprime les sessions expirées
func CleanupExpiredSessions() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP")
	return err
}
