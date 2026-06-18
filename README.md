# TP Forum

Forum de discussion autour de la chasse et de la peche.

## Equipe

- Dimitri Manfredonia
- Matthieu TOLISANO

## Stack

- **Langage :** Go 1.24+
- **Base de donnees :** MySQL
- **Driver SQL :** `github.com/go-sql-driver/mysql`
- **Routeur :** Gorilla Mux
- **Authentification :** JWT (HS256, `golang-jwt/jwt`)
- **Rendu :** `html/template` (SSR)
- **Frontend :** HTML, CSS, Vanilla JS

## Prerequis

- Go 1.24+
- MySQL 8+ ou compatible

## Configuration

Le projet lit un fichier `.env` a la racine du dossier `Forum/`.

Variables principales :

- `PORT` : port HTTP du serveur
- `JWT_SECRET` : cle secrete pour les tokens JWT
- `DB_DSN` : chaine de connexion MySQL

Exemple fourni dans `Forum/.env.example` :

```env
PORT=8080
JWT_SECRET=change_this_secret_key_in_production
DB_DSN=root:password@tcp(127.0.0.1:3306)/forum?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local
```

## Base De Donnees

Le schema MySQL est defini dans `Forum/migration/migration.sql`.

Le jeu de donnees de test est dans `Forum/migration/script.sql`.

## Installation

```bash
# Se placer dans le projet
cd Forum

# Copier la configuration
cp .env.example .env

# Installer les dependances
go mod tidy

# Creer la base et les tables dans MySQL
mysql -u root -p < migration/migration.sql

# Optionnel : charger les donnees de test
mysql -u root -p < migration/script.sql

# Lancer le serveur
go run .
```

Le serveur demarre sur `http://localhost:8080`.

## Comptes De Test

Le script de test recharge plusieurs comptes, dont :

- `admin` / `Admin1234!`
- `RenardRouge` / `Admin1234!`
- `PecheurDuSud` / `Admin1234!`

## Structure Du Projet

```text
Forum/
|-- main.go
|-- app/
|-- config/
|-- auth/
|-- middleware/
|-- router/
|-- controllers/
|-- services/
|-- repositories/
|-- models/
|-- dto/
|-- templates/
|-- static/
`-- migration/
```

## Fonctionnalites

- Inscription avec username et email uniques
- Connexion avec JWT
- Creation et consultation de fils de discussion
- Messages dans les fils ouverts
- Like / dislike sur les messages
- Edition et suppression par proprietaire ou admin
- Tri, pagination et filtrage
- Dashboard admin pour moderer les contenus
