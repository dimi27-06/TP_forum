package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"exemple_api/auth"
	"exemple_api/dto"
	"exemple_api/models"
	"exemple_api/repositories"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func InitAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (s *AuthService) Register(data dto.RegisterRequestDto) (*dto.LoginResponseDto, error) {
	nom := strings.TrimSpace(data.Nom)
	email := strings.TrimSpace(strings.ToLower(data.Email))
	password := strings.TrimSpace(data.Password)

	if nom == "" {
		return nil, errors.New("le nom est requis")
	}
	if email == "" {
		return nil, errors.New("l'email est requis")
	}
	if password == "" {
		return nil, errors.New("le mot de passe est requis")
	}

	exists, err := s.userRepository.ExistsByEmailOrName(email, nom)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ce compte existe deja")
	}

	salt, err := generateSalt(16)
	if err != nil {
		return nil, err
	}

	user := &models.ForumUser{
		Nom:              nom,
		Email:            email,
		PasswordSalt:     salt,
		PasswordHash:     hashPassword(password, salt),
		Bio:              strings.TrimSpace(data.Bio),
		Localisation:     strings.TrimSpace(data.Localisation),
		Role:             "user",
		PointsReputation: 0,
		Actif:            true,
	}

	if err := s.userRepository.CreateUser(user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(fmt.Sprintf("%d", user.ID), user.Role)
	if err != nil {
		return nil, err
	}

	_ = s.userRepository.UpdateLastLogin(user.ID)

	return &dto.LoginResponseDto{
		Type:        "Bearer",
		AccessToken: token,
		ExpiresIn:   900,
	}, nil
}

func (s *AuthService) Login(data dto.LoginRequestDto) (*dto.LoginResponseDto, error) {
	identifier := strings.TrimSpace(strings.ToLower(data.Username))
	password := strings.TrimSpace(data.Password)

	if identifier == "" || password == "" {
		return nil, errors.New("identifiants invalides")
	}

	user, err := s.userRepository.GetUserByIdentifier(identifier)
	if err != nil {
		return nil, errors.New("identifiants invalides")
	}

	expectedHash := hashPassword(password, user.PasswordSalt)
	if expectedHash != user.PasswordHash {
		return nil, errors.New("identifiants invalides")
	}

	token, err := auth.GenerateToken(fmt.Sprintf("%d", user.ID), user.Role)
	if err != nil {
		return nil, err
	}

	_ = s.userRepository.UpdateLastLogin(user.ID)

	return &dto.LoginResponseDto{
		Type:        "Bearer",
		AccessToken: token,
		ExpiresIn:   900,
	}, nil
}

func generateSalt(length int) (string, error) {
	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func hashPassword(password string, salt string) string {
	sum := sha256.Sum256([]byte(salt + ":" + password))
	return hex.EncodeToString(sum[:])
}
