package utils

import "fmt"

// HumanRule принимает timestamp как строку (t) и уже обогащённые детали
type HumanRule func(t string, m map[string]any) string

var HumanRules = map[string]HumanRule{
	"click_cta_top": func(t string, m map[string]any) string {
		label, _ := m["label"].(string)
		if label == "" {
			label = "Получить доступ"
		}
		return fmt.Sprintf("%s — нажал кнопку «%s» (верхняя)", t, label)
	},

	"click_cta_bottom": func(t string, m map[string]any) string {
		label, _ := m["label"].(string)
		if label == "" {
			label = "Получить доступ"
		}
		return fmt.Sprintf("%s — нажал кнопку «%s» (нижняя)", t, label)
	},

	"scroll_depth": func(t string, m map[string]any) string {
		sec, _ := m["section_title"].(string)
		if sec == "" {
			sec, _ = m["section_id"].(string)
		}
		if sec != "" {
			return fmt.Sprintf("%s — остановился на секции «%s»", t, sec)
		}
		return fmt.Sprintf("%s — пролистал страницу", t)
	},

	"scroll_section_view": func(t string, m map[string]any) string {
		sec, _ := m["section_title"].(string)
		if sec == "" {
			sec, _ = m["section_id"].(string)
		}
		if sec != "" {
			return fmt.Sprintf("%s — просмотр секции «%s»", t, sec)
		}
		return fmt.Sprintf("%s — просмотр секции", t)
	},

	"gallery_scroll_right": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — пролистал галерею учителей вправо", t)
	},

	"gallery_scroll_left": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — пролистал галерею учителей влево", t)
	},

	"gallery_raids_scroll": func(t string, m map[string]any) string {
		idx, _ := m["index"].(string)
		if idx == "" {
			return fmt.Sprintf("%s — прокрутка галереи рейдов", t)
		}
		return fmt.Sprintf("%s — прокрутка галереи рейдов (слайд %s)", t, idx)
	},

	"external_link_raid": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = "Raid"
		}
		url, _ := m["url"].(string)
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)
	},

	"external_link_raids": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = "Наши рейды"
		}
		url, _ := m["url"].(string)
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)
	},

	"external_link_details": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = "Подробнее"
		}
		url, _ := m["url"].(string)
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)
	},

	"external_link_mentor_page": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = "страница ментора"
		}
		url, _ := m["url"].(string)
		if url != "" {
			return fmt.Sprintf("%s — перешёл на «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл на «%s»", t, lbl)
	},

	"external_link_social": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		lbl, _ := m["label"].(string)
		if lbl == "" {
			lbl = "соцссылка"
		}
		url, _ := m["url"].(string)
		if url != "" {
			return fmt.Sprintf("%s — перешёл по внешней ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по внешней ссылке «%s»", t, lbl)
	},

	"faq_open_answer": func(t string, m map[string]any) string {
		if det, _ := m["details"].(string); det != "" {
			return fmt.Sprintf("%s — %s", t, det)
		}
		q, _ := m["question"].(string)
		if q == "" {
			return fmt.Sprintf("%s — открыл ответ в FAQ", t)
		}
		return fmt.Sprintf("%s — открыл ответ в FAQ: «%s»", t, q)
	},

	"click_links_buy_access": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке «Купить доступ»", t)
	},

	"click_links_telegram": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке Telegram из https://retry.school/links", t)
	},

	"click_links_youtube_entertainment": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке YouTube (развлекательные видео) из https://retry.school/links", t)
	},

	"click_links_youtube_streams": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке YouTube (стримы) из https://retry.school/links", t)
	},

	"click_links_instagram": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке Instagram из https://retry.school/links", t)
	},

	"click_links_tiktok": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке TikTok из https://retry.school/links", t)
	},

	"click_links_artstation": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке ArtStation из https://retry.school/links", t)
	},

	"click_links_3d_guide": func(t string, _ map[string]any) string {
		return fmt.Sprintf("%s — клик по ссылке «3D-гайд» из https://retry.school/links", t)
	},
}
