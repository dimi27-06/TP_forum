DROP DATABASE IF EXISTS forum_chasse_peche;
CREATE DATABASE forum_chasse_peche
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE forum_chasse_peche;

CREATE TABLE utilisateurs (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nom VARCHAR(100) NOT NULL UNIQUE,
  email VARCHAR(150) NOT NULL UNIQUE,
  mot_de_passe VARCHAR(64) NOT NULL,
  mot_de_passe_salt VARCHAR(64) NOT NULL,
  bio TEXT NULL,
  avatar VARCHAR(255) NULL,
  localisation VARCHAR(120) NULL,
  role ENUM('user', 'moderateur', 'admin') NOT NULL DEFAULT 'user',
  points_reputation INT NOT NULL DEFAULT 0,
  actif BOOLEAN NOT NULL DEFAULT TRUE,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_derniere_connexion TIMESTAMP NULL DEFAULT NULL,
  INDEX idx_utilisateurs_email (email),
  INDEX idx_utilisateurs_role (role)
) ENGINE=InnoDB;

CREATE TABLE forum_categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nom VARCHAR(100) NOT NULL UNIQUE,
  description TEXT NULL,
  slug VARCHAR(120) NOT NULL UNIQUE,
  icon VARCHAR(50) NULL,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_forum_categories_slug (slug)
) ENGINE=InnoDB;

CREATE TABLE forum_topics (
  id INT AUTO_INCREMENT PRIMARY KEY,
  titre VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  description VARCHAR(500) NULL,
  contenu LONGTEXT NOT NULL,
  utilisateur_id INT NOT NULL,
  categorie_id INT NOT NULL,
  vues INT NOT NULL DEFAULT 0,
  nombre_reponses INT NOT NULL DEFAULT 0,
  epingle BOOLEAN NOT NULL DEFAULT FALSE,
  ferme BOOLEAN NOT NULL DEFAULT FALSE,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modification TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  date_derniere_reponse TIMESTAMP NULL DEFAULT NULL,
  INDEX idx_forum_topics_categorie (categorie_id),
  INDEX idx_forum_topics_utilisateur (utilisateur_id),
  INDEX idx_forum_topics_date (date_creation DESC),
  INDEX idx_forum_topics_epingle (epingle),
  INDEX idx_forum_topics_slug (slug),
  CONSTRAINT fk_forum_topics_utilisateurs
    FOREIGN KEY (utilisateur_id) REFERENCES utilisateurs(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_forum_topics_categories
    FOREIGN KEY (categorie_id) REFERENCES forum_categories(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE forum_comments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  contenu TEXT NOT NULL,
  utilisateur_id INT NOT NULL,
  topic_id INT NOT NULL,
  nombre_likes INT NOT NULL DEFAULT 0,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date_modification TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_forum_comments_topic (topic_id),
  INDEX idx_forum_comments_utilisateur (utilisateur_id),
  INDEX idx_forum_comments_date (date_creation DESC),
  CONSTRAINT fk_forum_comments_utilisateurs
    FOREIGN KEY (utilisateur_id) REFERENCES utilisateurs(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_forum_comments_topics
    FOREIGN KEY (topic_id) REFERENCES forum_topics(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE forum_likes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  utilisateur_id INT NOT NULL,
  topic_id INT NULL,
  comment_id INT NULL,
  type_like ENUM('topic', 'comment') NOT NULL,
  date_creation TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_forum_like_target CHECK (
    (topic_id IS NOT NULL AND comment_id IS NULL)
    OR (topic_id IS NULL AND comment_id IS NOT NULL)
  ),
  UNIQUE KEY unique_topic_like (utilisateur_id, topic_id),
  UNIQUE KEY unique_comment_like (utilisateur_id, comment_id),
  INDEX idx_forum_likes_utilisateur (utilisateur_id),
  INDEX idx_forum_likes_topic (topic_id),
  INDEX idx_forum_likes_comment (comment_id),
  CONSTRAINT fk_forum_likes_utilisateurs
    FOREIGN KEY (utilisateur_id) REFERENCES utilisateurs(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_forum_likes_topics
    FOREIGN KEY (topic_id) REFERENCES forum_topics(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_forum_likes_comments
    FOREIGN KEY (comment_id) REFERENCES forum_comments(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;
