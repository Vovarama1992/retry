package models

import "time"

// Корневой сценарий: "Получить доступ"
type ScenarioGetAccessSummary struct {
	TotalVisits              int                           `json:"total_visits"`        // всего визитов
	VisitsWithClick          int                           `json:"visits_with_click"`   // сколько визитов дошли до кнопки
	ConversionFromAll        float64                       `json:"conversion_from_all"` // % от всех визитов
	SessionIndexDistribution map[int]int                   `json:"session_index_distribution"`
	Details                  []ScenarioGetAccessDetail     `json:"details"`
	ProceedToPayment         ScenarioProceedToPaymentStats `json:"proceed_to_payment"`
}

// Детали визитов, где был клик "Получить доступ"
type ScenarioGetAccessDetail struct {
	VisitID      string    `json:"visit_id"`
	SessionIndex int       `json:"session_index"`
	ClickedAt    time.Time `json:"clicked_at"`
}

// Вложенный сценарий: "Перейти к оплате"
type ScenarioProceedToPaymentStats struct {
	VisitsWithPayment          int                              `json:"visitsWithPayment"`
	ConversionFromClicks       float64                          `json:"conversionFromClicks"`
	PaymentMethodsDistribution map[string]int                   `json:"paymentMethodsDistribution"`
	Details                    []ScenarioProceedToPaymentDetail `json:"details"`

	AvgMinutes float64 `json:"avgMinutes"`
	P50Minutes float64 `json:"p50Minutes"`
}

type ScenarioProceedToPaymentDetail struct {
	VisitID       string    `json:"visitId"`
	SessionIndex  int       `json:"sessionIndex"`
	ClickedAt     time.Time `json:"clickedAt"`
	PaymentMethod string    `json:"paymentMethod"`

	// если «оплата» была в той же сессии после CTA — разница в минутах
	DurationMinutes *float64 `json:"durationMinutes,omitempty"`
}
