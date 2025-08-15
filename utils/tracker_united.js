<script>
(function () {
  // Настройка
  const CTA_SELECTOR = '[data-track="click_cta_top"]'; // верхняя кнопка
  const SCROLL_THROTTLE_MS = 1500; // как часто слать scroll_depth

  // Ждём айдишники из первого скрипта, но не дольше 2с
  function waitIds(timeoutMs=2000){
    return new Promise((resolve)=>{
      if (window.visitId && window.sessionId) return resolve(true);
      const t0 = Date.now();
      const t = setInterval(()=>{
        if (window.visitId && window.sessionId){ clearInterval(t); resolve(true); }
        else if (Date.now()-t0 > timeoutMs){ clearInterval(t); resolve(false); }
      }, 50);
    });
  }

  // Отправка
  async function post(type, meta){
    await waitIds();
    const payload = {
      visit_id:   window.visitId || null,
      session_id: window.sessionId || null,
      type,
      source:     window.visitSource || null,
      timestamp:  new Date().toISOString(),
      ...(meta ? { meta } : {})
    };
    fetch("https://crm.retry.school/track/action", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
      keepalive: true
    }).catch(()=>{});
  }

  // Клик по верхней CTA
  document.addEventListener("click", (e) => {
    const el = e.target.closest(CTA_SELECTOR);
    if (el) post("click_cta_top");
  }, { capture: true, passive: true });

  // Scroll depth с троттлингом
  function getDepthPct(){
    const se = document.scrollingElement || document.documentElement;
    const scrollTop = se.scrollTop || 0;
    const docH = Math.max(se.scrollHeight, document.documentElement.scrollHeight, document.body.scrollHeight);
    const winH = window.innerHeight || document.documentElement.clientHeight || 0;
    const maxScroll = Math.max(docH - winH, 1);
    return Math.min(100, Math.round((scrollTop / maxScroll) * 100));
  }

  let lastSentAt = 0;
  let pending = false;
  function onScrollThrottled(){
    const now = Date.now();
    if (now - lastSentAt >= SCROLL_THROTTLE_MS){
      lastSentAt = now;
      post("scroll_depth", { depth_pct: getDepthPct() });
    } else if (!pending) {
      pending = true;
      const delay = SCROLL_THROTTLE_MS - (now - lastSentAt);
      setTimeout(() => {
        lastSentAt = Date.now();
        pending = false;
        post("scroll_depth", { depth_pct: getDepthPct() });
      }, Math.max(0, delay));
    }
  }

  // Подписки на скролл/жесты
  addEventListener("scroll", onScrollThrottled, { passive: true });
  addEventListener("wheel", onScrollThrottled, { passive: true });
  addEventListener("touchmove", onScrollThrottled, { passive: true });

  // Стартовая отправка (на случай короткой страницы)
  onScrollThrottled();

  // Финальный «хвост» при уходе со страницы
  addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") {
      try {
        const payload = JSON.stringify({
          visit_id: window.visitId || null,
          session_id: window.sessionId || null,
          type: "scroll_depth",
          source: window.visitSource || null,
          timestamp: new Date().toISOString(),
          meta: { depth_pct: getDepthPct(), final: true }
        });
        navigator.sendBeacon("https://crm.retry.school/track/action", new Blob([payload], {type: "application/json"}));
      } catch {}
    }
  });
})();
</script>