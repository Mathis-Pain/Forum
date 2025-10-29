package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServer démarre un serveur HTTP sécurisé avec gestion des timeouts et arrêt gracieux
func StartServer(addr string, handler http.Handler) {
	const maxRequestSize = 10 * 1024 * 1024                     // 10 MB
	limitedMux := http.MaxBytesHandler(handler, maxRequestSize) // limite requêtes à 10 MB

	srv := &http.Server{
		Addr:              addr,
		Handler:           limitedMux,
		ReadTimeout:       15 * time.Second, // lire la requete
		WriteTimeout:      15 * time.Second, // envoyer la reponse au client
		IdleTimeout:       60 * time.Second, // temp d'inactivité entre deux requete
		ReadHeaderTimeout: 10 * time.Second, // lire uniquement les entetes http
		MaxHeaderBytes:    1 * 1024 * 1024,  // 1MB maximum pour l'entete
	}
	// Permet d’arrêter le serveur de manière plus propre, en terminant les requêtes en cours avant l’interruption.
	stop := make(chan os.Signal, 1)                    //Crée un canal pour recevoir un signal d’arrêt
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM) //Informe Go de relayer ces signaux dans le canal: os.Interrupt recoi le terminal et syscall.SIGTERM recois le docker

	go func() {
		log.Printf(" Serveur démarré sur http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" Erreur serveur : %v", err)
		}
	}()

	<-stop // 	Attend le signal avant de lancer la procédure d’arrêt
	log.Println(" Signal d'arrêt reçu, fermeture gracieuse du serveur...")

	// context.Background() → point de départ d’un nouveau contexte “vide”. context.WithTimeout(...) → crée une copie de ce contexte
	// maisd qui s’annule automatiquement au bout de 15 secondes,
	// ou plus tôt si on appelle manuellement cancel().
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erreur fermeture serveur : %v", err)
	}
	log.Println("Serveur arrêté correctement")
}
