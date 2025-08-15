document.addEventListener("DOMContentLoaded", function () {
  var btn = document.querySelector('[data-track="click_cta_top"]');
  if (!btn) return;

  btn.addEventListener("click", function () {
    console.log("[track] Клик по кнопке 'click_cta_top'");

    try {
      console.log("[track] Отправка запроса на /track/action...");
      fetch("/track/action", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          visitId: window.visitId || null,
          sessionId: window.sessionId || null,
          action: "click_cta_top",
          ts: Date.now()
        })
      })
        .then(res => {
          console.log("[track] Ответ получен, статус:", res.status);
          return res.text();
        })
        .then(body => {
          console.log("[track] Тело ответа:", body);
        })
        .catch(err => {
          console.error("[track] Ошибка при запросе:", err);
        });
    } catch (e) {
      console.error("[track] Исключение в обработчике клика:", e);
    }
  });
});
