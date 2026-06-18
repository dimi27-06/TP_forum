package auth

import (
	"errors" // pour fabriquer nos propres messages d'erreur
	"os"     // pour lire la clé secrète dans les réglages
	"time"   // pour gérer les dates (création et expiration du jeton)

	"github.com/golang-jwt/jwt/v5" // librairie qui sait créer et lire les jetons JWT
)

// Un JWT (JSON Web Token) est une sorte de carte d'identité numérique.
// Quand on se connecte, le serveur nous donne ce jeton ; on le renvoie ensuite
// à chaque page pour prouver qui on est, sans avoir à retaper son mot de passe.

// Claims = les informations que l'on glisse à l'intérieur du jeton.
type Claims struct {
	UserID               int    `json:"user_id"`  // l'identifiant de l'utilisateur
	Username             string `json:"username"` // son pseudo
	Role                 string `json:"role"`     // son rôle : "user" ou "admin"
	Banned               bool   `json:"banned"`   // est-il banni ou pas
	jwt.RegisteredClaims        // champs standards fournis par la librairie (date, etc.)
}

// getSecret renvoie la clé secrète qui sert à signer les jetons.
// Cette signature empêche quelqu'un de fabriquer un faux jeton.
func getSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		// Valeur de secours si on a oublié de la définir (à changer en vrai !).
		s = "fallback_secret_change_me"
	}
	return []byte(s)
}

// GenerateToken fabrique un nouveau jeton pour un utilisateur qui vient de se connecter.
func GenerateToken(userID int, username, role string, banned bool) (string, error) {
	now := time.Now()
	// On remplit la carte d'identité avec les infos de l'utilisateur.
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Banned:   banned,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    "forum",
			IssuedAt:  jwt.NewNumericDate(now),                     // date de création
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // expire au bout de 24h
		},
	}
	// On crée le jeton et on le signe avec notre clé secrète.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// ValidateToken vérifie qu'un jeton reçu est bien authentique et pas expiré,
// puis nous redonne les infos qu'il contient.
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// On s'assure que le jeton utilise bien la méthode de signature attendue.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("méthode de signature invalide")
		}
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	// On récupère les infos et on vérifie que le jeton est valide.
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide")
	}
	return claims, nil
}
