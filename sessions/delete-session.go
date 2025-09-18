package sessions

import "database/sql"

// DeleteSession supprime une session
func DeleteSession(sessionID string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}
