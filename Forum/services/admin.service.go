package services

import (
	"forum/models"
	"forum/repositories"
)

// AdminService regroupe les actions d'administration sur les utilisateurs.
type AdminService struct {
	userRepo   *repositories.UserRepository
	threadRepo *repositories.ThreadRepository
}

func InitAdminService(userRepo *repositories.UserRepository, threadRepo *repositories.ThreadRepository) *AdminService {
	return &AdminService{userRepo: userRepo, threadRepo: threadRepo}
}

// GetAllUsers renvoie la liste de tous les membres (pour le tableau de bord admin).
func (s *AdminService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// BanUser marque un utilisateur comme banni.
func (s *AdminService) BanUser(userID int) error {
	return s.userRepo.SetBanned(userID, true)
}

// UnbanUser retire le bannissement d'un utilisateur.
func (s *AdminService) UnbanUser(userID int) error {
	return s.userRepo.SetBanned(userID, false)
}
