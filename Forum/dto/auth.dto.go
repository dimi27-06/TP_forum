// Le package "dto" (Data Transfer Object) regroupe des structures simples
// qui servent à transporter les données saisies dans les formulaires.
package dto

// RegisterRequest = ce qu'on récupère du formulaire d'inscription.
type RegisterRequest struct {
	Username string // le pseudo choisi
	Email    string // l'adresse email
	Password string // le mot de passe
	Confirm  string // la confirmation du mot de passe (doit être identique)
}

// LoginRequest = ce qu'on récupère du formulaire de connexion.
type LoginRequest struct {
	Identifier string // username ou email
	Password   string
}
