-- Migration pour le forum de chasse et pêche
USE boutique_exemple_code;

-- Table des catégories du forum
CREATE TABLE IF NOT EXISTS forum_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    slug VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(50),
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Table des topics (fils de discussion)
CREATE TABLE IF NOT EXISTS forum_topics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titre VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    contenu LONGTEXT NOT NULL,
    utilisateur_id INT NOT NULL,
    categorie_id INT NOT NULL,
    vues INT DEFAULT 0,
    nombre_reponses INT DEFAULT 0,
    epingle BOOLEAN DEFAULT FALSE,
    ferme BOOLEAN DEFAULT FALSE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_topics_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_topics_categories
        FOREIGN KEY (categorie_id)
        REFERENCES forum_categories(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des réponses (comments)
CREATE TABLE IF NOT EXISTS forum_comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    contenu TEXT NOT NULL,
    utilisateur_id INT NOT NULL,
    topic_id INT NOT NULL,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_comments_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_comments_topics
        FOREIGN KEY (topic_id)
        REFERENCES forum_topics(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- Table des likes (j'aime)
CREATE TABLE IF NOT EXISTS forum_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    utilisateur_id INT NOT NULL,
    topic_id INT NULL,
    comment_id INT NULL,
    type_like ENUM('topic', 'comment') NOT NULL,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_likes_utilisateurs
        FOREIGN KEY (utilisateur_id)
        REFERENCES utilisateurs(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_likes_topics
        FOREIGN KEY (topic_id)
        REFERENCES forum_topics(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_likes_comments
        FOREIGN KEY (comment_id)
        REFERENCES forum_comments(id)
        ON DELETE CASCADE,
    
    UNIQUE KEY unique_topic_like (utilisateur_id, topic_id, comment_id)
) ENGINE=InnoDB;

-- Indices pour améliorer les performances
CREATE INDEX idx_topics_categorie ON forum_topics(categorie_id);
CREATE INDEX idx_topics_utilisateur ON forum_topics(utilisateur_id);
CREATE INDEX idx_topics_date ON forum_topics(date_creation DESC);
CREATE INDEX idx_comments_topic ON forum_comments(topic_id);
CREATE INDEX idx_comments_utilisateur ON forum_comments(utilisateur_id);
CREATE INDEX idx_comments_date ON forum_comments(date_creation DESC);
CREATE INDEX idx_likes_utilisateur ON forum_likes(utilisateur_id);
