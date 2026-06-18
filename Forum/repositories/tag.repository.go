package repositories

import (
	"database/sql"
	"forum/models"
)

// TagRepository gère les étiquettes (tags) en base de données.
type TagRepository struct {
	db *sql.DB
}

func InitTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

// FindAll renvoie tous les tags existants, triés par ordre alphabétique.
func (r *TagRepository) FindAll() ([]models.Tag, error) {
	rows, err := r.db.Query(`SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		rows.Scan(&t.ID, &t.Name)
		tags = append(tags, t)
	}
	return tags, nil
}

// FindOrCreate renvoie l'id d'un tag s'il existe déjà, sinon le crée puis renvoie son id.
// Pratique quand l'utilisateur tape un tag libre qui n'existe pas encore.
func (r *TagRepository) FindOrCreate(name string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM tags WHERE name = ?`, name).Scan(&id)
	if err == sql.ErrNoRows {
		// Le tag n'existe pas : on le crée.
		res, err2 := r.db.Exec(`INSERT INTO tags (name) VALUES (?)`, name)
		if err2 != nil {
			return 0, err2
		}
		lid, _ := res.LastInsertId()
		return int(lid), nil
	}
	return id, err
}

// FindByThreadID renvoie tous les tags associés à un fil de discussion donné.
func (r *TagRepository) FindByThreadID(threadID int) ([]models.Tag, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.name FROM tags t
		JOIN thread_tags tt ON tt.tag_id = t.id
		WHERE tt.thread_id = ?`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		rows.Scan(&t.ID, &t.Name)
		tags = append(tags, t)
	}
	return tags, nil
}

// SetThreadTags remplace tous les tags d'un sujet par la nouvelle liste fournie.
// On efface d'abord les anciens, puis on rajoute les nouveaux un par un.
func (r *TagRepository) SetThreadTags(threadID int, tagIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM thread_tags WHERE thread_id = ?`, threadID); err != nil {
		tx.Rollback()
		return err
	}

	for _, tid := range tagIDs {
		// MySQL utilise INSERT IGNORE pour conserver le comportement "pas d'erreur si doublon".
		if _, err := tx.Exec(`INSERT IGNORE INTO thread_tags (thread_id, tag_id) VALUES (?, ?)`, threadID, tid); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
