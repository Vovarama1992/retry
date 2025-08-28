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
	VisitsWithPayment          int                              `json:"visits_with_payment"`
	ConversionFromClicks       float64                          `json:"conversion_from_clicks"` // % относительно "получить доступ"
	PaymentMethodsDistribution map[string]int                   `json:"payment_methods_distribution"`
	Details                    []ScenarioProceedToPaymentDetail `json:"details"`
}

type ScenarioProceedToPaymentDetail struct {
	VisitID       string    `json:"visit_id"`
	SessionIndex  int       `json:"session_index"`
	ClickedAt     time.Time `json:"clicked_at"`
	PaymentMethod string    `json:"payment_method,omitempty"`
}
