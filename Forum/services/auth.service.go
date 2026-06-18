// Le package "services" contient les règles métier : les vérifications et
// les décisions (a-t-on le droit ? les données sont-elles valides ?).
// Les services se placent entre les controllers et les repositories.
package services

import (
	"errors"
	"forum/auth"
	"forum/dto"
	"forum/models"
	"forum/repositories"
	"strings"
)

// AuthService s'occupe de l'inscription et de la connexion des membres.
type AuthService struct {
	userRepo *repositories.UserRepository
}

func InitAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register crée un nouveau compte après avoir vérifié que tout est correct.
func (s *AuthService) Register(req dto.RegisterRequest) error {
	// On nettoie les entrées : on enlève les espaces inutiles et on met l'email en minuscules.
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Vérification 1 : aucun champ ne doit être vide.
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("Tous les champs sont obligatoires")
	}
	// Vérification 2 : le mot de passe et sa confirmation doivent être identiques.
	if req.Password != req.Confirm {
		return errors.New("Les mots de passe ne correspondent pas")
	}
	// Vérification 3 : le mot de passe doit être assez solide (voir auth/validation).
	if err := auth.ValidatePassword(req.Password); err != nil {
		return err
	}
	// Vérification 4 : le pseudo ne doit pas déjà exister.
	if s.userRepo.ExistsByUsername(req.Username) {
		return errors.New("Ce nom d'utilisateur est déjà pris")
	}
	// Vérification 5 : l'email ne doit pas déjà exister.
	if s.userRepo.ExistsByEmail(req.Email) {
		return errors.New("Cette adresse email est déjà utilisée")
	}

	// Tout est bon : on hache le mot de passe avant de l'enregistrer (jamais en clair !).
	hashed := auth.HashPassword(req.Password)
	_, err := s.userRepo.Create(req.Username, req.Email, hashed)
	return err
}

// Login vérifie les identifiants et renvoie l'utilisateur + un jeton de connexion.
func (s *AuthService) Login(req dto.LoginRequest) (*models.User, string, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, "", errors.New("Identifiant et mot de passe requis")
	}

	// On retrouve l'utilisateur par son pseudo OU son email.
	user, err := s.userRepo.FindByIdentifier(req.Identifier)
	if err != nil || user == nil {
		// Message volontairement vague : on ne dit pas si c'est l'identifiant ou le mot de passe
		// qui est faux, pour ne pas aider un éventuel pirate.
		return nil, "", errors.New("Identifiant ou mot de passe incorrect")
	}

	// On compare le mot de passe tapé avec celui enregistré (en les hachant tous les deux).
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, "", errors.New("Identifiant ou mot de passe incorrect")
	}

	// Un compte banni ne peut pas se connecter.
	if user.Banned {
		return nil, "", errors.New("Votre compte a été banni de la plateforme")
	}

	// Identifiants corrects : on fabrique le jeton qui servira de "carte d'identité".
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role, user.Banned)
	if err != nil {
		return nil, "", errors.New("Erreur lors de la génération du token")
	}

	return user, token, nil
}
