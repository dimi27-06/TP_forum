package repositories

import (
	"database/sql"
	"forum/models"
)

// ReactionRepository gère les "j'aime" / "j'aime pas" en base de données.
type ReactionRepository struct {
	db *sql.DB
}

func InitReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

// FindByUserAndMessage cherche si un membre a déjà réagi à un message précis.
// Renvoie nil (rien) s'il n'a pas encore réagi.
func (r *ReactionRepository) FindByUserAndMessage(userID, messageID int) (*models.Reaction, error) {
	row := r.db.QueryRow(`SELECT id, user_id, message_id, type, created_at FROM reactions WHERE user_id=? AND message_id=?`, userID, messageID)
	react := &models.Reaction{}
	err := row.Scan(&react.ID, &react.UserID, &react.MessageID, &react.Type, &react.CreatedAt)
	if err == sql.ErrNoRows {
		// Aucune réaction trouvée : ce n'est pas une vraie erreur.
		return nil, nil
	}
	return react, err
}

// Upsert = "Update OR Insert".
// Si le membre n'a pas encore réagi, on insère la réaction.
// S'il avait déjà réagi, on remplace simplement son ancien choix par le nouveau.
func (r *ReactionRepository) Upsert(userID, messageID int, reactionType string) error {
	_, err := r.db.Exec(`
		INSERT INTO reactions (user_id, message_id, type) VALUES (?, ?, ?)
		ON CONFLICT(user_id, message_id) DO UPDATE SET type=excluded.type`,
		userID, messageID, reactionType)
	return err
}

// Delete retire la réaction d'un membre sur un message (par exemple s'il re-clique dessus).
func (r *ReactionRepository) Delete(userID, messageID int) error {
	_, err := r.db.Exec(`DELETE FROM reactions WHERE user_id=? AND message_id=?`, userID, messageID)
	return err
}

// GetScore calcule, pour un message, le nombre de likes, de dislikes, et le score (likes - dislikes).
func (r *ReactionRepository) GetScore(messageID int) (likes, dislikes, score int) {
	r.db.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN type='like' THEN 1 ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN type='dislike' THEN 1 ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN type='like' THEN 1 WHEN type='dislike' THEN -1 ELSE 0 END),0)
		FROM reactions WHERE message_id=?`, messageID).Scan(&likes, &dislikes, &score)
	return
}
