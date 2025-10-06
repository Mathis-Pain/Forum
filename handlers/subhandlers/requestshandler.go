package subhandlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/utils"
)

func RequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		action := r.FormValue("action")

		switch action {
		case "requestmod":

			db, err := sql.Open("sqlite3", "./data/forum.db")
			if err != nil {
				log.Printf("<messageactionshandler> Could not open database : %v\n", err)
				return
			}
			defer db.Close()

			username := r.FormValue("username")
			userID, _ := strconv.Atoi(r.FormValue("id"))
			url := "/profil"

			err = AskedToBeMod(db, userID)

			if err != nil {
				log.Print("Erreur dans l'envoi de la requête à l'administrateur, ", err)
				utils.InternalServError(w)
				return
			}

			log.Printf("REQUEST : L'utilisateur %s (n°%d) a demandé à rejoindre la modération.", username, userID)

			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	}

}
