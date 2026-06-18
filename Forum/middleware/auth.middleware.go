// Le package "middleware" contient du code qui s'exécute AVANT d'arriver sur une page.
// C'est comme un videur à l'entrée : il vérifie qui passe et peut bloquer l'accès.
package middleware

import (
	"context"   // pour transporter des infos d'une étape à l'autre pendant une requête
	"forum/auth"
	"net/http"
	"strings"
)

// On crée un type spécial pour la "clé" qui range l'utilisateur dans le contexte.
// Utiliser un type à nous évite les conflits avec d'autres librairies.
type contextKey string

const UserKey contextKey = "user"

// GetClaims récupère les claims depuis le contexte (peut être nil si non connecté)
func GetClaims(r *http.Request) *auth.Claims {
	v := r.Context().Value(UserKey)
	if v == nil {
		return nil // personne n'est connecté pour cette requête
	}
	c, _ := v.(*auth.Claims)
	return c
}

// LoadUser charge les claims depuis le cookie JWT (sans bloquer)
// Ce middleware tourne sur TOUTES les pages : il regarde juste si un visiteur
// est connecté, mais ne refuse jamais l'accès lui-même.
func LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// On cherche d'abord le jeton dans un cookie du navigateur.
		cookie, err := r.Cookie("jwt_token")
		if err == nil {
			claims, err := auth.ValidateToken(cookie.Value)
			if err == nil {
				// Jeton valide : on range l'utilisateur dans le contexte de la requête.
				ctx := context.WithValue(r.Context(), UserKey, claims)
				r = r.WithContext(ctx)
			}
		}
		// Sinon chercher dans le header Authorization (API)
		if GetClaims(r) == nil {
			header := r.Header.Get("Authorization")
			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")
				claims, err := auth.ValidateToken(token)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserKey, claims)
					r = r.WithContext(ctx)
				}
			}
		}
		// On laisse la requête continuer vers la page demandée.
		next.ServeHTTP(w, r)
	})
}

// RequireAuth bloque si non connecté
// À mettre sur les pages réservées aux membres (poster un message, etc.).
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil {
			// Pas connecté : on renvoie vers la page de connexion.
			// Le "?redirect=" permet de revenir sur la page voulue après s'etre connecté.
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusFound)
			return
		}
		if claims.Banned {
			// Connecté mais banni : accès interdit.
			http.Error(w, "Compte banni", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin bloque si pas admin
// À mettre sur les pages d'administration.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// Connecté mais ce n'est pas un admin : on refuse.
		if claims.Role != "admin" {
			http.Error(w, "Accès refusé", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
