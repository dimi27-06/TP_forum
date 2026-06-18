package auth

import (
	"errors"  // pour renvoyer un message quand le mot de passe ne convient pas
	"unicode" // pour examiner chaque caractère (majuscule ? chiffre ? etc.)
)

// ValidatePassword vérifie qu'un mot de passe est assez solide.
// Si tout va bien elle renvoie "nil" (aucune erreur).
// Sinon elle renvoie un message expliquant ce qui ne va pas.
func ValidatePassword(password string) error {
	// Règle 1 : au moins 12 caractères.
	if len(password) < 12 {
		return errors.New("Le mot de passe doit contenir au moins 12 caractères")
	}
	// On va parcourir le mot de passe lettre par lettre pour vérifier 2 choses.
	var hasUpper, hasSpecial bool
	for _, c := range password {
		// Y a-t-il au moins une majuscule ?
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		// Un caractère spécial = ni une lettre, ni un chiffre (ex: ! ou @).
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			hasSpecial = true
		}
	}
	// Règle 2 : il faut une majuscule.
	if !hasUpper {
		return errors.New("Le mot de passe doit contenir au moins une majuscule")
	}
	// Règle 3 : il faut un caractère spécial.
	if !hasSpecial {
		return errors.New("Le mot de passe doit contenir au moins un caractère spécial")
	}
	// Tout est bon, le mot de passe est accepté.
	return nil
}
