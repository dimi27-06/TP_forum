package services

import (
	"errors"
	"exemple_api/models"
	"exemple_api/repositories"
	"strings"
)

// ForumService gère la logique métier du forum
type ForumService struct {
	forumRepository *repositories.ForumRepository
}

// InitForumService initialise le service
func InitForumService(forumRepository *repositories.ForumRepository) *ForumService {
	return &ForumService{
		forumRepository: forumRepository,
	}
}

// CATEGORIES

// CreateCategory crée une nouvelle catégorie
func (s *ForumService) CreateCategory(req *models.CreateCategoryRequest) (*models.ForumCategory, error) {
	if strings.TrimSpace(req.Nom) == "" {
		return nil, errors.New("le nom de la catégorie est requis")
	}

	category := &models.ForumCategory{
		Nom:         req.Nom,
		Description: req.Description,
		Slug:        req.Slug,
		Icon:        req.Icon,
	}

	err := s.forumRepository.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetAllCategories récupère toutes les catégories
func (s *ForumService) GetAllCategories() ([]models.ForumCategory, error) {
	return s.forumRepository.GetAllCategories()
}

// GetCategoryByID récupère une catégorie par ID
func (s *ForumService) GetCategoryByID(id int) (*models.ForumCategory, error) {
	if id <= 0 {
		return nil, errors.New("ID de catégorie invalide")
	}
	return s.forumRepository.GetCategoryByID(id)
}

// TOPICS

// CreateTopic crée un nouveau topic
func (s *ForumService) CreateTopic(req *models.CreateTopicRequest, userID int) (*models.ForumTopic, error) {
	if strings.TrimSpace(req.Titre) == "" {
		return nil, errors.New("le titre du topic est requis")
	}
	if strings.TrimSpace(req.Contenu) == "" {
		return nil, errors.New("le contenu du topic est requis")
	}
	if req.CategorieID <= 0 {
		return nil, errors.New("ID de catégorie invalide")
	}

	// Vérifier que la catégorie existe
	_, err := s.forumRepository.GetCategoryByID(req.CategorieID)
	if err != nil {
		return nil, errors.New("catégorie non trouvée")
	}

	topic := &models.ForumTopic{
		Titre:         req.Titre,
		Description:   req.Description,
		Contenu:       req.Contenu,
		UtilisateurID: userID,
		CategorieID:   req.CategorieID,
		Epingle:       false,
		Ferme:         false,
	}

	err = s.forumRepository.CreateTopic(topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// GetAllTopics récupère tous les topics
func (s *ForumService) GetAllTopics(limit int, offset int) ([]models.ForumTopic, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.forumRepository.GetAllTopics(limit, offset)
}

// GetTopicsByCategory récupère les topics d'une catégorie
func (s *ForumService) GetTopicsByCategory(categoryID int, limit int, offset int) ([]models.ForumTopic, error) {
	if categoryID <= 0 {
		return nil, errors.New("ID de catégorie invalide")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.forumRepository.GetTopicsByCategory(categoryID, limit, offset)
}

// GetTopicByID récupère un topic par ID
func (s *ForumService) GetTopicByID(id int) (*models.ForumTopic, error) {
	if id <= 0 {
		return nil, errors.New("ID de topic invalide")
	}
	return s.forumRepository.GetTopicByID(id)
}

// UpdateTopic met à jour un topic
func (s *ForumService) UpdateTopic(id int, req *models.UpdateTopicRequest, userID int) (*models.ForumTopic, error) {
	topic, err := s.GetTopicByID(id)
	if err != nil {
		return nil, err
	}

	// Vérifier que l'utilisateur est l'auteur du topic
	if topic.UtilisateurID != userID {
		return nil, errors.New("vous n'êtes pas autorisé à modifier ce topic")
	}

	if strings.TrimSpace(req.Titre) != "" {
		topic.Titre = req.Titre
	}
	if strings.TrimSpace(req.Description) != "" {
		topic.Description = req.Description
	}
	if strings.TrimSpace(req.Contenu) != "" {
		topic.Contenu = req.Contenu
	}
	if req.Epingle != nil {
		topic.Epingle = *req.Epingle
	}
	if req.Ferme != nil {
		topic.Ferme = *req.Ferme
	}

	err = s.forumRepository.UpdateTopic(id, topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// DeleteTopic supprime un topic
func (s *ForumService) DeleteTopic(id int, userID int) error {
	topic, err := s.GetTopicByID(id)
	if err != nil {
		return err
	}

	// Vérifier que l'utilisateur est l'auteur du topic
	if topic.UtilisateurID != userID {
		return errors.New("vous n'êtes pas autorisé à supprimer ce topic")
	}

	return s.forumRepository.DeleteTopic(id)
}

// COMMENTS

// CreateComment crée un nouveau commentaire
func (s *ForumService) CreateComment(topicID int, contenu string, userID int) (*models.ForumComment, error) {
	if topicID <= 0 {
		return nil, errors.New("ID de topic invalide")
	}
	if strings.TrimSpace(contenu) == "" {
		return nil, errors.New("le contenu du commentaire est requis")
	}

	// Vérifier que le topic existe
	topic, err := s.GetTopicByID(topicID)
	if err != nil {
		return nil, errors.New("topic non trouvé")
	}

	// Vérifier que le topic n'est pas fermé
	if topic.Ferme {
		return nil, errors.New("ce topic est fermé, vous ne pouvez pas ajouter de commentaires")
	}

	comment := &models.ForumComment{
		Contenu:       contenu,
		UtilisateurID: userID,
		TopicID:       topicID,
	}

	err = s.forumRepository.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByTopic récupère les commentaires d'un topic
func (s *ForumService) GetCommentsByTopic(topicID int) ([]models.ForumComment, error) {
	if topicID <= 0 {
		return nil, errors.New("ID de topic invalide")
	}
	return s.forumRepository.GetCommentsByTopic(topicID)
}

// GetCommentByID récupère un commentaire par ID
func (s *ForumService) GetCommentByID(id int) (*models.ForumComment, error) {
	if id <= 0 {
		return nil, errors.New("ID de commentaire invalide")
	}
	return s.forumRepository.GetCommentByID(id)
}

// UpdateComment met à jour un commentaire
func (s *ForumService) UpdateComment(id int, contenu string, userID int) (*models.ForumComment, error) {
	comment, err := s.GetCommentByID(id)
	if err != nil {
		return nil, err
	}

	// Vérifier que l'utilisateur est l'auteur du commentaire
	if comment.UtilisateurID != userID {
		return nil, errors.New("vous n'êtes pas autorisé à modifier ce commentaire")
	}

	if strings.TrimSpace(contenu) == "" {
		return nil, errors.New("le contenu du commentaire est requis")
	}

	err = s.forumRepository.UpdateComment(id, contenu)
	if err != nil {
		return nil, err
	}

	comment.Contenu = contenu
	return comment, nil
}

// DeleteComment supprime un commentaire
func (s *ForumService) DeleteComment(id int, userID int) error {
	comment, err := s.GetCommentByID(id)
	if err != nil {
		return err
	}

	// Vérifier que l'utilisateur est l'auteur du commentaire
	if comment.UtilisateurID != userID {
		return errors.New("vous n'êtes pas autorisé à supprimer ce commentaire")
	}

	return s.forumRepository.DeleteComment(id)
}

// LIKES

// LikeTopic like un topic
func (s *ForumService) LikeTopic(topicID int, userID int) error {
	// Vérifier que le topic existe
	_, err := s.GetTopicByID(topicID)
	if err != nil {
		return err
	}

	// Vérifier si le like existe déjà
	existingLike, err := s.forumRepository.GetUserLike(userID, &topicID, nil)
	if err != nil {
		return err
	}

	if existingLike != nil {
		// Le like existe déjà, le supprimer
		return s.forumRepository.RemoveLike(userID, &topicID, nil)
	}

	// Créer le like
	like := &models.ForumLike{
		UtilisateurID: userID,
		TopicID:       &topicID,
		TypeLike:      "topic",
	}

	return s.forumRepository.AddLike(like)
}

// LikeComment like un commentaire
func (s *ForumService) LikeComment(commentID int, userID int) error {
	// Vérifier que le commentaire existe
	_, err := s.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	// Vérifier si le like existe déjà
	existingLike, err := s.forumRepository.GetUserLike(userID, nil, &commentID)
	if err != nil {
		return err
	}

	if existingLike != nil {
		// Le like existe déjà, le supprimer
		return s.forumRepository.RemoveLike(userID, nil, &commentID)
	}

	// Créer le like
	like := &models.ForumLike{
		UtilisateurID: userID,
		CommentID:     &commentID,
		TypeLike:      "comment",
	}

	return s.forumRepository.AddLike(like)
}

// SEARCH

// Search recherche des topics
func (s *ForumService) Search(query string, limit int, offset int) ([]models.ForumTopic, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("la requête de recherche ne peut pas être vide")
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.forumRepository.Search(query, limit, offset)
}

// GetPopularTopics récupère les topics populaires
func (s *ForumService) GetPopularTopics(limit int) ([]models.ForumTopic, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.forumRepository.GetPopularTopics(limit)
}
