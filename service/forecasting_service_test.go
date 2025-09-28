package service

import (
	"context"
	"testing"

	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/models"
)

func TestForecastingService_GenerateForecast(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		SupportedCurrencies:        []string{"USD", "EUR", "GBP"},
		DefaultForecastPeriods:     5,
		CurrencyExchangeServiceURL: "http://localhost:8081", // Mock URL
	}

	// Create test logger
	loggerInstance := logger.New("debug")

	// Create forecasting service
	service := NewForecastingService(cfg, loggerInstance)

	// Test cases
	tests := []struct {
		name    string
		request *models.ForecastRequest
		wantErr bool
	}{
		{
			name: "invalid base currency",
			request: &models.ForecastRequest{
				BaseCurrency:   "INVALID",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        5,
				ForecastType:   "linear",
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         -100,
				Periods:        5,
				ForecastType:   "linear",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GenerateForecast(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateForecast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForecastingService_validateForecastRequest(t *testing.T) {
	cfg := &config.Config{
		SupportedCurrencies: []string{"USD", "EUR", "GBP"},
	}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	tests := []struct {
		name    string
		request *models.ForecastRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        5,
			},
			wantErr: false,
		},
		{
			name: "empty base currency",
			request: &models.ForecastRequest{
				BaseCurrency:   "",
				TargetCurrency: "EUR",
				Amount:         1000,
			},
			wantErr: true,
		},
		{
			name: "unsupported currency",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "INVALID",
				Amount:         1000,
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         -100,
			},
			wantErr: true,
		},
		{
			name: "too many periods",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        400,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateForecastRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateForecastRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForecastingService_generateLinearForecast(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	req := &models.ForecastRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "EUR",
		Amount:         1000,
		Periods:        5,
	}

	forecasts, confidence := service.generateLinearForecast(1.2, req)

	if len(forecasts) != 5 {
		t.Errorf("Expected 5 forecasts, got %d", len(forecasts))
	}

	if confidence <= 0 || confidence > 1 {
		t.Errorf("Confidence score should be between 0 and 1, got %f", confidence)
	}

	// Check that forecasts are in order
	for i := 1; i < len(forecasts); i++ {
		if forecasts[i].Period <= forecasts[i-1].Period {
			t.Errorf("Forecasts should be in ascending order by period")
		}
	}
}

func TestForecastingService_isCurrencySupported(t *testing.T) {
	cfg := &config.Config{
		SupportedCurrencies: []string{"USD", "EUR", "GBP"},
	}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	tests := []struct {
		currency string
		expected bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"JPY", false},
		{"", false},
	}

	for _, tt := range tests {
		result := service.isCurrencySupported(tt.currency)
		if result != tt.expected {
			t.Errorf("isCurrencySupported(%s) = %v, expected %v", tt.currency, result, tt.expected)
		}
	}
}
