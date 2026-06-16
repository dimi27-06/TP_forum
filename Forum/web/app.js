const SESSION_KEY = 'forum_session_v1';
const API_BASE = '/api';

function readSession() {
  try {
    return JSON.parse(localStorage.getItem(SESSION_KEY) || '{}');
  } catch {
    return {};
  }
}

function writeSession(session) {
  localStorage.setItem(SESSION_KEY, JSON.stringify(session));
}

function clearSession() {
  localStorage.removeItem(SESSION_KEY);
}

function isMember() {
  const session = readSession();
  return Boolean(session.mode === 'member' && session.token && session.userId);
}

function getHeaders(extraHeaders = {}) {
  const session = readSession();
  const headers = { ...extraHeaders };

  if (isMember()) {
    headers.Authorization = `Bearer ${session.token}`;
    headers['X-User-ID'] = String(session.userId);
  }

  return headers;
}

async function fetchJson(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
  });

  const contentType = response.headers.get('content-type') || '';
  const payload = contentType.includes('application/json') ? await response.json() : await response.text();

  if (!response.ok) {
    const message = payload && typeof payload === 'object' && 'error' in payload ? payload.error : 'Request failed';
    throw new Error(message);
  }

  return payload;
}

function formatDate(value) {
  if (!value) {
    return '';
  }

  const date = new Date(value);
  return new Intl.DateTimeFormat('fr-FR', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date);
}

function setText(selector, value) {
  const target = document.querySelector(selector);
  if (target) {
    target.textContent = value;
  }
}

function setHTML(selector, value) {
  const target = document.querySelector(selector);
  if (target) {
    target.innerHTML = value;
  }
}

function showMessage(selector, message, tone = 'notice') {
  const target = document.querySelector(selector);
  if (!target) {
    return;
  }

  target.textContent = message;
  target.dataset.tone = tone;
}

async function hydrateSession() {
  const session = readSession();
  const chip = document.querySelector('[data-session-chip]');
  const loginLink = document.querySelector('[data-login-link]');
  const registerLink = document.querySelector('[data-register-link]');
  const logoutLink = document.querySelector('[data-logout-link]');

  if (!chip) {
    return session;
  }

  if (session.mode === 'member' && session.token) {
    try {
      const me = await fetchJson(`${API_BASE}/me`, {
        headers: { Authorization: `Bearer ${session.token}` },
      });

      const hydrated = {
        mode: 'member',
        token: session.token,
        userId: me.user_id,
        role: me.role,
      };

      writeSession(hydrated);
      chip.textContent = `Connecte: ${me.role} #${me.user_id}`;
      chip.classList.remove('muted');
      if (loginLink) loginLink.classList.add('hidden');
      if (registerLink) registerLink.classList.add('hidden');
      if (logoutLink) logoutLink.classList.remove('hidden');
      return hydrated;
    } catch {
      clearSession();
    }
  }

  chip.textContent = 'Mode invite';
  chip.classList.add('muted');
  if (loginLink) loginLink.classList.remove('hidden');
  if (registerLink) registerLink.classList.remove('hidden');
  if (logoutLink) logoutLink.classList.add('hidden');
  return { mode: 'guest' };
}

async function handleLogin(event) {
  event.preventDefault();

  const form = event.currentTarget;
  const errorBox = document.querySelector('[data-login-error]');
  const submitButton = form.querySelector('button[type="submit"]');
  const username = form.username.value.trim();
  const password = form.password.value.trim();

  if (!username || !password) {
    showMessage('[data-login-error]', 'Remplis le nom d utilisateur et le mot de passe.', 'error');
    return;
  }

  submitButton.disabled = true;
  submitButton.textContent = 'Connexion...';

  try {
    const auth = await fetchJson(`${API_BASE}/login`, {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });

    const me = await fetchJson(`${API_BASE}/me`, {
      headers: { Authorization: `${auth.type} ${auth.access_token}` },
    });

    writeSession({
      mode: 'member',
      token: auth.access_token,
      userId: me.user_id,
      role: me.role,
    });

    window.location.href = '/forum';
  } catch (error) {
    showMessage('[data-login-error]', error.message || 'Connexion impossible.', 'error');
  } finally {
    submitButton.disabled = false;
    submitButton.textContent = 'Se connecter';
  }
}

async function handleRegister(event) {
  event.preventDefault();

  const form = event.currentTarget;
  const submitButton = form.querySelector('button[type="submit"]');

  submitButton.disabled = true;
  submitButton.textContent = 'Creation...';

  try {
    const auth = await fetchJson(`${API_BASE}/register`, {
      method: 'POST',
      body: JSON.stringify({
        nom: form.nom.value.trim(),
        email: form.email.value.trim(),
        password: form.password.value.trim(),
        localisation: form.localisation.value.trim(),
        bio: form.bio.value.trim(),
      }),
    });

    const me = await fetchJson(`${API_BASE}/me`, {
      headers: { Authorization: `${auth.type} ${auth.access_token}` },
    });

    writeSession({
      mode: 'member',
      token: auth.access_token,
      userId: me.user_id,
      role: me.role,
    });

    window.location.href = '/forum';
  } catch (error) {
    showMessage('[data-register-error]', error.message || 'Inscription impossible.', 'error');
  } finally {
    submitButton.disabled = false;
    submitButton.textContent = 'Creer mon compte';
  }
}

function renderTopics(topics) {
  const container = document.querySelector('[data-topics]');
  if (!container) {
    return;
  }

  if (!topics.length) {
    container.innerHTML = `
      <div class="empty">
        <h3>Aucun message pour le moment</h3>
        <p>Crée le premier topic pour lancer la discussion.</p>
      </div>
    `;
    return;
  }

  container.innerHTML = topics.map((topic) => `
    <article class="topic">
      <div class="topic-top">
        <span class="category-pill">${topic.categorie_nom || 'Forum'}</span>
        <span class="topic-meta">${formatDate(topic.date_creation)}</span>
      </div>
      <h3><a href="/forum/topic/${topic.id}">${topic.titre}</a></h3>
      <p>${topic.description || topic.contenu || ''}</p>
      <div class="topic-actions">
        <span class="chip">${topic.nombre_reponses || 0} replies</span>
        <span class="chip">${topic.likes || 0} likes</span>
        <a class="btn small" href="/forum/topic/${topic.id}">Open</a>
        <button class="btn small" data-like-topic="${topic.id}" type="button">Like</button>
      </div>
    </article>
  `).join('');

  container.querySelectorAll('[data-like-topic]').forEach((button) => {
    button.addEventListener('click', async () => {
      if (!isMember()) {
        window.location.href = '/login';
        return;
      }

      const topicId = button.getAttribute('data-like-topic');
      try {
        await fetchJson(`${API_BASE}/forum/topics/${topicId}/like`, {
          method: 'POST',
          headers: getHeaders(),
        });
        await loadForumTopics();
      } catch (error) {
        alert(error.message);
      }
    });
  });
}

function renderCategories(categories) {
  const select = document.querySelector('[data-category-select]');
  const sidebar = document.querySelector('[data-categories]');

  if (select) {
    const options = ['<option value="">Choose a category</option>']
      .concat(categories.map((category) => `<option value="${category.id}">${category.nom}</option>`))
      .join('');
    select.innerHTML = options;
  }

  if (sidebar) {
    sidebar.innerHTML = categories.map((category) => `
      <li>
        <span>${category.nom}</span>
        <strong>${category.slug || ''}</strong>
      </li>
    `).join('');
  }
}

async function loadForumTopics() {
  const container = document.querySelector('[data-topics]');
  if (!container) {
    return [];
  }

  try {
    const topics = await fetchJson(`${API_BASE}/forum/topics?limit=30`);
    renderTopics(Array.isArray(topics) ? topics : []);
    return Array.isArray(topics) ? topics : [];
  } catch (error) {
    container.innerHTML = `
      <div class="empty">
        <h3>Impossible de charger le feed</h3>
        <p>${error.message}</p>
      </div>
    `;
    return [];
  }
}

async function loadForumCategories() {
  const categories = await fetchJson(`${API_BASE}/forum/categories`);
  renderCategories(Array.isArray(categories) ? categories : []);
  return Array.isArray(categories) ? categories : [];
}

async function handleCreateTopic(event) {
  event.preventDefault();

  if (!isMember()) {
    window.location.href = '/login';
    return;
  }

  const form = event.currentTarget;
  const submitButton = form.querySelector('button[type="submit"]');
  const status = document.querySelector('[data-compose-status]');

  submitButton.disabled = true;
  submitButton.textContent = 'Publication...';

  try {
    await fetchJson(`${API_BASE}/forum/topics`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({
        titre: form.titre.value.trim(),
        description: form.description.value.trim(),
        contenu: form.contenu.value.trim(),
        categorie_id: Number(form.categorie_id.value),
      }),
    });

    form.reset();
    showMessage('[data-compose-status]', 'Message publie avec succes.', 'success');
    await loadForumTopics();
  } catch (error) {
    showMessage('[data-compose-status]', error.message, 'error');
  } finally {
    submitButton.disabled = false;
    submitButton.textContent = 'Publier le message';
  }
}

async function loadForumPage() {
  const guestMode = new URLSearchParams(window.location.search).get('mode') === 'guest';
  if (guestMode) {
    writeSession({ mode: 'guest' });
  }

  const session = await hydrateSession();

  const composer = document.querySelector('[data-composer]');
  const guestBanner = document.querySelector('[data-guest-banner]');
  const welcome = document.querySelector('[data-welcome]');

  if (composer) {
    composer.setAttribute('aria-hidden', isMember() ? 'false' : 'true');
  }

  if (guestBanner) {
    guestBanner.classList.toggle('hidden', isMember());
  }

  if (welcome) {
    if (isMember()) {
      welcome.textContent = `Bienvenue ${readSession().role || 'member'}`;
    } else {
      welcome.textContent = 'Mode invite';
    }
  }

  const createForm = document.querySelector('[data-create-topic-form]');
  if (createForm) {
    createForm.addEventListener('submit', handleCreateTopic);
  }

  if (document.querySelector('[data-topics]')) {
    await loadForumCategories();
    await loadForumTopics();
  }

  const guestButton = document.querySelector('[data-guest-button]');
  if (guestButton) {
    guestButton.addEventListener('click', () => {
      writeSession({ mode: 'guest' });
      window.location.href = '/forum?mode=guest';
    });
  }

  return session;
}

async function loadTopicPage() {
  await hydrateSession();

  const topicId = document.body.getAttribute('data-topic-id');
  if (!topicId) {
    return;
  }

  const topicWrap = document.querySelector('[data-topic-view]');
  const commentsWrap = document.querySelector('[data-comments]');
  const replyForm = document.querySelector('[data-reply-form]');

  try {
    const topic = await fetchJson(`${API_BASE}/forum/topics/${topicId}`);
    const comments = await fetchJson(`${API_BASE}/forum/topics/${topicId}/comments`);

    if (topicWrap) {
      topicWrap.innerHTML = `
        <div class="topic">
          <div class="topic-top">
            <span class="category-pill">${topic.categorie_nom || 'Forum'}</span>
            <span class="topic-meta">${formatDate(topic.date_creation)}</span>
          </div>
          <h1>${topic.titre}</h1>
          <p>${topic.contenu || topic.description || ''}</p>
          <div class="topic-actions">
            <span class="chip">${topic.nombre_reponses || 0} replies</span>
            <span class="chip">${topic.likes || 0} likes</span>
            <button class="btn small" data-like-topic="${topic.id}" type="button">Like</button>
          </div>
        </div>
      `;
    }

    if (topicWrap) {
      const likeButton = topicWrap.querySelector('[data-like-topic]');
      if (likeButton) {
        likeButton.addEventListener('click', async () => {
          if (!isMember()) {
            window.location.href = '/login';
            return;
          }

          try {
            await fetchJson(`${API_BASE}/forum/topics/${topicId}/like`, {
              method: 'POST',
              headers: getHeaders(),
            });
            await loadTopicPage();
          } catch (error) {
            alert(error.message);
          }
        });
      }
    }

    if (commentsWrap) {
      if (!Array.isArray(comments) || comments.length === 0) {
        commentsWrap.innerHTML = '<div class="empty"><h3>No replies yet</h3><p>Be the first to answer.</p></div>';
      } else {
        commentsWrap.innerHTML = comments.map((comment) => `
          <article class="comment">
            <div class="comment-top">
              <strong>${comment.utilisateur_nom || `User ${comment.utilisateur_id}`}</strong>
              <span class="meta">${formatDate(comment.date_creation)}</span>
            </div>
            <p>${comment.contenu}</p>
          </article>
        `).join('');
      }
    }

    if (replyForm) {
      replyForm.classList.toggle('hidden', !isMember());
      if (!isMember()) {
        const notice = document.querySelector('[data-reply-notice]');
        if (notice) {
          notice.classList.remove('hidden');
        }
      }
    }
  } catch (error) {
    if (topicWrap) {
      topicWrap.innerHTML = `<div class="empty"><h3>Topic not found</h3><p>${error.message}</p></div>`;
    }
  }
}

async function handleReply(event) {
  event.preventDefault();

  if (!isMember()) {
    window.location.href = '/login';
    return;
  }

  const topicId = document.body.getAttribute('data-topic-id');
  const form = event.currentTarget;
  const submitButton = form.querySelector('button[type="submit"]');
  const status = document.querySelector('[data-reply-status]');

  submitButton.disabled = true;
  submitButton.textContent = 'Envoi...';

  try {
    await fetchJson(`${API_BASE}/forum/topics/${topicId}/comments`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ contenu: form.contenu.value.trim() }),
    });

    form.reset();
    showMessage('[data-reply-status]', 'Reply sent.', 'success');
    await loadTopicPage();
  } catch (error) {
    showMessage('[data-reply-status]', error.message, 'error');
  } finally {
    submitButton.disabled = false;
    submitButton.textContent = 'Envoyer la reponse';
  }
}

function initHeaderButtons() {
  const logoutButton = document.querySelector('[data-logout-link]');
  if (logoutButton) {
    logoutButton.addEventListener('click', () => {
      clearSession();
      window.location.href = '/';
    });
  }
}

document.addEventListener('DOMContentLoaded', async () => {
  initHeaderButtons();

  const page = document.body.getAttribute('data-page');
  if (page === 'login') {
    document.querySelector('[data-login-form]')?.addEventListener('submit', handleLogin);
    await hydrateSession();
    if (isMember()) {
      window.location.href = '/forum';
    }
  }

  if (page === 'register') {
    document.querySelector('[data-register-form]')?.addEventListener('submit', handleRegister);
    await hydrateSession();
    if (isMember()) {
      window.location.href = '/forum';
    }
  }

  if (page === 'forum') {
    await loadForumPage();
  }

  if (page === 'topic') {
    document.querySelector('[data-reply-form]')?.addEventListener('submit', handleReply);
    await loadTopicPage();
  }
});
