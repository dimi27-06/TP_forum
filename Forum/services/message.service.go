package services

import (
	"errors"
	"forum/dto"
	"forum/models"
	"forum/repositories"
	"strings"
)

// MessageService gère les règles autour des messages.
type MessageService struct {
	messageRepo *repositories.MessageRepository
	threadRepo  *repositories.ThreadRepository
}

func InitMessageService(messageRepo *repositories.MessageRepository, threadRepo *repositories.ThreadRepository) *MessageService {
	return &MessageService{messageRepo: messageRepo, threadRepo: threadRepo}
}

// Create poste un message, mais seulement si le fil est encore ouvert.
func (s *MessageService) Create(req dto.CreateMessageRequest, userID int) (int, error) {
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return 0, errors.New("Le message ne peut pas être vide")
	}

	// On vérifie que le fil existe.
	thread, err := s.threadRepo.FindByID(req.ThreadID)
	if err != nil {
		return 0, errors.New("Fil de discussion introuvable")
	}
	// On n'autorise les réponses que si le fil est "ouvert" (ni fermé, ni archivé).
	if thread.Status != models.StatusOpen {
		return 0, errors.New("Ce fil de discussion n'accepte plus de nouveaux messages")
	}

	return s.messageRepo.Create(req.Content, req.ThreadID, userID)
}

// GetByThread renvoie les messages d'un fil + les infos de pagination.
func (s *MessageService) GetByThread(threadID, page, limit int, sort string, currentUserID int) ([]models.Message, dto.PaginationMeta, error) {
	messages, total, err := s.messageRepo.FindByThreadID(threadID, page, limit, sort, currentUserID)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	// Même calcul d'arrondi vers le haut que pour les fils.
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	meta := dto.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    limit > 0 && page < totalPages,
	}

	return messages, meta, nil
}

// GetByID renvoie un seul message d'après son identifiant.
func (s *MessageService) GetByID(id int) (*models.Message, error) {
	return s.messageRepo.FindByID(id)
}

// Update modifie un message si la personne en a le droit (auteur ou admin).
func (s *MessageService) Update(id, userID int, role string, req dto.UpdateMessageRequest) error {
	msg, err := s.messageRepo.FindByID(id)
	if err != nil || msg == nil {
		return errors.New("Message introuvable")
	}
	if role != "admin" && msg.UserID != userID {
		return errors.New("Non autorisé")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return errors.New("Le message ne peut pas être vide")
	}

	return s.messageRepo.Update(id, content)
}

// Delete supprime un message. Un membre normal ne peut effacer que les siens.
func (s *MessageService) Delete(id, userID int, role string) error {
	if role != "admin" {
		msg, err := s.messageRepo.FindByID(id)
		if err != nil || msg == nil {
			return errors.New("Message introuvable")
		}
		if msg.UserID != userID {
			return errors.New("Non autorisé")
		}
	}
	return s.messageRepo.Delete(id)
}
