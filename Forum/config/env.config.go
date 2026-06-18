package config

import (
	"bufio"   // pour lire un fichier ligne par ligne
	"log"     // pour afficher des messages
	"os"      // pour lire/écrire les variables d'environnement
	"strings" // pour découper et nettoyer du texte
)

// LoadEnv lit le fichier ".env" et charge ses valeurs comme variables d'environnement.
// Un fichier .env contient des réglages du genre :  PORT=8080
func LoadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		// Pas grave si le fichier n'existe pas : on utilisera les réglages du système.
		log.Println("Pas de fichier .env trouvé, utilisation des variables système")
		return
	}
	defer file.Close() // on ferme le fichier quand on a fini de le lire

	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // on lit le fichier ligne par ligne
		line := strings.TrimSpace(scanner.Text()) // on enlève les espaces autour
		// On saute les lignes vides et les commentaires (qui commence par #).
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// On coupe la ligne en deux au niveau du "=" : la clé d'un côté, la valeur de l'autre.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// On ne définit la variable que si elle n'existe pas déjà dans le système.
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

// GetEnv renvoie la valeur d'un réglage, ou une valeur de secours si rien n'est défini.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetRequiredEnv renvoie un réglage OBLIGATOIRE.
// Si on ne le trouve pas, on arrête le programme car il ne peut pas tourner sans.
func GetRequiredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Variable d'environnement manquante : %s", key)
	}
	return v
}
