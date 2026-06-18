# TP_forum
# Forum Chasse & Pêche

Plateforme communautaire de discussion dédiée à la chasse et à la pêche.

## Équipe

- Dimitri Manfredonia
- Matthieu TOLISANO

## Technologies

- **Langage :** Go 1.22
- **Base de données :** SQLite (via `mattn/go-sqlite3`)
- **Routeur :** Gorilla Mux
- **Authentification :** JWT (HS256, `golang-jwt/jwt`)
- **Rendu :** html/template (SSR)
- **Frontend :** HTML + CSS + Vanilla JS

## Prérequis

- Go 1.22+
- GCC (requis par `go-sqlite3` pour la compilation CGO)

## Installation

```bash
# Cloner le dépôt
git clone <url-du-repo>
cd forum

# Copier la configuration
cp .env.example .env
# Éditer .env si nécessaire (JWT_SECRET, PORT, DB_PATH)

# Télécharger les dépendances
go mod tidy

# Lancer le serveur
go run main.go
```

Le serveur démarre sur **http://localhost:8080**

La base de données est créée automatiquement au premier démarrage avec des données de test.

## Comptes de test

| Utilisateur | Email | Mot de passe | Rôle |
|---|---|---|---|
| admin | admin@forum.fr | *Admin1234!* | admin |
| chasseur42 | chasseur42@forum.fr | *Admin1234!* | user |
| pecheur_du_sud | pecheur@forum.fr | *Admin1234!* | user |

> Pour créer votre propre compte, utilisez le formulaire d'inscription.
> Mot de passe requis : 12 caractères minimum, 1 majuscule, 1 caractère spécial.

## Structure du projet

```
forum/
├── main.go              # Point d'entrée
├── app/app.go           # Assemblage des dépendances
├── config/              # Chargement .env, connexion DB
├── auth/                # JWT, hash SHA-512, validation
├── middleware/          # Auth, rôles
├── router/              # Définition des routes
├── controllers/         # Réception HTTP, validation
├── services/            # Logique métier
├── repositories/        # Accès SQL
├── models/              # Structs Go
├── dto/                 # Objets de transfert
├── templates/           # Pages HTML (SSR)
├── static/              # CSS, JS, images
├── database/            # Migration auto
└── migration/           # Scripts SQL de référence
```

## Fonctionnalités

- FT-1 : Inscription (username unique, email unique, SHA-512, règles mot de passe)
- FT-2 : Connexion par username ou email + JWT
- FT-3 : Création de fil de discussion avec tags et statut
- FT-4 : Consultation des fils (open/closed visibles, archived masqués)
- FT-5 : Publication de messages dans les fils ouverts
- FT-6 : Like / Dislike sur les messages (score de popularité)
- FT-7 : Modification et suppression (propriétaire ou admin)
- FT-8 : Tri des messages (récent, ancien, popularité)
- FT-9 : Pagination (10 / 20 / 30 / tout)
- FT-10 : Filtrage par tag/catégorie
- FT-11 : Recherche par titre ou tag
- FT-12 : Dashboard admin (ban/unban, statut des fils, suppression)
