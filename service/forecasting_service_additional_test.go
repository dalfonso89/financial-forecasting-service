package service

import (
	"context"
	"testing"

	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/models"
)

func TestForecastingService_GenerateExponentialForecast(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	currentRate := 1.2
	req := &models.ForecastRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "EUR",
		Amount:         1000,
		Periods:        5,
	}

	forecasts, confidence := service.generateExponentialForecast(currentRate, req)

	if len(forecasts) != 5 {
		t.Errorf("Expected 5 forecasts, got %d", len(forecasts))
	}

	if confidence <= 0 || confidence > 1 {
		t.Errorf("Confidence score should be between 0 and 1, got %f", confidence)
	}

	// Check that forecasts are in ascending order by period
	for i := 1; i < len(forecasts); i++ {
		if forecasts[i].Period <= forecasts[i-1].Period {
			t.Errorf("Forecasts should be in ascending order by period")
		}
	}

	// Check that rates are increasing (exponential growth)
	for i := 1; i < len(forecasts); i++ {
		if forecasts[i].Rate <= forecasts[i-1].Rate {
			t.Errorf("Exponential forecast should have increasing rates")
		}
	}
}

func TestForecastingService_GenerateMovingAverageForecast(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	currentRate := 1.2
	req := &models.ForecastRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "EUR",
		Amount:         1000,
		Periods:        5,
	}

	forecasts, confidence := service.generateMovingAverageForecast(currentRate, req)

	if len(forecasts) != 5 {
		t.Errorf("Expected 5 forecasts, got %d", len(forecasts))
	}

	if confidence <= 0 || confidence > 1 {
		t.Errorf("Confidence score should be between 0 and 1, got %f", confidence)
	}

	// Check that forecasts are in ascending order by period
	for i := 1; i < len(forecasts); i++ {
		if forecasts[i].Period <= forecasts[i-1].Period {
			t.Errorf("Forecasts should be in ascending order by period")
		}
	}
}

func TestForecastingService_GenerateMultiCurrencyForecast(t *testing.T) {
	cfg := &config.Config{
		SupportedCurrencies:        []string{"USD", "EUR", "GBP"},
		DefaultForecastPeriods:     5,
		CurrencyExchangeServiceURL: "http://localhost:8081",
	}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	req := &models.MultiCurrencyForecastRequest{
		BaseCurrency: "USD",
		Currencies:   []string{"EUR", "GBP"},
		Amount:       1000,
		Periods:      5,
		ForecastType: "linear",
	}

	// This will fail because we can't actually call the currency service in tests
	// but we can test the validation logic
	_, err := service.GenerateMultiCurrencyForecast(context.Background(), req)
	if err == nil {
		t.Error("Expected error due to currency service unavailability, got nil")
	}
}

func TestForecastingService_AnalyzeTrend(t *testing.T) {
	cfg := &config.Config{
		SupportedCurrencies:        []string{"USD", "EUR"},
		CurrencyExchangeServiceURL: "http://localhost:8081",
	}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	// This will fail because we can't actually call the currency service in tests
	_, err := service.AnalyzeTrend(context.Background(), "USD", "EUR", 30)
	if err == nil {
		t.Error("Expected error due to currency service unavailability, got nil")
	}
}

func TestForecastingService_ClearCache(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	// Add something to cache first
	service.cacheMutex.Lock()
	service.cache["test_key"] = models.ForecastResponse{}
	service.cacheMutex.Unlock()

	// Verify cache has content
	service.cacheMutex.RLock()
	if len(service.cache) == 0 {
		t.Error("Expected cache to have content before clearing")
	}
	service.cacheMutex.RUnlock()

	// Clear cache
	service.ClearCache()

	// Verify cache is empty
	service.cacheMutex.RLock()
	if len(service.cache) != 0 {
		t.Error("Expected cache to be empty after clearing")
	}
	service.cacheMutex.RUnlock()
}

func TestForecastingService_GenerateCacheKey(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	req := &models.ForecastRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "EUR",
		Amount:         1000,
		Periods:        30,
		ForecastType:   "linear",
	}

	key1 := service.generateCacheKey(req)
	key2 := service.generateCacheKey(req)

	if key1 != key2 {
		t.Error("Expected same cache key for same request")
	}

	// Test different requests generate different keys
	req2 := &models.ForecastRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "GBP", // Different target currency
		Amount:         1000,
		Periods:        30,
		ForecastType:   "linear",
	}

	key3 := service.generateCacheKey(req2)
	if key1 == key3 {
		t.Error("Expected different cache keys for different requests")
	}
}

func TestForecastingService_ValidateForecastRequest_EdgeCases(t *testing.T) {
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
			name: "zero amount",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         0,
			},
			wantErr: true,
		},
		{
			name: "negative periods",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        -1,
			},
			wantErr: true,
		},
		{
			name: "exactly 365 periods",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        365,
			},
			wantErr: false,
		},
		{
			name: "more than 365 periods",
			request: &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        366,
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

func TestForecastingService_GenerateLinearForecast_EdgeCases(t *testing.T) {
	cfg := &config.Config{}
	loggerInstance := logger.New("debug")
	service := NewForecastingService(cfg, loggerInstance)

	tests := []struct {
		name        string
		currentRate float64
		periods     int
	}{
		{"single period", 1.2, 1},
		{"zero periods", 1.2, 0},
		{"large number of periods", 1.2, 100},
		{"very small rate", 0.001, 5},
		{"very large rate", 1000.0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &models.ForecastRequest{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
				Amount:         1000,
				Periods:        tt.periods,
			}

			forecasts, confidence := service.generateLinearForecast(tt.currentRate, req)

			if tt.periods == 0 {
				if len(forecasts) != 0 {
					t.Errorf("Expected 0 forecasts for 0 periods, got %d", len(forecasts))
				}
				return
			}

			if len(forecasts) != tt.periods {
				t.Errorf("Expected %d forecasts, got %d", tt.periods, len(forecasts))
			}

			if confidence <= 0 || confidence > 1 {
				t.Errorf("Confidence score should be between 0 and 1, got %f", confidence)
			}

			// Check that all forecasts have valid data
			for i, forecast := range forecasts {
				if forecast.Period != i+1 {
					t.Errorf("Expected period %d, got %d", i+1, forecast.Period)
				}
				if forecast.Rate <= 0 {
					t.Errorf("Expected positive rate, got %f", forecast.Rate)
				}
				if forecast.Amount <= 0 {
					t.Errorf("Expected positive amount, got %f", forecast.Amount)
				}
			}
		})
	}
}

func TestForecastingService_NewForecastingService(t *testing.T) {
	cfg := &config.Config{
		SupportedCurrencies:        []string{"USD", "EUR"},
		CurrencyExchangeServiceURL: "http://localhost:8081",
	}
	loggerInstance := logger.New("debug")

	service := NewForecastingService(cfg, loggerInstance)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.config != cfg {
		t.Error("Expected config to be set correctly")
	}

	if service.logger != loggerInstance {
		t.Error("Expected logger to be set correctly")
	}

	if service.currencyClient == nil {
		t.Error("Expected currency client to be created")
	}

	if service.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}
