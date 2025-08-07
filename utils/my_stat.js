(function () {
  console.log('[tracker] Скрипт запущен');

  function getCookie(name) {
    const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return match ? decodeURIComponent(match[2]) : null;
  }

  function setCookie(name, value, days = 365) {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; path=/; expires=${expires}`;
    console.log(`[tracker] Установлена кука ${name}=${value}`);
  }

  function generateVisitId() {
    const id = 'visit_' + Math.random().toString(36).substring(2) + Date.now();
    console.log('[tracker] Сгенерирован visit_id:', id);
    return id;
  }

  function getSource() {
    const params = new URLSearchParams(window.location.search);
    if (params.has('utm_source')) {
      const utm = 'utm:' + params.get('utm_source');
      console.log('[tracker] Источник по utm:', utm);
      return utm;
    }
    if (document.referrer) {
      const ref = 'ref:' + new URL(document.referrer).hostname;
      console.log('[tracker] Источник по referrer:', ref);
      return ref;
    }
    console.log('[tracker] Источник: direct');
    return 'direct';
  }

  const visit_id = getCookie('visit_id');
  const visit_source = getCookie('visit_source');

  if (!visit_id) {
    console.log('[tracker] Куки отсутствуют, создаём заново...');
    const newVisitId = generateVisitId();
    const newSource = getSource();
    setCookie('visit_id', newVisitId);
    setCookie('visit_source', newSource);
  } else {
    console.log('[tracker] Куки уже есть:');
    console.log('  visit_id:', visit_id);
    console.log('  visit_source:', visit_source);
  }

  const finalVisitId = getCookie('visit_id');
  const finalSource = getCookie('visit_source');

  console.log('[tracker] Отправляем запрос на backend...');
  fetch('https://crm.retry.school/track/visit', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      visit_id: finalVisitId,
      source: finalSource,
      timestamp: new Date().toISOString(),
    })
  })
  .then(res => console.log('[tracker] Ответ от backend:', res.status))
  .catch(err => console.error('[tracker] Ошибка при отправке:', err));
})();
