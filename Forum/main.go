// Le "package main" c'est le point de départ obligatoire de tout programme Go.
// C'est ici que tout commence quand on lance l'application.
package main

import (
	"forum/app" // notre propre code (le dossier app qui prépare le serveur)
	"log"       // sert à afficher des messages dans la console (logs)
	"net/http"  // la boite à outils de Go pour faire un serveur web
	"os"        // sert à lire des infos du système, ici les variables d'environnement
)

// La fonction main est LA fonction lancée en premier au démarrage.
// Sans elle le programme ne sait pas quoi faire.
func main() {
	// On construit toute l'application (base de données, routes, etc.).
	application := app.InitApp()
	// "defer" veut dire : fais ça à la toute fin, juste avant de quitter.
	// Ici on ferme proprement la base de données pour ne rien laisser ouvert.
	defer application.Close()

	// On regarde si un numéro de port a été choisi dans les réglages.
	port := os.Getenv("PORT")
	if port == "" {
		// Si rien n'est précisé, on utilise le port 8080 par défaut.
		port = "8080"
	}

	// On affiche l'adresse à ouvrir dans le navigateur pour voir le site.
	log.Printf("Serveur lancé : http://localhost:%s", port)
	// On démarre vraiment le serveur. Cette ligne tourne en boucle tant que le site est en ligne.
	// Si quelque chose se passe mal (port déjà pris par exemple), on coupe tout avec un message d'erreur.
	if err := http.ListenAndServe(":"+port, application.Router); err != nil {
		log.Fatalf("Erreur lancement serveur : %s", err.Error())
	}
}
