package repositories

import (
	"database/sql"
	"errors"

	"exemple_api/models"
)

type UserRepository struct {
	Db *sql.DB
}

func InitUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) CreateUser(user *models.ForumUser) error {
	query := `
		INSERT INTO utilisateurs (
			nom, email, mot_de_passe, mot_de_passe_salt, bio, avatar, localisation,
			role, points_reputation, actif
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.Db.Exec(
		query,
		user.Nom,
		user.Email,
		user.PasswordHash,
		user.PasswordSalt,
		user.Bio,
		user.Avatar,
		user.Localisation,
		user.Role,
		user.PointsReputation,
		user.Actif,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.ForumUser, error) {
	query := `
		SELECT id, nom, email, mot_de_passe, mot_de_passe_salt, bio, avatar, localisation,
			role, points_reputation, actif, date_creation, date_derniere_connexion
		FROM utilisateurs
		WHERE email = ?
	`

	var user models.ForumUser
	err := r.Db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Nom,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordSalt,
		&user.Bio,
		&user.Avatar,
		&user.Localisation,
		&user.Role,
		&user.PointsReputation,
		&user.Actif,
		&user.DateCreation,
		&user.DateDerniereConnexion,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByIdentifier(identifier string) (*models.ForumUser, error) {
	query := `
		SELECT id, nom, email, mot_de_passe, mot_de_passe_salt, bio, avatar, localisation,
			role, points_reputation, actif, date_creation, date_derniere_connexion
		FROM utilisateurs
		WHERE email = ? OR nom = ?
		LIMIT 1
	`

	var user models.ForumUser
	err := r.Db.QueryRow(query, identifier, identifier).Scan(
		&user.ID,
		&user.Nom,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordSalt,
		&user.Bio,
		&user.Avatar,
		&user.Localisation,
		&user.Role,
		&user.PointsReputation,
		&user.Actif,
		&user.DateCreation,
		&user.DateDerniereConnexion,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.ForumUser, error) {
	query := `
		SELECT id, nom, email, mot_de_passe, mot_de_passe_salt, bio, avatar, localisation,
			role, points_reputation, actif, date_creation, date_derniere_connexion
		FROM utilisateurs
		WHERE id = ?
	`

	var user models.ForumUser
	err := r.Db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nom,
		&user.Email,
		&user.PasswordHash,
		&user.PasswordSalt,
		&user.Bio,
		&user.Avatar,
		&user.Localisation,
		&user.Role,
		&user.PointsReputation,
		&user.Actif,
		&user.DateCreation,
		&user.DateDerniereConnexion,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ExistsByEmailOrName(email string, name string) (bool, error) {
	query := `SELECT 1 FROM utilisateurs WHERE email = ? OR nom = ? LIMIT 1`
	var found int
	err := r.Db.QueryRow(query, email, name).Scan(&found)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepository) UpdateLastLogin(id int) error {
	_, err := r.Db.Exec("UPDATE utilisateurs SET date_derniere_connexion = NOW() WHERE id = ?", id)
	return err
}
