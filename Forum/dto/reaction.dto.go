package dto

// ReactionRequest = ce qu'on envoie quand on clique sur "j'aime" / "j'aime pas".
type ReactionRequest struct {
	MessageID int
	Type      string // "like" ou "dislike"
}

// ReactionResponse = ce que le serveur renvoie après une réaction :
// le nouveau total de likes/dislikes et la réaction de l'utilisateur courant.
// Les `json:"..."` indiquent le nom à utiliser quand on envoie ça au navigateur.
type ReactionResponse struct {
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Score    int    `json:"score"` // likes moins dislikes
	UserReaction string `json:"user_reaction"`
}
