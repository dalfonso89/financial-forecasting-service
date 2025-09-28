package models

import "time"

// HealthCheck represents the health check response
type HealthCheck struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ForecastRequest represents a request for financial forecasting
type ForecastRequest struct {
	BaseCurrency   string  `json:"base_currency" binding:"required"`
	TargetCurrency string  `json:"target_currency" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Periods        int     `json:"periods,omitempty"`       // Number of periods to forecast
	ForecastType   string  `json:"forecast_type,omitempty"` // "linear", "exponential", "moving_average"
}

// ForecastResponse represents a financial forecast response
type ForecastResponse struct {
	BaseCurrency    string           `json:"base_currency"`
	TargetCurrency  string           `json:"target_currency"`
	CurrentRate     float64          `json:"current_rate"`
	Amount          float64          `json:"amount"`
	ForecastType    string           `json:"forecast_type"`
	Periods         int              `json:"periods"`
	Forecasts       []ForecastPeriod `json:"forecasts"`
	GeneratedAt     time.Time        `json:"generated_at"`
	ConfidenceScore float64          `json:"confidence_score"`
}

// ForecastPeriod represents a single period in the forecast
type ForecastPeriod struct {
	Period        int     `json:"period"`
	Date          string  `json:"date"`
	Rate          float64 `json:"rate"`
	Amount        float64 `json:"amount"`
	Change        float64 `json:"change"`         // Change from previous period
	ChangePercent float64 `json:"change_percent"` // Percentage change from previous period
}

// TrendAnalysis represents trend analysis data
type TrendAnalysis struct {
	CurrencyPair   string    `json:"currency_pair"`
	Trend          string    `json:"trend"` // "upward", "downward", "sideways"
	Volatility     float64   `json:"volatility"`
	AverageRate    float64   `json:"average_rate"`
	MinRate        float64   `json:"min_rate"`
	MaxRate        float64   `json:"max_rate"`
	AnalysisPeriod int       `json:"analysis_period"`
	GeneratedAt    time.Time `json:"generated_at"`
}

// MultiCurrencyForecastRequest represents a request for multi-currency forecasting
type MultiCurrencyForecastRequest struct {
	BaseCurrency string   `json:"base_currency" binding:"required"`
	Currencies   []string `json:"currencies" binding:"required,min=1"`
	Amount       float64  `json:"amount" binding:"required,gt=0"`
	Periods      int      `json:"periods,omitempty"`
	ForecastType string   `json:"forecast_type,omitempty"`
}

// MultiCurrencyForecastResponse represents a multi-currency forecast response
type MultiCurrencyForecastResponse struct {
	BaseCurrency string                      `json:"base_currency"`
	Amount       float64                     `json:"amount"`
	ForecastType string                      `json:"forecast_type"`
	Periods      int                         `json:"periods"`
	Currencies   map[string][]ForecastPeriod `json:"currencies"`
	GeneratedAt  time.Time                   `json:"generated_at"`
}
