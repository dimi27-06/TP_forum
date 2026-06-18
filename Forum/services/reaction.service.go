package services

import (
	"errors"
	"forum/dto"
	"forum/repositories"
)

// ReactionService gère la logique des "j'aime" / "j'aime pas".
type ReactionService struct {
	reactionRepo *repositories.ReactionRepository
}

func InitReactionService(reactionRepo *repositories.ReactionRepository) *ReactionService {
	return &ReactionService{reactionRepo: reactionRepo}
}

// React enregistre (ou retire) la réaction d'un membre sur un message,
// puis renvoie les compteurs à jour.
func (s *ReactionService) React(userID, messageID int, reactionType string) (*dto.ReactionResponse, error) {
	// On n'accepte que "like" ou "dislike".
	if reactionType != "like" && reactionType != "dislike" {
		return nil, errors.New("Type de réaction invalide")
	}

	// On regarde si le membre avait déjà réagi à ce message.
	existing, err := s.reactionRepo.FindByUserAndMessage(userID, messageID)
	if err != nil {
		return nil, err
	}

	if existing != nil && string(existing.Type) == reactionType {
		// Même réaction → on la retire (toggle)
		// (recliquer sur "like" alors qu'on avait déjà liké = annuler le like)
		s.reactionRepo.Delete(userID, messageID)
	} else {
		// Nouvelle réaction ou changement
		// (première réaction, ou on passe de like à dislike)
		s.reactionRepo.Upsert(userID, messageID, reactionType)
	}

	// On récupère les nouveaux totaux après la modification.
	likes, dislikes, score := s.reactionRepo.GetScore(messageID)

	// On regarde aussi quelle est, maintenant, la réaction de ce membre (ou aucune).
	userReaction := ""
	updated, _ := s.reactionRepo.FindByUserAndMessage(userID, messageID)
	if updated != nil {
		userReaction = string(updated.Type)
	}

	// On renvoie le tout au controller, qui le transformera en JSON.
	return &dto.ReactionResponse{
		Likes:        likes,
		Dislikes:     dislikes,
		Score:        score,
		UserReaction: userReaction,
	}, nil
}
