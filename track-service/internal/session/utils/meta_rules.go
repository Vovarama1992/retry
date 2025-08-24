package utils

import (
	"fmt"
)

type Rule struct {
	Fields   []string // какие поля читаем из meta
	FormatFn func(m map[string]any) string
}

var EnrichRules = map[string]Rule{
	"external_link_raid": {
		Fields: []string{"text", "raid_name"},
		FormatFn: func(m map[string]any) string {
			name, _ := m["text"].(string)
			if name == "" {
				name, _ = m["raid_name"].(string)
			}
			if name == "" {
				name = "рейд"
			}
			return fmt.Sprintf("перешёл на рейд «%s»", name)
		},
	},

	"external_link_mentor_page": {
		Fields: []string{"mentor_name"},
		FormatFn: func(m map[string]any) string {
			name, _ := m["mentor_name"].(string)
			if name == "" {
				name = "страница ментора"
			}
			return fmt.Sprintf("перешёл на страницу ментора «%s»", name)
		},
	},

	"external_link_social": {
		Fields: []string{"platform", "text"},
		FormatFn: func(m map[string]any) string {
			platform, _ := m["platform"].(string)
			if platform == "" {
				platform = "соцсеть"
			}
			text, _ := m["text"].(string)
			if text != "" && text != platform {
				return fmt.Sprintf("перешёл в %s (%s)", platform, text)
			}
			return fmt.Sprintf("перешёл в %s", platform)
		},
	},

	"faq_open_answer": {
		Fields: []string{"question_text"},
		FormatFn: func(m map[string]any) string {
			q, _ := m["question_text"].(string)
			if q == "" {
				return "открыл ответ в FAQ"
			}
			return fmt.Sprintf("открыл ответ в FAQ: «%s»", q)
		},
	},
}
