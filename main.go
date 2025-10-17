package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mathis-Pain/Forum/builddb"
	"github.com/Mathis-Pain/Forum/routes"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils/external"
)

func main() {
	// Initialisation de la BDD du forum
	dbPath := "./data/forum.db"
	schemaPath := "./data/forumdbschema.sql"
	log.Printf("Initialisation des bases de données en cours.")

	db, err := builddb.InitDB(dbPath, schemaPath)
	if err != nil {
		fmt.Println("Erreur creation bdd :", err)
		return
	}

	db.Close()

	// Initialisation de la BDD des notifications
	dbPath = "./data/notifications/notifications.db"
	schemaPath = "./data/notifications/notifschema.sql"
	logs, err := builddb.InitDB(dbPath, schemaPath)
	if err != nil {
		fmt.Println("Erreur creation bdd :", err)
		return
	}
	logs.Close()

	log.Print("Projet lancé, bases de données prêtes à l'emploi")

	external.InitGoogleOAuth()
	external.InitGitHubOAuth()
	external.InitDiscordOAuth()

	// Nettoyage des sessions expirées toutes les 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			sessions.CleanupExpiredSessions()
		}
	}()

	// initialisation des routes
	mux := routes.InitRoutes()

	// démarrage serveur
	fmt.Println("Serveur démarré sur http://localhost:5080 ...")
	if err := http.ListenAndServe(":5080", mux); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}

}
