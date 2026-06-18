package controllers

import (
	"encoding/json" // pour renvoyer la réponse au format JSON (lisible par le JavaScript)
	"forum/middleware"
	"forum/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ReactionController gère les "j'aime" / "j'aime pas" sur les messages.
type ReactionController struct {
	reactionService *services.ReactionService
}

func InitReactionController(rs *services.ReactionService) *ReactionController {
	return &ReactionController{reactionService: rs}
}

// React - endpoint JSON pour les likes/dislikes (appelé en JS fetch)
// Contrairement aux autres pages, ici on ne renvoie pas du HTML mais des données JSON,
// car c'est le JavaScript de la page qui met à jour le compteur sans recharger.
func (c *ReactionController) React(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		// Il faut être connecté pour réagir.
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	msgID, _ := strconv.Atoi(mux.Vars(r)["id"])
	reactionType := mux.Vars(r)["type"] // "like" ou "dislike"

	// Le service enregistre la réaction et renvoie les nouveaux compteurs.
	resp, err := c.reactionService.React(claims.UserID, msgID, reactionType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// On indique que la réponse est du JSON, puis on l'envoie.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
