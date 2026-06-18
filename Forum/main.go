package main

import (
	"forum/app"
	"log"
	"net/http"
	"os"
)

// La fonction main est LA fonction lancée en premier au démarrage
func main() {
	// On construit toute l'application (base de données, routes, etc.).
	application := app.InitApp()

	defer application.Close()

	// On regarde si un numéro de port a été choisi dans les réglages.
	port := os.Getenv("PORT")
	if port == "" {
		// Si rien n'est précisé, on utilise le port 8080 par défaut.
		port = "8080"
	}
	log.Printf("Serveur lancé : http://localhost:%s", port)
	// On démarre vraiment le serveur. Cette ligne tourne en boucle tant que le site est en ligne.
	// Si quelque chose se passe mal (port déjà pris par exemple), on coupe tout avec un message d'erreur.
	if err := http.ListenAndServe(":"+port, application.Router); err != nil {
		log.Fatalf("Erreur lancement serveur : %s", err.Error())
	}
}
