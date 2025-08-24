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

func HumanActionLine(ts time.Time, actionType string, metaRaw []byte, loc *time.Location) string {
	var m map[string]any
	_ = json.Unmarshal(metaRaw, &m)

	t := ruTime(ts, loc)

	// 1. enrichment
	var details string
	if enrichRule, ok := EnrichRules[actionType]; ok {
		details = enrichRule.FormatFn(m)
	}

	// 2. humanization
	if humanRule, ok := HumanRules[actionType]; ok {
		// пробрасываем details внутрь
		return humanRule(t, map[string]any{
			"details": details,
			"raw":     m,
		})
	}

	// 3. fallback
	url, _ := m["url"].(string)
	if url != "" {
		return fmt.Sprintf("%s — %s (%s)", t, actionType, url)
	}
	return fmt.Sprintf("%s — %s", t, actionType)
}
