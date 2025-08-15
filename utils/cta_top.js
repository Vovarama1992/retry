<script>
(function () {
  const LOG = (...a) => console.log("[track]", ...a);
  const ERR = (...a) => console.error("[track]", ...a);

  // Под свои реальные варианты текста/якоря при необходимости добавь ещё
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

  function sendAction() {
    try {
      const payload = {
        visit_id:   window.visitId || null,
        session_id: window.sessionId || null,
        type:       "click_cta_top",
        source:     window.visitSource || null,
        timestamp:  new Date().toISOString()
      };
      LOG("POST /track/action", payload);
      fetch("https://crm.retry.school/track/action", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      })
      .then(r => { LOG("status:", r.status); return r.text(); })
      .then(b => LOG("body:", b))
      .catch(ERR);
    } catch (e) {
      ERR("exception while sending:", e);
    }
  }

  // Делегирование клика (ловит и текущие, и будущие элементы)
  document.addEventListener("click", (e) => {
    let n = e.target instanceof Element ? e.target : null;
    while (n && n !== document.body) {
      if (isCandidate(n)) {
        LOG("delegated click on candidate", n);
        markTracked(n);
        sendAction();
        break;
      }
      n = n.parentElement;
    }
  }, { capture: true, passive: true });

  // Первичный скан (на случай, если кнопка уже есть)
  document.querySelectorAll(BASE_SELECTOR).forEach(el => {
    if (isCandidate(el)) markTracked(el);
  });

  // Observer — помечаем, когда Tilda дорисует кнопку
  const observer = new MutationObserver(muts => {
    for (const m of muts) {
      for (const node of m.addedNodes) {
        if (node.nodeType !== 1) continue;
        if (isCandidate(node)) markTracked(node);
        node.querySelectorAll?.(BASE_SELECTOR).forEach(el => {
          if (isCandidate(el)) markTracked(el);
        });
      }
    }
  });
  observer.observe(document.documentElement, { childList: true, subtree: true });

  LOG("boot (delegation+observer). Waiting for:", { BTN_TEXTS, HREF_HINTS });
})();
</script>