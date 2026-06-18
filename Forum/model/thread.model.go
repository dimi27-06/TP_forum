package models

import "time"

// ThreadStatus est un type spécial qui ne peut prendre que certaines valeurs de texte.
// Ça évite les fautes de frappe : on est sûr d'utiliser un statut valide.
type ThreadStatus string

// Les 3 statuts possibles pour un fil de discussion.
const (
	StatusOpen     ThreadStatus = "open"     // ouvert : tout le monde peut répondre
	StatusClosed   ThreadStatus = "closed"   // fermé : on ne peut plus répondre
	StatusArchived ThreadStatus = "archived" // archivé : caché aux membres normaux
)

// Thread représente un fil de discussion (un sujet).
type Thread struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Status    ThreadStatus `json:"status"`
	UserID    int          `json:"user_id"` // id de l'auteur
	Author    string       `json:"author"`  // pseudo de l'auteur (pratique pour l'affichage)
	Tags      []Tag        `json:"tags"`    // les étiquettes du sujet
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	MsgCount  int          `json:"msg_count"` // nombre de messages dans le sujet
}
