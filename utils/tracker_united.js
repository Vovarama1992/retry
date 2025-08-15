<script>
(function () {
  const LOG = (...a) => console.log("[track]", ...a);
  const ERR = (...a) => console.error("[track]", ...a);

  // --- кнопки (как у тебя)
  const BTN_TEXTS = ["получить доступ"];              // сравнивается в нижнем регистре
  const HREF_HINTS = ["#rec573904816"];               // фрагменты href
  const BASE_SELECTOR = "a.tn-atom, button.tn-atom";  // типовые Tilda-кнопки

  function normText(el) {
    return (el.innerText || el.textContent || "").trim().toLowerCase();
  }
  function isCandidate(el) {
    if (!(el instanceof Element)) return false;
    if (!el.matches(BASE_SELECTOR)) return false;
    const txt = normText(el);
    const href = el.getAttribute("href") || "";
    const textMatch = txt && BTN_TEXTS.some(t => txt.includes(t));
    const hrefMatch = href && HREF_HINTS.some(h => href.includes(h));
    return textMatch || hrefMatch || el.dataset.track === "click_cta_top";
  }
  function markTracked(el) {
    if (el.dataset.track !== "click_cta_top") {
      el.setAttribute("data-track", "click_cta_top");
      LOG("marked button with data-track=click_cta_top", el);
    }
  }

  // --- отправка
  function post(type, extra) {
    const payload = {
      visit_id:   window.visitId || null,
      session_id: window.sessionId || null,
      type,
      source:     window.visitSource || null,
      timestamp:  new Date().toISOString(),
      ...(extra ? { meta: extra } : {})
    };
    LOG("POST /track/action", payload);
    fetch("https://crm.retry.school/track/action", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
      keepalive: true
    })
    .then(r => { LOG("status:", r.status); return r.text(); })
    .then(b => LOG("body:", b))
    .catch(ERR);
  }

  // --- клики по CTA
  document.addEventListener("click", (e) => {
    let n = e.target instanceof Element ? e.target : null;
    while (n && n !== document.body) {
      if (isCandidate(n)) {
        LOG("delegated click on candidate", n);
        markTracked(n);
        post("click_cta_top");
        break;
      }
      n = n.parentElement;
    }
  }, { capture: true, passive: true });

  // первичная маркировка
  document.querySelectorAll(BASE_SELECTOR).forEach(el => {
    if (isCandidate(el)) markTracked(el);
  });

  // --- SCROLL DEPTH
  const THRESHOLDS = [25, 50, 75, 100];
  const fired = new Set();
  function getDepthPct(){
    const y = window.pageYOffset || document.documentElement.scrollTop || 0;
    const docH = Math.max(
      document.body.scrollHeight, document.documentElement.scrollHeight,
      document.body.offsetHeight,  document.documentElement.offsetHeight,
      document.body.clientHeight,  document.documentElement.clientHeight
    );
    const winH = window.innerHeight || document.documentElement.clientHeight || 0;
    const maxScroll = Math.max(docH - winH, 1);
    return Math.min(100, Math.round((y / maxScroll) * 100));
  }
  let ticking = false;
  function onScroll(){
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(() => {
      const pct = getDepthPct();
      for (const t of THRESHOLDS){
        if (pct >= t && !fired.has(t)){
          fired.add(t);
          LOG("scroll_depth threshold reached:", t);
          // если бэку не нужен meta — убери второй аргумент
          post("scroll_depth", { depth_pct: t });
        }
      }
      ticking = false;
    });
  }
  // стартовая проверка и подписка
  onScroll();
  document.addEventListener("scroll", onScroll, { passive: true });

  LOG("boot (delegation + scroll_depth). Waiting for:", { BTN_TEXTS, HREF_HINTS, THRESHOLDS });
})();
</script>