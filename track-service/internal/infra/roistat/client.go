package roistat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Vovarama1992/go-utils/logger"
	"github.com/Vovarama1992/retry/pkg/domain"
)

type RoistatClient struct {
	apiBase   string
	apiKey    string
	projectID string
	http      *http.Client
	logger    logger.Logger
}

func NewRoistatClient(l logger.Logger) *RoistatClient {
	apiKey := os.Getenv("ROISTAT_KEY")
	projectID := os.Getenv("ROISTAT_PROJECT_ID")
	apiBase := os.Getenv("ROISTAT_URL")
	if apiBase == "" {
		apiBase = "https://cloud.roistat.com/api/v1/project"
	}

	if apiKey == "" {
		l.Log(logger.LogEntry{
			Level:   "warn",
			Service: "track",
			Method:  "NewRoistatClient",
			Message: "ROISTAT_KEY не задан в ENV",
		})
	}
	if projectID == "" {
		l.Log(logger.LogEntry{
			Level:   "warn",
			Service: "track",
			Method:  "NewRoistatClient",
			Message: "ROISTAT_PROJECT_ID не задан в ENV",
		})
	}

	return &RoistatClient{
		apiBase:   apiBase,
		apiKey:    apiKey,
		projectID: projectID,
		http:      &http.Client{Timeout: 10 * time.Second},
		logger:    l,
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

	url := fmt.Sprintf("%s/%s/integration/order/add", c.apiBase, c.projectID)

	bodyObj := map[string]any{
		"roistat": visit,
		"email":   email,
		"title":   "Перейти к оплате",
		"comment": fmt.Sprintf("Метод оплаты: %s, страница: %s", method, page),
		"fields": map[string]any{
			"social_link":    social,
			"payment_method": method,
			"page":           page,
		},
	}
	body, _ := json.Marshal(bodyObj)

	c.logger.Log(logger.LogEntry{
		Level:   "info",
		Service: "track",
		Method:  "SendProceedToPayment",
		Message: fmt.Sprintf("[Roistat] 🚀 POST %s body=%s", url, string(body)),
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
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
