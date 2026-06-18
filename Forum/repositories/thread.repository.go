package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
)

// ThreadRepository gère les fils de discussion en base de données.
type ThreadRepository struct {
	db *sql.DB
}

func InitThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// Create enregistre un nouveau fil de discussion et renvoie son id.
func (r *ThreadRepository) Create(title, content string, userID int) (int, error) {
	res, err := r.db.Exec(
		`INSERT INTO threads (title, content, user_id) VALUES (?, ?, ?)`,
		title, content, userID,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// FindByID récupère un fil précis grâce à son id.
// La sous-requête entre parenthèses compte au passage ses messages (msg_count).
func (r *ThreadRepository) FindByID(id int) (*models.Thread, error) {
	row := r.db.QueryRow(`
		SELECT t.id, t.title, t.content, t.status, t.user_id, u.username,
		       t.created_at, t.updated_at,
		       (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) as msg_count
		FROM threads t JOIN users u ON u.id = t.user_id
		WHERE t.id = ?`, id)
	return r.scanThread(row)
}

// FindAll renvoie la liste des fils visibles, avec filtres, recherche, tri et pagination.
// Les sujets "archived" sont volontairement exclus pour les visiteurs normaux.
func (r *ThreadRepository) FindAll(page, limit int, tag, search, sort string) ([]models.Thread, int, error) {
	// On commence par une condition de base, puis on l'enrichit selon les filtres.
	whereClause := "WHERE t.status != 'archived'"
	args := []interface{}{}

	// Filtre par tag : on ne garde que les sujets qui possèdent ce tag.
	if tag != "" {
		whereClause += " AND EXISTS (SELECT 1 FROM thread_tags tt JOIN tags tg ON tg.id = tt.tag_id WHERE tt.thread_id = t.id AND tg.name = ?)"
		args = append(args, tag)
	}
	// Filtre par recherche : on regarde dans le titre OU dans les tags.
	if search != "" {
		whereClause += " AND (t.title LIKE ? OR EXISTS (SELECT 1 FROM thread_tags tt JOIN tags tg ON tg.id = tt.tag_id WHERE tt.thread_id = t.id AND tg.name LIKE ?))"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	// Count total
	// On compte d'abord le nombre total de résultats (pour la pagination).
	// On recopie les arguments car la requête de comptage utilise les mêmes filtres.
	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	r.db.QueryRow("SELECT COUNT(*) FROM threads t "+whereClause, countArgs...).Scan(&total)

	// Tri : par défaut du plus récent au plus ancien.
	orderBy := "t.created_at DESC"
	if sort == "oldest" {
		orderBy = "t.created_at ASC"
	}

	// On assemble la requête finale avec les filtres et le tri.
	query := fmt.Sprintf(`
		SELECT t.id, t.title, t.content, t.status, t.user_id, u.username,
		       t.created_at, t.updated_at,
		       (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) as msg_count
		FROM threads t JOIN users u ON u.id = t.user_id
		%s ORDER BY %s`, whereClause, orderBy)

	// On limite au besoin pour n'afficher qu'une page.
	if limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		t, err := r.scanThread(rows)
		if err == nil {
			threads = append(threads, *t)
		}
	}
	return threads, total, nil
}

// FindByUserID renvoie tous les fils créés par un utilisateur donné.
func (r *ThreadRepository) FindByUserID(userID int) ([]models.Thread, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.content, t.status, t.user_id, u.username,
		       t.created_at, t.updated_at,
		       (SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id) as msg_count
		FROM threads t JOIN users u ON u.id = t.user_id
		WHERE t.user_id = ? ORDER BY t.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var threads []models.Thread
	for rows.Next() {
		t, err := r.scanThread(rows)
		if err == nil {
			threads = append(threads, *t)
		}
	}
	return threads, nil
}

// Update modifie le titre, le contenu et le statut d'un fil.
func (r *ThreadRepository) Update(id int, title, content, status string) error {
	_, err := r.db.Exec(
		`UPDATE threads SET title=?, content=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		title, content, status, id,
	)
	return err
}

// UpdateStatus change uniquement le statut d'un fil (open / closed / archived).
func (r *ThreadRepository) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE threads SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

// Delete supprime un fil de discussion.
func (r *ThreadRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM threads WHERE id=?`, id)
	return err
}

// Même astuce que pour les messages : on accepte aussi bien QueryRow que Query,
// car les deux savent faire "Scan".
type threadScanner interface {
	Scan(dest ...interface{}) error
}

// scanThread recopie une ligne de la base dans un objet Thread.
func (r *ThreadRepository) scanThread(s threadScanner) (*models.Thread, error) {
	t := &models.Thread{}
	err := s.Scan(&t.ID, &t.Title, &t.Content, &t.Status, &t.UserID, &t.Author,
		&t.CreatedAt, &t.UpdatedAt, &t.MsgCount)
	return t, err
}
