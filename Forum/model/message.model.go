package models

import "time"

// Message représente une réponse postée dans un fil de discussion.
type Message struct {
	ID            int       `json:"id"`
	Content       string    `json:"content"`
	ThreadID      int       `json:"thread_id"` // dans quel sujet
	UserID        int       `json:"user_id"`   // id de l'auteur
	Author        string    `json:"author"`    // pseudo de l'auteur
	Likes         int       `json:"likes"`     // nombre de "j'aime"
	Dislikes      int       `json:"dislikes"`  // nombre de "j'aime pas"
	Score         int       `json:"score"`     // likes - dislikes
	UserReaction  string    `json:"user_reaction"` // "like", "dislike", ou ""
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
