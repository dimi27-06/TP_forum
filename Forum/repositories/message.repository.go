// Le package "repositories" contient le seul code autorisé à parler directement
// à la base de données (lire, écrire, modifier, supprimer).
package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
)

// MessageRepository garde une référence vers la base pour pouvoir l'interroger.
type MessageRepository struct {
	db *sql.DB
}

// InitMessageRepository fabrique le repository en lui donnant la base de données.
func InitMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create enregistre un nouveau message et renvoie son identifiant.
// Les "?" dans la requête sont remplacés par les valeurs : c'est plus sûr,
// ça protège contre les injections SQL (du code malveillant glissé dans un texte).
func (r *MessageRepository) Create(content string, threadID, userID int) (int, error) {
	res, err := r.db.Exec(
		`INSERT INTO messages (content, thread_id, user_id) VALUES (?, ?, ?)`,
		content, threadID, userID,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId() // on récupère l'id généré automatiquement
	return int(id), nil
}

// FindByID va chercher un seul message grâce à son identifiant.
// La requête calcule aussi au passage le nombre de likes, de dislikes et le score.
func (r *MessageRepository) FindByID(id int) (*models.Message, error) {
	row := r.db.QueryRow(`
		SELECT m.id, m.content, m.thread_id, m.user_id, u.username,
		       COALESCE(SUM(CASE WHEN r.type='like' THEN 1 ELSE 0 END),0) as likes,
		       COALESCE(SUM(CASE WHEN r.type='dislike' THEN 1 ELSE 0 END),0) as dislikes,
		       COALESCE(SUM(CASE WHEN r.type='like' THEN 1 WHEN r.type='dislike' THEN -1 ELSE 0 END),0) as score,
		       m.created_at, m.updated_at
		FROM messages m
		JOIN users u ON u.id = m.user_id
		LEFT JOIN reactions r ON r.message_id = m.id
		WHERE m.id = ?
		GROUP BY m.id`, id)
	return r.scanMessage(row)
}

// FindByThreadID renvoie les messages d'un fil de discussion, page par page.
// Elle renvoie aussi le nombre total de messages (utile pour la pagination).
func (r *MessageRepository) FindByThreadID(threadID, page, limit int, sort string, currentUserID int) ([]models.Message, int, error) {
	// On compte d'abord combien il y a de messages en tout dans ce fil.
	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE thread_id = ?`, threadID).Scan(&total)

	// On choisit l'ordre d'affichage selon ce que l'utilisateur a demandé.
	orderBy := "m.created_at DESC" // par défaut : du plus récent au plus ancien
	switch sort {
	case "oldest":
		orderBy = "m.created_at ASC" // du plus ancien au plus récent
	case "popular":
		orderBy = "score DESC, m.created_at DESC" // les plus appréciés d'abord
	}

	// On construit la requête. fmt.Sprintf insère l'ordre de tri choisi juste au-dessus.
	query := fmt.Sprintf(`
		SELECT m.id, m.content, m.thread_id, m.user_id, u.username,
		       COALESCE(SUM(CASE WHEN r.type='like' THEN 1 ELSE 0 END),0) as likes,
		       COALESCE(SUM(CASE WHEN r.type='dislike' THEN 1 ELSE 0 END),0) as dislikes,
		       COALESCE(SUM(CASE WHEN r.type='like' THEN 1 WHEN r.type='dislike' THEN -1 ELSE 0 END),0) as score,
		       m.created_at, m.updated_at
		FROM messages m
		JOIN users u ON u.id = m.user_id
		LEFT JOIN reactions r ON r.message_id = m.id
		WHERE m.thread_id = ?
		GROUP BY m.id
		ORDER BY %s`, orderBy)

	args := []interface{}{threadID}
	// Si une limite est demandée, on ne récupère que la "tranche" correspondant à la page.
	if limit > 0 {
		offset := (page - 1) * limit // combien d'éléments sauter avant la page voulue
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close() // on ferme le résultat une fois la lecture terminée

	// On parcourt chaque ligne renvoyée et on la transforme en objet Message.
	var messages []models.Message
	for rows.Next() {
		m, err := r.scanMessage(rows)
		if err == nil {
			// Charger la réaction de l'utilisateur courant
			// (pour pouvoir afficher si LUI a déjà liké ce message).
			if currentUserID > 0 {
				var rType string
				r.db.QueryRow(`SELECT type FROM reactions WHERE user_id=? AND message_id=?`, currentUserID, m.ID).Scan(&rType)
				m.UserReaction = rType
			}
			messages = append(messages, *m)
		}
	}
	return messages, total, nil
}

// Update modifie le texte d'un message et met à jour sa date de modification.
func (r *MessageRepository) Update(id int, content string) error {
	_, err := r.db.Exec(`UPDATE messages SET content=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, content, id)
	return err
}

// Delete supprime un message de la base.
func (r *MessageRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM messages WHERE id=?`, id)
	return err
}

// msgScanner est une petite astuce : QueryRow et Query renvoient des types différents,
// mais tous les deux savent faire "Scan". On accepte donc n'importe lequel des deux.
type msgScanner interface {
	Scan(dest ...interface{}) error
}

// scanMessage recopie une ligne de la base dans un objet Message bien rangé.
// L'ordre des champs doit correspondre exactement à l'ordre du SELECT.
func (r *MessageRepository) scanMessage(s msgScanner) (*models.Message, error) {
	m := &models.Message{}
	err := s.Scan(&m.ID, &m.Content, &m.ThreadID, &m.UserID, &m.Author,
		&m.Likes, &m.Dislikes, &m.Score, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}
