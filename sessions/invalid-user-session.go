package sessions

import "database/sql"

// InvalidateUserSessions supprime toutes les sessions d’un utilisateur
func InvalidateUserSessions(userID int) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
