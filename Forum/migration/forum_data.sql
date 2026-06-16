USE forum_chasse_peche;

INSERT INTO utilisateurs (
  nom, email, mot_de_passe, mot_de_passe_salt, bio, localisation, role, points_reputation, actif
) VALUES
('admin', 'admin@forum.local', SHA2(CONCAT('forum-demo-salt-1', ':', 'password'), 256), 'forum-demo-salt-1', 'Admin du forum chasse et peche', 'France', 'admin', 1000, TRUE),
('camille', 'camille@forum.local', SHA2(CONCAT('forum-demo-salt-2', ':', 'peche123'), 256), 'forum-demo-salt-2', 'Passionnee de peche en eau douce', 'Bourgogne', 'user', 220, TRUE),
('leo', 'leo@forum.local', SHA2(CONCAT('forum-demo-salt-3', ':', 'chasse123'), 256), 'forum-demo-salt-3', 'Chasseur en foret et membre actif', 'Sologne', 'user', 180, TRUE);

INSERT INTO forum_categories (nom, description, slug, icon) VALUES
('Chasse generale', 'Discussions sur la chasse, les techniques et les saisons', 'chasse-generale', 'chasse'),
('Peche generale', 'Discussions sur la peche, les montages et les spots', 'peche-generale', 'peche'),
('Materiel', 'Armes, cannes, moulinets, sacs et equipements', 'materiel', 'materiel'),
('Debutants', 'Questions pour commencer sans se tromper', 'debutants', 'debutants'),
('Spots', 'Partage de lieux, riviere, foret et bord de mer', 'spots', 'spots');

INSERT INTO forum_topics (
  titre, slug, description, contenu, utilisateur_id, categorie_id, vues, nombre_reponses, epingle, ferme
) VALUES
(
  'Quel leurre pour la truite en riviere rapide ?',
  'quel-leurre-pour-la-truite-en-riviere-rapide',
  'Conseils de peche en eau vive',
  'Je cherche des retours sur les leurres qui fonctionnent le mieux en riviere rapide. Vous partez sur quelle taille et quelle couleur ?',
  2,
  2,
  48,
  3,
  FALSE,
  FALSE
),
(
  'Sortie de chasse au lever du jour',
  'sortie-de-chasse-au-lever-du-jour',
  'Organisation et securite avant la sortie',
  'Quels sont vos checklists avant une sortie de chasse le matin ? Je veux construire une routine simple et sure.',
  3,
  1,
  36,
  2,
  FALSE,
  FALSE
),
(
  'Quel sac leger pour garder les mains libres ?',
  'quel-sac-leger-pour-garder-les-mains-libres',
  'Discussion materiel',
  'Je cherche un sac ou gilet pratique pour les longues sorties. Vous utilisez quoi au quotidien ?',
  1,
  3,
  29,
  1,
  FALSE,
  FALSE
);

INSERT INTO forum_comments (contenu, utilisateur_id, topic_id, nombre_likes) VALUES
('Pour la truite, je commence souvent avec des petits leurres coulant lentement.', 2, 1, 4),
('En riviere rapide, les coloris naturels marchent bien chez moi.', 3, 1, 2),
('La veille, je verifie toujours le terrain et la meteo.', 1, 2, 3),
('Un gilet de peche avec des poches fermes peut vraiment faire la difference.', 2, 3, 1);

INSERT INTO forum_likes (utilisateur_id, topic_id, comment_id, type_like) VALUES
(1, 1, NULL, 'topic'),
(2, 1, NULL, 'topic'),
(3, 1, NULL, 'topic'),
(1, 2, NULL, 'topic'),
(2, 3, NULL, 'topic'),
(3, NULL, 1, 'comment'),
(1, NULL, 2, 'comment');
