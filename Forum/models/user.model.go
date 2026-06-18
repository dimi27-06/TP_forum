// Le package "models" décrit la forme des données : à quoi ressemble un utilisateur,
// un message, un fil de discussion, etc. C'est la "carte" de chaque chose en base.
package models

import "time" // pour le type date/heure

// User représente un membre du forum.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	// Le `json:"-"` cache le mot de passe : il ne sera jamais envoyé au navigateur.
	Password  string    `json:"-"`
	Role      string    `json:"role"`   // "user" ou "admin"
	Banned    bool      `json:"banned"` // banni ou non
	CreatedAt time.Time `json:"created_at"`
}
