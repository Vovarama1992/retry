package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type Meta map[string]any

var ruMonths = [...]string{
	"", "января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func ruTime(ts time.Time, loc *time.Location) string {
	if loc == nil {
		loc, _ = time.LoadLocation("Europe/Amsterdam")
	}
	t := ts.In(loc)
	return fmt.Sprintf("%d %s %02d:%02d", t.Day(), ruMonths[int(t.Month())], t.Hour(), t.Minute())
}

func str(m Meta, k string) string {
	if v, ok := m[k]; ok {
		switch t := v.(type) {
		case string:
			return t
		case float64:
			if t == float64(int64(t)) {
				return fmt.Sprintf("%d", int64(t))
			}
			return fmt.Sprintf("%.0f", t)
		}
	}
	return ""
}

func HumanActionLine(ts time.Time, actionType string, metaRaw []byte, loc *time.Location) string {
	var m Meta
	_ = json.Unmarshal(metaRaw, &m)

	t := ruTime(ts, loc)

	switch actionType {
	case "click_cta_top":
		label := str(m, "label")
		if label == "" {
			label = "Получить доступ"
		}
		return fmt.Sprintf("%s — нажал кнопку «%s» (верхняя)", t, label)

	case "click_cta_bottom":
		label := str(m, "label")
		if label == "" {
			label = "Получить доступ"
		}
		return fmt.Sprintf("%s — нажал кнопку «%s» (нижняя)", t, label)

	case "scroll_depth":
		p := str(m, "percent")
		if p != "" && p[len(p)-1] != '%' {
			p += "%"
		}
		if p == "" {
			p = "—"
		}
		return fmt.Sprintf("%s — пролистал до %s", t, p)

	case "scroll_section_view":
		sec := str(m, "section")
		if sec == "" {
			sec = str(m, "id")
		}
		if sec == "" {
			sec = str(m, "path")
		}
		if sec == "" {
			return fmt.Sprintf("%s — просмотр секции", t)
		}
		return fmt.Sprintf("%s — просмотр секции «%s»", t, sec)

	case "gallery_scroll_right":
		return fmt.Sprintf("%s — пролистал галерею учителей вправо", t)

	case "gallery_scroll_left":
		return fmt.Sprintf("%s — пролистал галерею учителей влево", t)

	case "gallery_raids_scroll":
		idx := str(m, "index")
		if idx == "" {
			return fmt.Sprintf("%s — прокрутка галереи рейдов", t)
		}
		return fmt.Sprintf("%s — прокрутка галереи рейдов (слайд %s)", t, idx)

	case "external_link_raid":
		lbl := str(m, "label")
		if lbl == "" {
			lbl = "Raid"
		}
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)

	case "external_link_raids":
		lbl := str(m, "label")
		if lbl == "" {
			lbl = "Наши рейды"
		}
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)

	case "external_link_details":
		lbl := str(m, "label")
		if lbl == "" {
			lbl = "Подробнее"
		}
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — перешёл по ссылке «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл по ссылке «%s»", t, lbl)

	case "external_link_mentor_page":
		lbl := str(m, "label")
		if lbl == "" {
			lbl = "страница ментора"
		}
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — перешёл на «%s» (%s)", t, lbl, url)
		}
		return fmt.Sprintf("%s — перешёл на «%s»", t, lbl)

	case "external_link_social":
		label := str(m, "label")
		if label == "" {
			label = "соцссылка"
		}
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — перешёл по внешней ссылке «%s» (%s)", t, label, url)
		}
		return fmt.Sprintf("%s — перешёл по внешней ссылке «%s»", t, label)

	case "faq_open_answer":
		q := str(m, "question")
		if q == "" {
			return fmt.Sprintf("%s — открыл ответ в FAQ", t)
		}
		return fmt.Sprintf("%s — открыл ответ в FAQ: «%s»", t, q)

	case "click_links_buy_access":
		return fmt.Sprintf("%s — клик по ссылке «Купить доступ»", t)

	case "click_links_telegram":
		return fmt.Sprintf("%s — клик по ссылке Telegram из https://retry.school/links", t)

	case "click_links_youtube_entertainment":
		return fmt.Sprintf("%s — клик по ссылке YouTube (развлекательные видео) из https://retry.school/links", t)

	case "click_links_youtube_streams":
		return fmt.Sprintf("%s — клик по ссылке YouTube (стримы) из https://retry.school/links", t)

	case "click_links_instagram":
		return fmt.Sprintf("%s — клик по ссылке Instagram из https://retry.school/links", t)

	case "click_links_tiktok":
		return fmt.Sprintf("%s — клик по ссылке TikTok из https://retry.school/links", t)

	case "click_links_artstation":
		return fmt.Sprintf("%s — клик по ссылке ArtStation из https://retry.school/links", t)

	case "click_links_3d_guide":
		return fmt.Sprintf("%s — клик по ссылке «3D-гайд» из https://retry.school/links", t)

	default:
		url := str(m, "url")
		if url != "" {
			return fmt.Sprintf("%s — %s (%s)", t, actionType, url)
		}
		return fmt.Sprintf("%s — %s", t, actionType)
	}
}
