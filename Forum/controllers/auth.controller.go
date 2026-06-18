package controllers

import (
	"encoding/json"
	"exemple_api/auth"
	"exemple_api/dto"
	"exemple_api/helper"
	"exemple_api/services"
	"net/http"
)

type AuthControllers struct {
	service *services.AuthService
}

// Controle l'authentification
func AuthProductController(authService *services.AuthService) *AuthControllers {
	return &AuthControllers{service: authService}
}

// verifie les identifiants de connexion et renvoie un token si ils sont valides.
func (c *AuthControllers) Login(w http.ResponseWriter, r *http.Request) {
	var data dto.LoginRequestDto

	// On decode le JSON envoyé dans le body de la requete.
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	// Le service verifie les identifiants et genere la reponse de connexion.
	response, err := c.service.Login(data)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiants invalides")
		return
	}

	// Si tout est valide renvoie le token au client.
	helper.WriteJSON(w, http.StatusOK, response)
}

// Crée un utilisateur et connecte l'utilisateur automatiquement après l'inscription.
func (c *AuthControllers) Register(w http.ResponseWriter, r *http.Request) {
	var data dto.RegisterRequestDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	response, err := c.service.Register(data)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, response)
}

// Retourne les informations de l'utilisateur authentifie.
func (c *AuthControllers) Me(w http.ResponseWriter, r *http.Request) {
	// Le middleware d'authentification ajoute les claims JWT dans le contexte.
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		helper.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]any{
		"code":     http.StatusOK,
		"message":  "authenticated",
		"user_id":  claims.UserID,
		"role":     claims.Role,
		"subject":  claims.Subject,
		"issuer":   claims.Issuer,
		"audience": claims.Audience,
	})
}
