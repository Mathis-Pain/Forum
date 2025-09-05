package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/routes"
	"github.com/Mathis-Pain/Forum/utils"
)

func main() {
	//initialisation de la bdd
	db := utils.InitDB()
	defer db.Close()
	fmt.Println("Projet lancé, DB prête à l'emploi")

	// initialisation des routes
	mux := routes.InitRoutes()

	// démarrage serveur
	fmt.Println("Serveur démarré sur http://localhost:5080 ...")
	if err := http.ListenAndServe(":5080", mux); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}
