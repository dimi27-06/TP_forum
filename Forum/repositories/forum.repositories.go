package repositories

import (
	"database/sql"
	"exemple_api/models"
	"fmt"
	"strings"
)

// ForumRepository gère l'accès aux données du forum
type ForumRepository struct {
	Db *sql.DB
}

// InitForumRepository initialise le repository
func InitForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{
		Db: db,
	}
}

// CATEGORIES

// CreateCategory crée une nouvelle catégorie
func (r *ForumRepository) CreateCategory(category *models.ForumCategory) error {
	query := "INSERT INTO forum_categories (nom, description, slug, icon) VALUES (?, ?, ?, ?)"
	result, err := r.Db.Exec(query, category.Nom, category.Description, category.Slug, category.Icon)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	category.ID = int(id)
	return err
}

// GetAllCategories récupère toutes les catégories
func (r *ForumRepository) GetAllCategories() ([]models.ForumCategory, error) {
	query := "SELECT id, nom, description, slug, icon, date_creation FROM forum_categories ORDER BY nom"
	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.ForumCategory
	for rows.Next() {
		var cat models.ForumCategory
		err := rows.Scan(&cat.ID, &cat.Nom, &cat.Description, &cat.Slug, &cat.Icon, &cat.DateCreation)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetCategoryByID récupère une catégorie par ID
func (r *ForumRepository) GetCategoryByID(id int) (*models.ForumCategory, error) {
	query := "SELECT id, nom, description, slug, icon, date_creation FROM forum_categories WHERE id = ?"
	var cat models.ForumCategory
	err := r.Db.QueryRow(query, id).Scan(&cat.ID, &cat.Nom, &cat.Description, &cat.Slug, &cat.Icon, &cat.DateCreation)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// TOPICS

// CreateTopic crée un nouveau topic
func (r *ForumRepository) CreateTopic(topic *models.ForumTopic) error {
	slug := strings.ToLower(strings.ReplaceAll(topic.Titre, " ", "-"))
	query := "INSERT INTO forum_topics (titre, slug, description, contenu, utilisateur_id, categorie_id) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.Db.Exec(query, topic.Titre, slug, topic.Description, topic.Contenu, topic.UtilisateurID, topic.CategorieID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	topic.ID = int(id)
	return err
}

// GetAllTopics récupère tous les topics avec pagination
func (r *ForumRepository) GetAllTopics(limit int, offset int) ([]models.ForumTopic, error) {
	query := `SELECT t.id, t.titre, t.slug, t.description, t.contenu, t.utilisateur_id, u.nom, 
			t.categorie_id, c.nom, t.vues, t.nombre_reponses, t.epingle, t.ferme, t.date_creation, t.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_topics t
		LEFT JOIN utilisateurs u ON t.utilisateur_id = u.id
		LEFT JOIN forum_categories c ON t.categorie_id = c.id
		LEFT JOIN forum_likes l ON t.id = l.topic_id
		GROUP BY t.id
		ORDER BY t.epingle DESC, t.date_creation DESC
		LIMIT ? OFFSET ?`

	rows, err := r.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.ForumTopic
	for rows.Next() {
		var topic models.ForumTopic
		err := rows.Scan(&topic.ID, &topic.Titre, &topic.Slug, &topic.Description, &topic.Contenu,
			&topic.UtilisateurID, &topic.UtilisateurNom, &topic.CategorieID, &topic.CategorieNom,
			&topic.Vues, &topic.NombreReponses, &topic.Epingle, &topic.Ferme, &topic.DateCreation, &topic.DateModification, &topic.Likes)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

// GetTopicsByCategory récupère tous les topics d'une catégorie
func (r *ForumRepository) GetTopicsByCategory(categoryID int, limit int, offset int) ([]models.ForumTopic, error) {
	query := `SELECT t.id, t.titre, t.slug, t.description, t.contenu, t.utilisateur_id, u.nom, 
			t.categorie_id, c.nom, t.vues, t.nombre_reponses, t.epingle, t.ferme, t.date_creation, t.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_topics t
		LEFT JOIN utilisateurs u ON t.utilisateur_id = u.id
		LEFT JOIN forum_categories c ON t.categorie_id = c.id
		LEFT JOIN forum_likes l ON t.id = l.topic_id
		WHERE t.categorie_id = ?
		GROUP BY t.id
		ORDER BY t.epingle DESC, t.date_creation DESC
		LIMIT ? OFFSET ?`

	rows, err := r.Db.Query(query, categoryID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.ForumTopic
	for rows.Next() {
		var topic models.ForumTopic
		err := rows.Scan(&topic.ID, &topic.Titre, &topic.Slug, &topic.Description, &topic.Contenu,
			&topic.UtilisateurID, &topic.UtilisateurNom, &topic.CategorieID, &topic.CategorieNom,
			&topic.Vues, &topic.NombreReponses, &topic.Epingle, &topic.Ferme, &topic.DateCreation, &topic.DateModification, &topic.Likes)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

// GetTopicByID récupère un topic par ID et incrémente les vues
func (r *ForumRepository) GetTopicByID(id int) (*models.ForumTopic, error) {
	query := `SELECT t.id, t.titre, t.slug, t.description, t.contenu, t.utilisateur_id, u.nom, 
			t.categorie_id, c.nom, t.vues, t.nombre_reponses, t.epingle, t.ferme, t.date_creation, t.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_topics t
		LEFT JOIN utilisateurs u ON t.utilisateur_id = u.id
		LEFT JOIN forum_categories c ON t.categorie_id = c.id
		LEFT JOIN forum_likes l ON t.id = l.topic_id
		WHERE t.id = ?
		GROUP BY t.id`

	var topic models.ForumTopic
	err := r.Db.QueryRow(query, id).Scan(&topic.ID, &topic.Titre, &topic.Slug, &topic.Description, &topic.Contenu,
		&topic.UtilisateurID, &topic.UtilisateurNom, &topic.CategorieID, &topic.CategorieNom,
		&topic.Vues, &topic.NombreReponses, &topic.Epingle, &topic.Ferme, &topic.DateCreation, &topic.DateModification, &topic.Likes)
	if err != nil {
		return nil, err
	}

	// Incrémenter les vues
	r.Db.Exec("UPDATE forum_topics SET vues = vues + 1 WHERE id = ?", id)

	return &topic, nil
}

// UpdateTopic met à jour un topic
func (r *ForumRepository) UpdateTopic(id int, topic *models.ForumTopic) error {
	query := "UPDATE forum_topics SET titre = ?, description = ?, contenu = ?, epingle = ?, ferme = ? WHERE id = ?"
	_, err := r.Db.Exec(query, topic.Titre, topic.Description, topic.Contenu, topic.Epingle, topic.Ferme, id)
	return err
}

// DeleteTopic supprime un topic
func (r *ForumRepository) DeleteTopic(id int) error {
	_, err := r.Db.Exec("DELETE FROM forum_topics WHERE id = ?", id)
	return err
}

// COMMENTS

// CreateComment crée un nouveau commentaire
func (r *ForumRepository) CreateComment(comment *models.ForumComment) error {
	query := "INSERT INTO forum_comments (contenu, utilisateur_id, topic_id) VALUES (?, ?, ?)"
	result, err := r.Db.Exec(query, comment.Contenu, comment.UtilisateurID, comment.TopicID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	comment.ID = int(id)

	// Incrémenter le nombre de réponses du topic
	r.Db.Exec("UPDATE forum_topics SET nombre_reponses = nombre_reponses + 1 WHERE id = ?", comment.TopicID)

	return err
}

// GetCommentsByTopic récupère tous les commentaires d'un topic
func (r *ForumRepository) GetCommentsByTopic(topicID int) ([]models.ForumComment, error) {
	query := `SELECT c.id, c.contenu, c.utilisateur_id, u.nom, c.topic_id, c.date_creation, c.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_comments c
		LEFT JOIN utilisateurs u ON c.utilisateur_id = u.id
		LEFT JOIN forum_likes l ON c.id = l.comment_id
		WHERE c.topic_id = ?
		GROUP BY c.id
		ORDER BY c.date_creation ASC`

	rows, err := r.Db.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.ForumComment
	for rows.Next() {
		var comment models.ForumComment
		err := rows.Scan(&comment.ID, &comment.Contenu, &comment.UtilisateurID, &comment.UtilisateurNom, &comment.TopicID,
			&comment.DateCreation, &comment.DateModification, &comment.Likes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetCommentByID récupère un commentaire par ID
func (r *ForumRepository) GetCommentByID(id int) (*models.ForumComment, error) {
	query := `SELECT c.id, c.contenu, c.utilisateur_id, u.nom, c.topic_id, c.date_creation, c.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_comments c
		LEFT JOIN utilisateurs u ON c.utilisateur_id = u.id
		LEFT JOIN forum_likes l ON c.id = l.comment_id
		WHERE c.id = ?
		GROUP BY c.id`

	var comment models.ForumComment
	err := r.Db.QueryRow(query, id).Scan(&comment.ID, &comment.Contenu, &comment.UtilisateurID, &comment.UtilisateurNom, &comment.TopicID,
		&comment.DateCreation, &comment.DateModification, &comment.Likes)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateComment met à jour un commentaire
func (r *ForumRepository) UpdateComment(id int, contenu string) error {
	_, err := r.Db.Exec("UPDATE forum_comments SET contenu = ? WHERE id = ?", contenu, id)
	return err
}

// DeleteComment supprime un commentaire
func (r *ForumRepository) DeleteComment(id int) error {
	// Récupérer le topic_id avant de supprimer
	var topicID int
	r.Db.QueryRow("SELECT topic_id FROM forum_comments WHERE id = ?", id).Scan(&topicID)

	_, err := r.Db.Exec("DELETE FROM forum_comments WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Décrémenter le nombre de réponses du topic
	r.Db.Exec("UPDATE forum_topics SET nombre_reponses = nombre_reponses - 1 WHERE id = ?", topicID)

	return nil
}

// LIKES

// AddLike ajoute un like
func (r *ForumRepository) AddLike(like *models.ForumLike) error {
	query := "INSERT INTO forum_likes (utilisateur_id, topic_id, comment_id, type_like) VALUES (?, ?, ?, ?)"
	result, err := r.Db.Exec(query, like.UtilisateurID, like.TopicID, like.CommentID, like.TypeLike)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	like.ID = int(id)
	return err
}

// RemoveLike supprime un like
func (r *ForumRepository) RemoveLike(utilisateurID int, topicID *int, commentID *int) error {
	query := "DELETE FROM forum_likes WHERE utilisateur_id = ? AND topic_id <=> ? AND comment_id <=> ?"
	_, err := r.Db.Exec(query, utilisateurID, topicID, commentID)
	return err
}

// GetUserLike récupère un like d'un utilisateur
func (r *ForumRepository) GetUserLike(utilisateurID int, topicID *int, commentID *int) (*models.ForumLike, error) {
	query := "SELECT id, utilisateur_id, topic_id, comment_id, type_like, date_creation FROM forum_likes WHERE utilisateur_id = ? AND topic_id <=> ? AND comment_id <=> ?"
	var like models.ForumLike
	err := r.Db.QueryRow(query, utilisateurID, topicID, commentID).Scan(&like.ID, &like.UtilisateurID, &like.TopicID, &like.CommentID, &like.TypeLike, &like.DateCreation)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// Search recherche des topics
func (r *ForumRepository) Search(query string, limit int, offset int) ([]models.ForumTopic, error) {
	searchQuery := `SELECT t.id, t.titre, t.slug, t.description, t.contenu, t.utilisateur_id, u.nom, 
			t.categorie_id, c.nom, t.vues, t.nombre_reponses, t.epingle, t.ferme, t.date_creation, t.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_topics t
		LEFT JOIN utilisateurs u ON t.utilisateur_id = u.id
		LEFT JOIN forum_categories c ON t.categorie_id = c.id
		LEFT JOIN forum_likes l ON t.id = l.topic_id
		WHERE t.titre LIKE ? OR t.description LIKE ? OR t.contenu LIKE ?
		GROUP BY t.id
		ORDER BY t.date_creation DESC
		LIMIT ? OFFSET ?`

	searchTerm := fmt.Sprintf("%%%s%%", query)
	rows, err := r.Db.Query(searchQuery, searchTerm, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.ForumTopic
	for rows.Next() {
		var topic models.ForumTopic
		err := rows.Scan(&topic.ID, &topic.Titre, &topic.Slug, &topic.Description, &topic.Contenu,
			&topic.UtilisateurID, &topic.UtilisateurNom, &topic.CategorieID, &topic.CategorieNom,
			&topic.Vues, &topic.NombreReponses, &topic.Epingle, &topic.Ferme, &topic.DateCreation, &topic.DateModification, &topic.Likes)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}

// GetPopularTopics récupère les topics les plus populaires
func (r *ForumRepository) GetPopularTopics(limit int) ([]models.ForumTopic, error) {
	query := `SELECT t.id, t.titre, t.slug, t.description, t.contenu, t.utilisateur_id, u.nom, 
			t.categorie_id, c.nom, t.vues, t.nombre_reponses, t.epingle, t.ferme, t.date_creation, t.date_modification,
			COALESCE(COUNT(l.id), 0)
		FROM forum_topics t
		LEFT JOIN utilisateurs u ON t.utilisateur_id = u.id
		LEFT JOIN forum_categories c ON t.categorie_id = c.id
		LEFT JOIN forum_likes l ON t.id = l.topic_id
		GROUP BY t.id
		ORDER BY t.vues DESC
		LIMIT ?`

	rows, err := r.Db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []models.ForumTopic
	for rows.Next() {
		var topic models.ForumTopic
		err := rows.Scan(&topic.ID, &topic.Titre, &topic.Slug, &topic.Description, &topic.Contenu,
			&topic.UtilisateurID, &topic.UtilisateurNom, &topic.CategorieID, &topic.CategorieNom,
			&topic.Vues, &topic.NombreReponses, &topic.Epingle, &topic.Ferme, &topic.DateCreation, &topic.DateModification, &topic.Likes)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}
