# Forum Chasse & Peche

Plateforme de discussion dédiée à la chasse et à la pêche, avec sujets, réponses, tags, réactions et espace admin.

## Équipe

- Dimi
- Matthieu Tolisano

## Tech

- Go 1.24
- MySQL
- Gorilla Mux
- JWT
- `html/template`
- HTML, CSS, JavaScript

## Prérequis

- Go 1.24+
- MySQL 8+ ou compatible
- Un client SQL comme DBeaver si besoin

## Mise en route

```bash
git clone <url-du-repo>
cd forum

copy .env.example .env
go mod tidy

mysql -u root -p < migration/migration.sql
mysql -u root -p < migration/script.sql

go run main.go
```

Le serveur tourne sur `http://localhost:8080`.

## Base de données

- `migration/migration.sql` crée la base et les tables
- `migration/script.sql` ajoute les données de départ

Le seed contient 16 membres, 16 fils de discussion et au moins 2 réponses par fil.

## Comptes de test

| Utilisateur | Email | Mot de passe | Rôle |
|---|---|---|---|
| admin | admin@forum.fr | `Admin1234!` | admin |
| chasseur42 | chasseur42@forum.fr | `Admin1234!` | user |
| pecheur_du_sud | pecheur@forum.fr | `Admin1234!` | user |

Pour créer un compte, passe par l’inscription. Le mot de passe doit faire au moins 12 caractères, avec une majuscule et un caractère spécial.

## Structure

```text
forum/
├── main.go
├── app/
├── auth/
├── config/
├── controllers/
├── dto/
├── middleware/
├── migration/
├── models/
├── repositories/
├── router/
├── services/
├── static/
└── templates/
```

## Fonctionnalités

- Inscription et connexion
- Création de fils avec tags
- Lecture des fils et des réponses
- Réactions like / dislike
- Pagination, tri et recherche
- Administration des utilisateurs et des fils
