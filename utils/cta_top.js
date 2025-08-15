document.addEventListener('DOMContentLoaded', function () {
  if (window.__ctaTopHook) return; window.__ctaTopHook = true;

  function lc(s){ return (s||'').trim().toLowerCase(); }
  function bumpLater(){ setTimeout(()=>{ try{ localStorage.setItem('last_action_ts', String(Date.now())); }catch(e){} }, 1000); }

  const candidates = Array.from(document.querySelectorAll('a.tn-atom, a.t-btn, .t-btn a, a[href]'));
  const btn = candidates.find(el => lc(el.textContent) === 'получить доступ');

  if (!btn) { console.warn('[tracker] CTA TOP not found'); return; }

  btn.addEventListener('click', function () {
    bumpLater();
    const rec = this.closest('[id^="rec"]')?.id || null;
    if (typeof window.trackAction === 'function') {
      window.trackAction('click_cta_top', { text: 'ПОЛУЧИТЬ ДОСТУП', rec });
    } else {
      console.warn('[tracker] trackAction not ready');
    }
  });
});
