(function() {
  function getCookie(name) {
    const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return match ? decodeURIComponent(match[2]) : null;
  }

  function setCookie(name, value, days = 365) {
    const expires = new Date(Date.now() + days*864e5).toUTCString();
    document.cookie = ${name}=${encodeURIComponent(value)}; path=/; expires=${expires};
  }

  function generateVisitId() {
    return 'visit_' + Math.random().toString(36).substring(2) + Date.now();
  }

  function getSource() {
    const params = new URLSearchParams(window.location.search);
    if (params.has('utm_source')) return 'utm:' + params.get('utm_source');
    if (document.referrer) return 'ref:' + new URL(document.referrer).hostname;
    return 'direct';
  }

  // --- логика ---
  if (!getCookie('visit_id')) {
    setCookie('visit_id', generateVisitId());
    setCookie('visit_source', getSource());
  }

  // // опционально — шлём в бек
  // fetch('https://your-backend.com/api/track', {
  //   method: 'POST',
  //   headers: {'Content-Type': 'application/json'},
  //   body: JSON.stringify({
  //     visit_id: getCookie('visit_id'),
  //     source: getCookie('visit_source'),
  //     timestamp: new Date().toISOString(),
  //   })
  // });

})();
