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

// SendProceedToPayment отправляет событие "Перейти к оплате" в Ройстат
func (c *RoistatClient) SendProceedToPayment(ctx context.Context, action domain.Action) error {
	payload := map[string]any{
		"roistat_visit": extractFromMeta(action.Meta, "roistat_visit"),
		"email":         extractFromMeta(action.Meta, "email"),
		"social_link":   extractFromMeta(action.Meta, "social_link"),
		"payment":       extractFromMeta(action.Meta, "name"),
	}

	body, _ := json.Marshal(payload)

	c.logger.Log(logger.LogEntry{
		Level:   "info",
		Service: "track",
		Method:  "SendProceedToPayment",
		Message: fmt.Sprintf("[Roistat] 🚀 отправка: %s", string(body)),
	})

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"?key="+c.apiKey, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Log(logger.LogEntry{
			Level:   "error",
			Service: "track",
			Method:  "SendProceedToPayment",
			Message: "failed to create request",
			Error:   err,
		})
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		c.logger.Log(logger.LogEntry{
			Level:   "error",
			Service: "track",
			Method:  "SendProceedToPayment",
			Message: "http request failed",
			Error:   err,
		})
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	c.logger.Log(logger.LogEntry{
		Level:   "info",
		Service: "track",
		Method:  "SendProceedToPayment",
		Message: fmt.Sprintf("[Roistat] ✅ ответ (%d): %s", resp.StatusCode, string(respBody)),
	})

	if resp.StatusCode >= 300 {
		return fmt.Errorf("roistat error: %s", resp.Status)
	}
	return nil
}

// достаём значение по ключу из action.Meta (json.RawMessage)
func extractFromMeta(meta json.RawMessage, key string) string {
	if len(meta) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(meta, &m); err != nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
