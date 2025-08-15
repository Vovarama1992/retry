(function () {
  console.log('[tracker] start');

  // --- utils
  function getCookie(name) {
    const m = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return m ? decodeURIComponent(m[2]) : null;
  }
  function setCookie(name, value, days = 365) {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; path=/; expires=${expires}`;
    // console.log(`[tracker] cookie ${name}=${value}`);
  }
  function generateId(prefix) {
    return `${prefix}_${Math.random().toString(36).slice(2)}_${Date.now()}`;
  }
  function nowIso() { return new Date().toISOString(); }
  function nowMs() { return Date.now(); }

  // --- source
  function getSource() {
    const params = new URLSearchParams(window.location.search);
    if (params.has('utm_source')) return 'utm:' + params.get('utm_source');
    if (document.referrer) {
      try { return 'ref:' + new URL(document.referrer).hostname; }
      catch { /* ignore */ }
    }
    return 'direct';
  }

  // --- visit/session bootstrap
  let visitId = getCookie('visit_id');
  let source  = getCookie('visit_source');
  if (!visitId) {
    visitId = generateId('visit');
    source = getSource();
    setCookie('visit_id', visitId);
    setCookie('visit_source', source);
  }

  let sessionId = getCookie('session_id');
  if (!sessionId) {
    sessionId = generateId('session');
    setCookie('session_id', sessionId);
  }

  // --- last action ts (ms since epoch)
  const LA_KEY = 'last_action_ts';
  function readLastAction() {
    const v = localStorage.getItem(LA_KEY);
    const n = v ? parseInt(v, 10) : 0;
    return Number.isFinite(n) ? n : 0;
  }
  function bumpLastAction() {
    localStorage.setItem(LA_KEY, String(nowMs()));
  }
  // на загрузке страницы считаем это действием
  bumpLastAction();

  // --- session TTL: 10 минут
  const SESSION_TTL_MS = 10 * 60 * 1000;
  // проверка раз в минуту
  setInterval(() => {
    const last = readLastAction();
    if (nowMs() - last > SESSION_TTL_MS) {
      // ре-инициализируем session_id
      sessionId = generateId('session');
      setCookie('session_id', sessionId);
      bumpLastAction(); // начать новую сессию с текущего момента
      console.log('[tracker] session rotated due to inactivity >10m');
    }
  }, 60 * 1000);

  // лёгкое обновление last_action_ts на базовые юзер-сигналы
  const bump = () => bumpLastAction();
  ['click','keydown','scroll','touchstart','visibilitychange'].forEach(evt => {
    window.addEventListener(evt, bump, { passive: true });
  });

  // --- send visit
  fetch('https://crm.retry.school/track/visit', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      visit_id:   visitId,
      session_id: sessionId,
      source:     source || getSource(),
      timestamp:  nowIso(),
    })
  }).then(res => {
    console.log('[tracker] /track/visit status:', res.status);
  }).catch(err => console.error('[tracker] visit error:', err));
})();