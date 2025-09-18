package sessions

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

// ouvre juste la db pour sauvegarder la session via SaveSession
func SaveSessionToDB(session models.Session) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	return SaveSession(db, session)
}
