package models

import "time"

// ReactionType ne peut être que "like" ou "dislike" (même principe que ThreadStatus).
type ReactionType string

const (
	ReactionLike    ReactionType = "like"
	ReactionDislike ReactionType = "dislike"
)

// Reaction représente le clic d'un membre sur un message ("j'aime" ou pas).
type Reaction struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`    // qui a réagi
	MessageID int          `json:"message_id"` // sur quel message
	Type      ReactionType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
}
