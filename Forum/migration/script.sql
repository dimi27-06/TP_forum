-- Suppression et recréation de la base
DROP DATABASE IF EXISTS forum_chasse_peche;
CREATE DATABASE forum_chasse_peche
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE forum_chasse_peche;

-- Table des utilisateurs
CREATE TABLE utilisateurs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    mot_de_passe VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar VARCHAR(255),
    localisation VARCHAR(100),
    role ENUM('user', 'moderateur', 'admin') DEFAULT 'user',
    points_reputation INT DEFAULT 0,
    actif BOOLEAN DEFAULT TRUE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_derniere_connexion TIMESTAMP NULL
) ENGINE=InnoDB;

-- Table des catégories du forum
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    slug VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(50),
    couleur VARCHAR(7),
    ordre INT DEFAULT 0,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Table des topics (fils de discussion)
CREATE TABLE topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(500),
    contenu LONGTEXT NOT NULL,
    utilisateur_id INT NOT NULL,
    categorie_id INT NOT NULL,
    vues INT DEFAULT 0,
    nombre_reponses INT DEFAULT 0,
    nombre_likes INT DEFAULT 0,
    epingle BOOLEAN DEFAULT FALSE,
    ferme BOOLEAN DEFAULT FALSE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    date_derniere_reponse TIMESTAMP NULL,
    
    INDEX idx_categorie (categorie_id),
    INDEX idx_utilisateur (utilisateur_id),
    INDEX idx_date (date_creation DESC),
    INDEX idx_epingle (epingle),
    
    CONSTRAINT fk_topics_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_topics_categories
        FOREIGN KEY (categorie_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des réponses (comments)
CREATE TABLE replies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    contenu TEXT NOT NULL,
    utilisateur_id INT NOT NULL,
    topic_id INT NOT NULL,
    numero_reponse INT,
    nombre_likes INT DEFAULT 0,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_topic (topic_id),
    INDEX idx_utilisateur (utilisateur_id),
    INDEX idx_date (date_creation DESC),
    
    CONSTRAINT fk_replies_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_replies_topics
        FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des likes
CREATE TABLE likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    utilisateur_id INT NOT NULL,
    topic_id INT NULL,
    reply_id INT NULL,
    type_like ENUM('topic', 'reply') NOT NULL,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_like (utilisateur_id, topic_id, reply_id),
    
    INDEX idx_utilisateur (utilisateur_id),
    
    CONSTRAINT fk_likes_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_likes_topics
        FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_likes_replies
        FOREIGN KEY (reply_id)
        REFERENCES replies(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des signalements
CREATE TABLE signalements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    utilisateur_id INT NOT NULL,
    topic_id INT NULL,
    reply_id INT NULL,
    raison VARCHAR(100) NOT NULL,
    description TEXT,
    statut ENUM('en_attente', 'accepte', 'rejete') DEFAULT 'en_attente',
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_signalements_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_signalements_topics
        FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_signalements_replies
        FOREIGN KEY (reply_id)
        REFERENCES replies(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des abonnements aux catégories
CREATE TABLE abonnements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    utilisateur_id INT NOT NULL,
    categorie_id INT NOT NULL,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE KEY unique_abonnement (utilisateur_id, categorie_id),
    
    CONSTRAINT fk_abonnements_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_abonnements_categories
        FOREIGN KEY (categorie_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des tags
CREATE TABLE tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Table de liaison topics_tags
CREATE TABLE topics_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    topic_id INT NOT NULL,
    tag_id INT NOT NULL,
    
    UNIQUE KEY unique_topic_tag (topic_id, tag_id),
    
    CONSTRAINT fk_topics_tags_topics
        FOREIGN KEY (topic_id)
        REFERENCES topics(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_topics_tags_tags
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Créer des indices supplémentaires pour les performances
CREATE INDEX idx_topics_slug ON topics(slug);
CREATE INDEX idx_topics_vues ON topics(vues DESC);
CREATE INDEX idx_replies_date_creation ON replies(date_creation DESC);
CREATE INDEX idx_likes_topic ON likes(topic_id);
CREATE INDEX idx_likes_reply ON likes(reply_id);