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
    console.log(`[tracker] cookie set: ${name}=${value} (expires: ${expires})`);
  }
  function generateId(prefix) {
    const id = `${prefix}_${Math.random().toString(36).slice(2)}_${Date.now()}`;
    console.log(`[tracker] generated ${prefix}_id: ${id}`);
    return id;
  }
  function nowIso() { return new Date().toISOString(); }
  function nowMs() { return Date.now(); }

  // --- source
  function getSource() {
    const params = new URLSearchParams(window.location.search);
    if (params.has('utm_source')) {
      const src = 'utm:' + params.get('utm_source');
      console.log('[tracker] source detected from utm:', src);
      return src;
    }
    if (document.referrer) {
      try {
        const src = 'ref:' + new URL(document.referrer).hostname;
        console.log('[tracker] source detected from referrer:', src);
        return src;
      } catch {
        console.warn('[tracker] failed to parse referrer:', document.referrer);
      }
    }
    console.log('[tracker] source detected as direct');
    return 'direct';
  }

  // --- visit bootstrap
  let visitId = getCookie('visit_id');
  let source  = getCookie('visit_source');
  let isNewVisit = false;

  if (!visitId) {
    console.log('[tracker] visit_id not found, generating new');
    visitId = generateId('visit');
    source = getSource();
    setCookie('visit_id', visitId);
    setCookie('visit_source', source);
    isNewVisit = true;
  } else {
    console.log(`[tracker] existing visit_id found: ${visitId}, source: ${source}`);
  }

  // --- session bootstrap
  let sessionId = getCookie('session_id');
  if (!sessionId) {
    console.log('[tracker] session_id not found, generating new');
    sessionId = generateId('session');
    setCookie('session_id', sessionId);
  } else {
    console.log(`[tracker] existing session_id found: ${sessionId}`);
  }

  // --- last action ts
  const LA_KEY = 'last_action_ts';
  function readLastAction() {
    const v = localStorage.getItem(LA_KEY);
    const n = v ? parseInt(v, 10) : 0;
    return Number.isFinite(n) ? n : 0;
  }
  function bumpLastAction(reason = '') {
    localStorage.setItem(LA_KEY, String(nowMs()));
    console.log(`[tracker] last_action_ts updated to ${nowIso()} ${reason ? `(${reason})` : ''}`);
  }

  // --- throttled bump (min 1.5s between actions)
  let lastBumpTime = 0;
  const MIN_BUMP_INTERVAL = 1500; // ms
  function throttledBump(reason) {
    const now = nowMs();
    if (now - lastBumpTime >= MIN_BUMP_INTERVAL) {
      bumpLastAction(reason);
      lastBumpTime = now;
    }
  }

  // на загрузке страницы считаем это действием
  throttledBump('page load');

  // --- session TTL: 10 минут
  const SESSION_TTL_MS = 10 * 60 * 1000;

  // проверка раз в минуту
  setInterval(() => {
    const last = readLastAction();
    const diff = nowMs() - last;
    const idleSeconds = Math.round(diff / 1000);

    if (diff > SESSION_TTL_MS) {
      const oldSession = sessionId;
      sessionId = generateId('session');
      setCookie('session_id', sessionId);
      throttledBump('session rotated');
      console.log(`[tracker] session rotated due to inactivity >10m (old: ${oldSession}, new: ${sessionId}, idle: ${idleSeconds}s)`);
    } else {
      console.log(`[tracker] session still active, idle=${idleSeconds}s, session_id=${sessionId}`);
    }
  }, 60 * 1000);

  // --- attach events with throttling
  ['click','keydown','scroll','touchstart','visibilitychange'].forEach(evt => {
    window.addEventListener(evt, (e) => throttledBump(`event: ${e.type}`), { passive: true });
  });

  // --- send visit only if new
  if (isNewVisit) {
    console.log(`[tracker] sending /track/visit: visit_id=${visitId}, session_id=${sessionId}, source=${source || getSource()}, ts=${nowIso()}`);
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
      console.log('[tracker] /track/visit response status:', res.status);
    }).catch(err => console.error('[tracker] visit error:', err));
  } else {
    console.log('[tracker] /track/visit skipped — existing visit_id in cookies');
  }
})();