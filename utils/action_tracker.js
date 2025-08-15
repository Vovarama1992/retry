(function () {
  function getCookie(name) {
    const m = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return m ? decodeURIComponent(m[2]) : null;
  }
  function nowIso() { return new Date().toISOString(); }

  const LA_KEY = 'last_action_ts';
  function bumpLastAction() {
    setTimeout(() => {
      localStorage.setItem(LA_KEY, String(Date.now()));
    }, 1000); // задержка 1 секунда
  }

  window.trackAction = function (actionType, metaObj) {
    const visitId   = getCookie('visit_id');
    const sessionId = getCookie('session_id');

    if (!visitId || !sessionId) {
      console.warn('[tracker] нет visit_id или session_id — действие не отправлено');
      return;
    }

    bumpLastAction();

    fetch('https://crm.retry.school/track/action', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action_type: actionType,
        visit_id:    visitId,
        session_id:  sessionId,
        timestamp:   nowIso(),
        meta:        metaObj || {}
      })
    }).then(res => {
      console.log('[tracker] action', actionType, 'status:', res.status);
    }).catch(err => console.error('[tracker] action error:', err));
  };
})();