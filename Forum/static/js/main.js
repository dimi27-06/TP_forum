// ==========================================
// Forum Chasse & Pêche – JavaScript minimal
// Utilisé uniquement pour améliorer l'UX
// (UX = "expérience utilisateur", le confort d'utilisation)
// ==========================================

/**
 * FT-6 : Like / Dislike via fetch (pas de rechargement de page)
 *
 * Le but : quand on clique sur 👍 ou 👎, on prévient le serveur et on met à jour
 * les compteurs SANS recharger toute la page. C'est plus fluide.
 */
function react(messageID, type) {
    // fetch envoie une requête au serveur en arrière-plan, sans quitter la page.
    fetch('/messages/' + messageID + '/react/' + type, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
    })
    // 1er ".then" : on reçoit la réponse du serveur.
    .then(function(resp) {
        // 401 = pas connecté : on redirige vers la page de connexion.
        if (resp.status === 401) {
            window.location.href = '/login';
            return null;
        }
        if (!resp.ok) return null;   // une autre erreur : on abandonne
        return resp.json();          // sinon on transforme la réponse JSON en objet utilisable
    })
    // 2e ".then" : on a les nouvelles données (likes, dislikes, score...).
    .then(function(data) {
        if (!data) return;

        // On retrouve le bloc de réactions du bon message grâce à son data-msg-id.
        var container = document.querySelector('.reactions[data-msg-id="' + messageID + '"]');
        if (!container) return;

        // On met à jour les chiffres affichés.
        container.querySelector('.likes-count').textContent    = data.likes;
        container.querySelector('.dislikes-count').textContent = data.dislikes;
        container.querySelector('.score-val').textContent      = data.score;

        // On met en surbrillance le bouton correspondant au choix actuel de l'utilisateur.
        container.querySelectorAll('.react-btn').forEach(function(btn) {
            btn.classList.remove('active'); // on enlève d'abord la surbrillance partout
            if (btn.dataset.type === data.user_reaction) {
                btn.classList.add('active'); // puis on l'ajoute sur le bon bouton
            }
        });
    })
    // ".catch" : en cas de souci réseau, on note l'erreur dans la console du navigateur.
    .catch(function(err) {
        console.error('Erreur réaction :', err);
    });
}

// Marque le lien nav actif
// Au chargement de la page, on repère le lien du menu qui correspond à la page actuelle
// et on lui ajoute la classe "active" (pour le souligner / le colorer).
document.addEventListener('DOMContentLoaded', function() {
    var path = window.location.pathname; // l'adresse de la page actuelle
    document.querySelectorAll('.main-nav a').forEach(function(a) {
        if (a.getAttribute('href') === path) {
            a.classList.add('active');
        }
    });
});
