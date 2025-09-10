package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mathis-Pain/Forum/routes"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
)

func main() {
	//initialisation de la bdd
	db := utils.InitDB()
	defer db.Close()
	fmt.Println("Projet lancé, DB prête à l'emploi")
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
