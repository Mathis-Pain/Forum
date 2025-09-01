package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mathis-Pain/Forum/routes"
)

func main() {
	mux := routes.InitRoutes()
	fmt.Println("Serveur démarré sur http://localhost:5080 ...")
	if err := http.ListenAndServe(":5080", mux); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}
