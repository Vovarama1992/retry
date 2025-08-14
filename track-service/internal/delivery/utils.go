package actionhttp

import (
	"net"
	"net/http"

	"net/url"
	"strings"
)

func ExtractIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

func NormalizeSource(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "Неизвестно"
	}

	// Прямой заход
	if raw == "direct" {
		return "Прямой заход"
	}

	// UTM-метки
	if strings.HasPrefix(raw, "utm:") {
		utmValue := strings.TrimPrefix(raw, "utm:")
		switch utmValue {
		case "tg":
			return "Telegram"
		case "youtube":
			return "YouTube"
		case "vk":
			return "ВКонтакте"
		case "social":
			return "Соцсети"
		case "video":
			return "Видео"
		default:
			return utmValue
		}
	}

	// Рефералы (ref:)
	if strings.HasPrefix(raw, "ref:") {
		refValue := strings.TrimPrefix(raw, "ref:")
		u, err := url.Parse(refValue)
		if err == nil && u.Host != "" {
			refValue = u.Host
		}
		// убираем www. и m.
		refValue = strings.TrimPrefix(refValue, "www.")
		refValue = strings.TrimPrefix(refValue, "m.")
		switch refValue {
		case "google.com", "google.ru":
			return "Google"
		case "yandex.ru", "ya.ru":
			return "Яндекс"
		case "vk.com":
			return "ВКонтакте"
		case "youtube.com":
			return "YouTube"
		case "facebook.com":
			return "Facebook"
		case "web.telegram.org", "org.telegram.messenger", "t.me":
			return "Telegram"
		default:
			return refValue
		}
	}

	// Всё остальное как есть
	return raw
}
