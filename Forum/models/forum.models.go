package models

import "time"

// ForumCategory représente une catégorie du forum
type ForumCategory struct {
	ID           int       `json:"id"`
	Nom          string    `json:"nom"`
	Description  string    `json:"description"`
	Slug         string    `json:"slug"`
	Icon         string    `json:"icon"`
	DateCreation time.Time `json:"date_creation"`
}

// ForumTopic représente un fil de discussion
type ForumTopic struct {
	ID               int       `json:"id"`
	Titre            string    `json:"titre"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description"`
	Contenu          string    `json:"contenu"`
	UtilisateurID    int       `json:"utilisateur_id"`
	UtilisateurNom   string    `json:"utilisateur_nom,omitempty"`
	CategorieID      int       `json:"categorie_id"`
	CategorieNom     string    `json:"categorie_nom,omitempty"`
	Vues             int       `json:"vues"`
	NombreReponses   int       `json:"nombre_reponses"`
	Epingle          bool      `json:"epingle"`
	Ferme            bool      `json:"ferme"`
	DateCreation     time.Time `json:"date_creation"`
	DateModification time.Time `json:"date_modification"`
	Likes            int       `json:"likes"`
}

// ForumComment représente une réponse à un topic
type ForumComment struct {
	ID               int       `json:"id"`
	Contenu          string    `json:"contenu"`
	UtilisateurID    int       `json:"utilisateur_id"`
	UtilisateurNom   string    `json:"utilisateur_nom,omitempty"`
	TopicID          int       `json:"topic_id"`
	DateCreation     time.Time `json:"date_creation"`
	DateModification time.Time `json:"date_modification"`
	Likes            int       `json:"likes"`
}

// ForumLike représente un like
type ForumLike struct {
	ID            int       `json:"id"`
	UtilisateurID int       `json:"utilisateur_id"`
	TopicID       *int      `json:"topic_id,omitempty"`
	CommentID     *int      `json:"comment_id,omitempty"`
	TypeLike      string    `json:"type_like"`
	DateCreation  time.Time `json:"date_creation"`
}

// CreateTopicRequest DTO pour créer un topic
type CreateTopicRequest struct {
	Titre       string `json:"titre" binding:"required"`
	Description string `json:"description"`
	Contenu     string `json:"contenu" binding:"required"`
	CategorieID int    `json:"categorie_id" binding:"required"`
}

// UpdateTopicRequest DTO pour mettre à jour un topic
type UpdateTopicRequest struct {
	Titre       string `json:"titre"`
	Description string `json:"description"`
	Contenu     string `json:"contenu"`
	Epingle     *bool  `json:"epingle,omitempty"`
	Ferme       *bool  `json:"ferme,omitempty"`
}

// CreateCommentRequest DTO pour créer un commentaire
type CreateCommentRequest struct {
	Contenu string `json:"contenu" binding:"required"`
}

// CreateCategoryRequest DTO pour créer une catégorie
type CreateCategoryRequest struct {
	Nom         string `json:"nom" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug" binding:"required"`
	Icon        string `json:"icon"`
}
