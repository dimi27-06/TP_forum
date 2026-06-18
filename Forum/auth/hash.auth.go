// Le package "auth" regroupe tout ce qui touche à la sécurité :
// les mots de passe et les jetons de connexion.
package auth

import (
	"crypto/sha512" // un algorithme qui transforme un texte en empreinte impossible à inverser
	"fmt"           // pour mettre en forme du texte
)

// HashPassword transforme un mot de passe en une longue suite de caractères illisibles.
// On ne stocke JAMAIS le vrai mot de passe en base, seulement cette empreinte.
// Comme ça, même si quelqu'un vole la base, il ne voit pas les mots de passe.
func HashPassword(password string) string {
	h := sha512.New()         // on prépare la machine à fabriquer l'empreinte
	h.Write([]byte(password)) // on lui donne le mot de passe à digérer
	// On renvoie l'empreinte sous forme de texte (en hexadécimal, d'où le %x).
	return fmt.Sprintf("%x", h.Sum(nil))
}

// CheckPassword vérifie si un mot de passe tapé correspond à celui enregistré.
// Astuce : on ré-hashe le mot de passe tapé et on compare les deux empreintes.
func CheckPassword(plain, hashed string) bool {
	return HashPassword(plain) == hashed
}
