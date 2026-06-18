package dto

// CreateMessageRequest = les infos pour poster un nouveau message.
type CreateMessageRequest struct {
	Content  string // le texte du message
	ThreadID int    // dans quel fil de discussion il est posté
}

// UpdateMessageRequest = les infos pour modifier un message existant.
type UpdateMessageRequest struct {
	Content string
}
