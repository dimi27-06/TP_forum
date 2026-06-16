package controllers

import (
	"encoding/json"
	"exemple_api/helper"
	"exemple_api/models"
	"exemple_api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ForumControllers gère les requêtes du forum
type ForumControllers struct {
	forumService *services.ForumService
}

// InitForumControllers initialise le controller
func InitForumControllers(forumService *services.ForumService) *ForumControllers {
	return &ForumControllers{
		forumService: forumService,
	}
}

// CATEGORIES

// CreateCategory crée une nouvelle catégorie
func (c *ForumControllers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Données invalides")
		return
	}

	category, err := c.forumService.CreateCategory(&req)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusCreated, category)
}

// GetAllCategories récupère toutes les catégories
func (c *ForumControllers) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.forumService.GetAllCategories()
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des catégories")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, categories)
}

// GetCategoryByID récupère une catégorie par ID
func (c *ForumControllers) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID invalide")
		return
	}

	category, err := c.forumService.GetCategoryByID(id)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusNotFound, "Catégorie non trouvée")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, category)
}

// TOPICS

// CreateTopic crée un nouveau topic
func (c *ForumControllers) CreateTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	var req models.CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Données invalides")
		return
	}

	topic, err := c.forumService.CreateTopic(&req, uid)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusCreated, topic)
}

// GetAllTopics récupère tous les topics
func (c *ForumControllers) GetAllTopics(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	topics, err := c.forumService.GetAllTopics(limit, offset)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des topics")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topics)
}

// GetTopicsByCategory récupère les topics d'une catégorie
func (c *ForumControllers) GetTopicsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de catégorie invalide")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	topics, err := c.forumService.GetTopicsByCategory(categoryID, limit, offset)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des topics")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topics)
}

// GetTopicByID récupère un topic par ID
func (c *ForumControllers) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	topic, err := c.forumService.GetTopicByID(topicID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusNotFound, "Topic non trouvé")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topic)
}

// UpdateTopic met à jour un topic
func (c *ForumControllers) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	var req models.UpdateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Données invalides")
		return
	}

	topic, err := c.forumService.UpdateTopic(topicID, &req, uid)
	if err != nil {
		if err.Error() == "vous n'êtes pas autorisé à modifier ce topic" {
			helper.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else {
			helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topic)
}

// DeleteTopic supprime un topic
func (c *ForumControllers) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	err = c.forumService.DeleteTopic(topicID, uid)
	if err != nil {
		if err.Error() == "vous n'êtes pas autorisé à supprimer ce topic" {
			helper.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else {
			helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// COMMENTS

// CreateComment crée un nouveau commentaire
func (c *ForumControllers) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	var req struct {
		Contenu string `json:"contenu"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Données invalides")
		return
	}

	comment, err := c.forumService.CreateComment(topicID, req.Contenu, uid)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusCreated, comment)
}

// GetCommentsByTopic récupère les commentaires d'un topic
func (c *ForumControllers) GetCommentsByTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	comments, err := c.forumService.GetCommentsByTopic(topicID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des commentaires")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, comments)
}

// UpdateComment met à jour un commentaire
func (c *ForumControllers) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	commentID, err := strconv.Atoi(mux.Vars(r)["commentId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de commentaire invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	var req struct {
		Contenu string `json:"contenu"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Données invalides")
		return
	}

	comment, err := c.forumService.UpdateComment(commentID, req.Contenu, uid)
	if err != nil {
		if err.Error() == "vous n'êtes pas autorisé à modifier ce commentaire" {
			helper.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else {
			helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, comment)
}

// DeleteComment supprime un commentaire
func (c *ForumControllers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	commentID, err := strconv.Atoi(mux.Vars(r)["commentId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de commentaire invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	err = c.forumService.DeleteComment(commentID, uid)
	if err != nil {
		if err.Error() == "vous n'êtes pas autorisé à supprimer ce commentaire" {
			helper.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else {
			helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// LIKES

// LikeTopic like un topic
func (c *ForumControllers) LikeTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topicId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de topic invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	err = c.forumService.LikeTopic(topicID, uid)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, map[string]string{"message": "Like ajouté/supprimé avec succès"})
}

// LikeComment like un commentaire
func (c *ForumControllers) LikeComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		helper.WriteErrorResponse(w, http.StatusUnauthorized, "Authentification requise")
		return
	}

	commentID, err := strconv.Atoi(mux.Vars(r)["commentId"])
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID de commentaire invalide")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "ID utilisateur invalide")
		return
	}

	err = c.forumService.LikeComment(commentID, uid)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, map[string]string{"message": "Like ajouté/supprimé avec succès"})
}

// SEARCH

// Search recherche des topics
func (c *ForumControllers) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		helper.WriteErrorResponse(w, http.StatusBadRequest, "Paramètre de recherche manquant")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	topics, err := c.forumService.Search(query, limit, offset)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topics)
}

// GetPopularTopics récupère les topics populaires
func (c *ForumControllers) GetPopularTopics(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	topics, err := c.forumService.GetPopularTopics(limit)
	if err != nil {
		helper.WriteErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la récupération des topics")
		return
	}

	helper.WriteSuccessResponse(w, http.StatusOK, topics)
}
