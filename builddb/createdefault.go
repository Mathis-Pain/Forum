package builddb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func DefaultDatabase(db *sql.DB) {
	// Check if the passed-in connection is valid
	if db == nil {
		log.Print("<createdefault.go> Database connection is nil.")
		return
	}

	// This helper function centralizes the logic for executing a query and handling errors.
	// It's a much cleaner way to avoid repeating the same logic for every query.
	execAndLog := func(query string) {
		_, err := db.Exec(query)

		// Check the error immediately and handle it before trying to use `result`.
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				fmt.Printf("ERREUR : <createdefault.go> Etape d'insertion ignorée : le rôle ou la catégorie existe déjà.")
			} else {
				log.Printf("<createdefault.go> Fatal error during insertion: %v\n", err)
				panic(err) // Panic on a fatal error so you can see the full stack trace.
			}
			return // Don't proceed if there was an error
		}

	}

	// Inserts
	log.Println("INIT : Création de la base de données initiale : Création des rôles")
	execAndLog(`INSERT INTO role (name) VALUES ('ADMIN')`)
	execAndLog(`INSERT INTO role (name) VALUES ('MODO')`)
	execAndLog(`INSERT INTO role (name) VALUES ('MEMBRE')`)
	execAndLog(`INSERT INTO role (name) VALUES ('BANNI')`)
	execAndLog(`INSERT INTO role (name) VALUES ('MEMBRE ayant demandé à devenir modo')`)

	log.Println("INIT : Création de la base de données initiale : Création des catégories")
	execAndLog(`INSERT INTO category (name, description) VALUES ('Pensée positive', 
	'Partagez vos astuces pour se libérer des pensées négatives, de l''anxiété et de tout ce qui pollue l''esprit.
	Ou venez juste nous partagez vos moments de joie et vos pensées positives.')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Productivité et bonnes habitudes', 
	'Besoin de booster votre productivité ? De changer une habitude qui vous pourrit la vie ? 
	Envie de donner des conseils aux autres pour profiter au mieux de son quotidien ?
	Partagez vos bonnes habitudes !')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Alimentation et exercice', 
	'Être bien dans sa tête ne suffit pas. 
	"Un esprit sain dans un corps sain" implique aussi de bien gérer son alimentation et son exercice pour éviter les excès ou les manques d''un côté ou de l''autre. 
	Venez demander des conseils ou partager vos astuces !')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Gérer son budget', 
	'L''argent n''est pas facile pour tout le monde. 
	Que vous soyez extrêmement doué pour le gérer ou, au contraire, que vous ayez besoin d''aide pour ne pas vous retrouver dépassé, c''est ici !')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Minimalisme et écologie', 
	'Ce n''est pas toujours facile de faire la part des choses. 
	Venez trouver ou donner des petites astuces pour un quotidien plus sain, pour vous comme pour la planète.')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Améliorer sa vie sociale', 
	'Communiquer n''est pas un don inné. 
	Tout le monde peut avoir besoin de conseil pour gérer une situation, se sentir mieux dans ses rapports avec son entourage ou simplement s''adapter à de nouvelles personnes.')`)
	execAndLog(`INSERT INTO category (name, description) VALUES ('Voyage, découverte et aventure', 
	'Vous avez des projets ou des souvenirs de voyage ? Des envie d''aventure ? Partagez-les avec nous !')`)
}
