package roistat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/domain"
)

// RoistatClient – клиент для отправки событий в Ройстат
type RoistatClient struct {
	apiURL string
	apiKey string
	http   *http.Client
	logger logger.Logger
}

// NewRoistatClient создаёт клиента, подтягивая ключ и URL из ENV
func NewRoistatClient(l logger.Logger) *RoistatClient {
	apiKey := os.Getenv("ROISTAT_KEY")
	apiURL := os.Getenv("ROISTAT_URL")
	if apiURL == "" {
		apiURL = "https://cloud.roistat.com/api/proxy/1.0/leads"
		l.Log(logger.LogEntry{
			Level:   "warn",
			Service: "track",
			Method:  "NewRoistatClient",
			Message: "ROISTAT_URL не задан, используем дефолтный https://cloud.roistat.com/api/proxy/1.0/leads",
		})
	}
	if apiKey == "" {
		l.Log(logger.LogEntry{
			Level:   "warn",
			Service: "track",
			Method:  "NewRoistatClient",
			Message: "ROISTAT_KEY не задан в ENV",
		})
	}

	return &RoistatClient{
		apiURL: apiURL,
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
		logger: l,
	}
}

func (c *RoistatClient) SendProceedToPayment(ctx context.Context, action domain.Action) error {
	visit := extractFromMeta(action.Meta, "roistat_visit")
	email := extractFromMeta(action.Meta, "email")
	social := extractFromMeta(action.Meta, "social_link")
	method := extractFromMeta(action.Meta, "name")
	page := extractFromMeta(action.Meta, "page")

	if visit == "" {
		c.logger.Log(logger.LogEntry{
			Level:   "warn",
			Service: "track",
			Method:  "SendProceedToPayment",
			Message: "[Roistat] пропускаем отправку: roistat_visit пуст",
		})
		return nil
	}

	// Собираем URL с query: key и roistat_visit
	u, err := url.Parse(c.apiURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("roistat_visit", visit)
	u.RawQuery = q.Encode()

	// Тело — как у proxy lead: title/email/fields
	bodyObj := map[string]any{
		"title": "Перейти к оплате",
		"email": email,
		"fields": map[string]any{
			"social_link":    social,
			"payment_method": method,
			"page":           page,
		},
	}
	body, _ := json.Marshal(bodyObj)

	// Логируем без утечки ключа
	redacted := *u
	rq := redacted.Query()
	rq.Set("key", "***")
	redacted.RawQuery = rq.Encode()

	c.logger.Log(logger.LogEntry{
		Level:   "info",
		Service: "track",
		Method:  "SendProceedToPayment",
		Message: fmt.Sprintf("[Roistat] 🚀 POST %s body=%s", redacted.String(), string(body)),
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		c.logger.Log(logger.LogEntry{Level: "error", Service: "track", Method: "SendProceedToPayment", Message: "failed to create request", Error: err})
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Log(logger.LogEntry{Level: "error", Service: "track", Method: "SendProceedToPayment", Message: "http request failed", Error: err})
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.logger.Log(logger.LogEntry{
		Level:   "info",
		Service: "track",
		Method:  "SendProceedToPayment",
		Message: fmt.Sprintf("[Roistat] ✅ статус %d, ответ: %s", resp.StatusCode, string(respBody)),
	})

	if resp.StatusCode >= 300 {
		return fmt.Errorf("roistat error: %s", resp.Status)
	}
	return nil
}

// достаём значение по ключу из action.Meta (json.RawMessage)
// поддерживаем как плоский вид, так и вложенный объект "meta"
func extractFromMeta(meta json.RawMessage, key string) string {
	if len(meta) == 0 {
		return ""
	}

	var m map[string]any
	if err := json.Unmarshal(meta, &m); err != nil {
		return ""
	}

	// 1. Пробуем найти на верхнем уровне
	if v, ok := m[key].(string); ok {
		return v
	}

	// 2. Если есть вложенный meta — копаем туда
	if inner, ok := m["meta"].(map[string]any); ok {
		if v, ok := inner[key].(string); ok {
			return v
		}
	}

	return ""
}
