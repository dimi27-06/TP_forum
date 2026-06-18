package services

import (
	"errors"
	"forum/dto"
	"forum/models"
	"forum/repositories"
	"strings"
)

// ThreadService gère les règles autour des fils de discussion.
// Il a besoin de deux repositories : celui des fils et celui des tags.
type ThreadService struct {
	threadRepo *repositories.ThreadRepository
	tagRepo    *repositories.TagRepository
}

func InitThreadService(threadRepo *repositories.ThreadRepository, tagRepo *repositories.TagRepository) *ThreadService {
	return &ThreadService{threadRepo: threadRepo, tagRepo: tagRepo}
}

// List renvoie la liste des fils + les infos de pagination (combien de pages, etc.).
func (s *ThreadService) List(req dto.PaginationRequest) ([]models.Thread, dto.PaginationMeta, error) {
	threads, total, err := s.threadRepo.FindAll(req.Page, req.Limit, req.Tag, req.Search, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	// Charger les tags pour chaque fil
	// (le repository ne ramène que les fils, on complète ici avec leurs étiquettes).
	for i := range threads {
		tags, _ := s.tagRepo.FindByThreadID(threads[i].ID)
		threads[i].Tags = tags
	}

	// Calcul du nombre total de pages. La petite formule "(total + limit - 1) / limit"
	// est une astuce pour arrondir vers le haut (ex: 25 éléments par 10 = 3 pages).
	totalPages := 0
	if req.Limit > 0 {
		totalPages = (total + req.Limit - 1) / req.Limit
	}

	meta := dto.PaginationMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    req.Page > 1,                           // y a-t-il une page avant ?
		HasNext:    req.Limit > 0 && req.Page < totalPages, // y a-t-il une page après ?
	}

	return threads, meta, nil
}

// GetByID renvoie un fil pour un visiteur normal.
// Les fils archivés sont refusés (invisibles pour les membres).
func (s *ThreadService) GetByID(id int) (*models.Thread, error) {
	t, err := s.threadRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Fil de discussion introuvable")
	}
	if t.Status == models.StatusArchived {
		return nil, errors.New("Ce fil de discussion n'est plus accessible")
	}
	tags, _ := s.tagRepo.FindByThreadID(id)
	t.Tags = tags
	return t, nil
}

// GetByIDAdmin renvoie un fil SANS bloquer les archives.
// Réservé aux pages d'édition / administration.
func (s *ThreadService) GetByIDAdmin(id int) (*models.Thread, error) {
	t, err := s.threadRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Fil de discussion introuvable")
	}
	tags, _ := s.tagRepo.FindByThreadID(id)
	t.Tags = tags
	return t, nil
}

// Create vérifie les champs puis enregistre un nouveau fil, avec ses tags.
func (s *ThreadService) Create(req dto.CreateThreadRequest, userID int) (int, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" {
		return 0, errors.New("Le titre est obligatoire")
	}
	if req.Content == "" {
		return 0, errors.New("Le contenu est obligatoire")
	}

	// On crée d'abord le fil pour obtenir son id.
	id, err := s.threadRepo.Create(req.Title, req.Content, userID)
	if err != nil {
		return 0, err
	}

	// Gérer les tags
	// Pour chaque tag fourni, on le crée s'il n'existe pas, puis on l'associe au fil.
	if len(req.Tags) > 0 {
		var tagIDs []int
		for _, name := range req.Tags {
			name = strings.TrimSpace(name)
			if name == "" {
				continue // on ignore les tags vides
			}
			tagID, err := s.tagRepo.FindOrCreate(name)
			if err == nil {
				tagIDs = append(tagIDs, tagID)
			}
		}
		if err := s.tagRepo.SetThreadTags(id, tagIDs); err != nil {
			return 0, err
		}
	}

	return id, nil
}

// Update modifie un fil après avoir vérifié que la personne a le droit.
func (s *ThreadService) Update(id, userID int, role string, req dto.UpdateThreadRequest) error {
	thread, err := s.threadRepo.FindByID(id)
	if err != nil {
		return errors.New("Fil introuvable")
	}

	// Seul l'auteur ou un admin peut modifier.
	if role != "admin" && thread.UserID != userID {
		return errors.New("Non autorisé")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" {
		return errors.New("Le titre est obligatoire")
	}

	// Par défaut on garde le statut actuel. Seul un admin peut le changer.
	status := string(thread.Status)
	if role == "admin" && req.Status != "" {
		status = req.Status
	}

	if err := s.threadRepo.Update(id, req.Title, req.Content, status); err != nil {
		return err
	}

	// Tags
	// Même logique que dans Create : on recrée la liste des tags du fil.
	if len(req.Tags) > 0 {
		var tagIDs []int
		for _, name := range req.Tags {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			tagID, err := s.tagRepo.FindOrCreate(name)
			if err == nil {
				tagIDs = append(tagIDs, tagID)
			}
		}
		if err := s.tagRepo.SetThreadTags(id, tagIDs); err != nil {
			return err
		}
	}

	return nil
}

// Delete supprime un fil. Un membre normal ne peut supprimer que les siens ;
// un admin peut tout supprimer.
func (s *ThreadService) Delete(id, userID int, role string) error {
	if role != "admin" {
		thread, err := s.threadRepo.FindByID(id)
		if err != nil {
			return errors.New("Fil introuvable")
		}
		if thread.UserID != userID {
			return errors.New("Non autorisé")
		}
	}
	return s.threadRepo.Delete(id)
}

// UpdateStatus change le statut d'un fil (utilisé côté admin).
func (s *ThreadService) UpdateStatus(id int, status string) error {
	return s.threadRepo.UpdateStatus(id, status)
}

// GetAllTags renvoie tous les tags existants (pour les listes déroulantes).
func (s *ThreadService) GetAllTags() ([]models.Tag, error) {
	return s.tagRepo.FindAll()
}
