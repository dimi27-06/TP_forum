package models

import (
	"database/sql"
	"time"
)

type ForumUser struct {
	ID                    int          `json:"id"`
	Nom                   string       `json:"nom"`
	Email                 string       `json:"email"`
	PasswordHash          string       `json:"-"`
	PasswordSalt          string       `json:"-"`
	Bio                   string       `json:"bio,omitempty"`
	Avatar                string       `json:"avatar,omitempty"`
	Localisation          string       `json:"localisation,omitempty"`
	Role                  string       `json:"role"`
	PointsReputation      int          `json:"points_reputation"`
	Actif                 bool         `json:"actif"`
	DateCreation          time.Time    `json:"date_creation"`
	DateDerniereConnexion sql.NullTime `json:"date_derniere_connexion"`
}
