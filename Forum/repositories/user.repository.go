package repositories

import (
	"database/sql"
	"forum/models"
)

// UserRepository gère les utilisateurs en base de données.
type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create enregistre un nouvel utilisateur et renvoie son identifiant.
// Le mot de passe reçu ici est déjà haché (le service s'en est chargé avant).
func (r *UserRepository) Create(username, email, password string) (int, error) {
	res, err := r.db.Exec(
		`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`,
		username, email, password,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// FindByID retrouve un utilisateur grâce à son numéro.
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, email, password, role, banned, created_at FROM users WHERE id = ?`, id)
	return r.scanUser(row)
}

// FindByUsername retrouve un utilisateur grâce à son pseudo.
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, email, password, role, banned, created_at FROM users WHERE username = ?`, username)
	return r.scanUser(row)
}

// FindByEmail retrouve un utilisateur grâce à son email.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, email, password, role, banned, created_at FROM users WHERE email = ?`, email)
	return r.scanUser(row)
}

// FindByIdentifier retrouve un utilisateur par son pseudo OU son email.
// Utile à la connexion, où l'on accepte les deux.
func (r *UserRepository) FindByIdentifier(identifier string) (*models.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password, role, banned, created_at FROM users WHERE username = ? OR email = ?`,
		identifier, identifier,
	)
	return r.scanUser(row)
}

// FindAll renvoie tous les utilisateurs, du plus récemment inscrit au plus ancien.
func (r *UserRepository) FindAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, username, email, password, role, banned, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err == nil {
			users = append(users, *u)
		}
	}
	return users, nil
}

// ExistsByUsername indique si un pseudo est déjà utilisé (vrai/faux).
// Pratique pour empêcher deux comptes d'avoir le même pseudo.
func (r *UserRepository) ExistsByUsername(username string) bool {
	var count int
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	return count > 0
}

// ExistsByEmail indique si un email est déjà utilisé.
func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&count)
	return count > 0
}

// SetBanned bannit (true) ou débannit (false) un utilisateur.
// En base, on stocke 1 pour banni et 0 pour autorisé.
func (r *UserRepository) SetBanned(userID int, banned bool) error {
	val := 0
	if banned {
		val = 1
	}
	_, err := r.db.Exec(`UPDATE users SET banned = ? WHERE id = ?`, val, userID)
	return err
}

// Même astuce que pour les autres repositories : on accepte aussi bien QueryRow que Query.
type userScanner interface {
	Scan(dest ...interface{}) error
}

// scanUser recopie une ligne de la base dans un objet User bien rangé.
func (r *UserRepository) scanUser(s userScanner) (*models.User, error) {
	u := &models.User{}
	err := s.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.Banned, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
