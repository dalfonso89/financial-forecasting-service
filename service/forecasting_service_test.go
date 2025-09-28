package service

import (
	"context"
	"testing"

	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/models"
)

// TestNewForecastingService tests the service constructor
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

// TestForecastingService_GenerateForecast tests the main forecast generation functionality
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

// TestForecastingService_GenerateMultiCurrencyForecast tests multi-currency forecast generation
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

// TestForecastingService_AnalyzeTrend tests trend analysis functionality
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

// TestForecastingService_validateForecastRequest tests request validation
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

// TestForecastingService_generateLinearForecast tests linear forecast generation
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

// TestForecastingService_generateLinearForecast_EdgeCases tests edge cases for linear forecasting
func TestForecastingService_generateLinearForecast_EdgeCases(t *testing.T) {
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

// TestForecastingService_generateExponentialForecast tests exponential forecast generation
func TestForecastingService_generateExponentialForecast(t *testing.T) {
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

// TestForecastingService_generateMovingAverageForecast tests moving average forecast generation
func TestForecastingService_generateMovingAverageForecast(t *testing.T) {
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

// TestForecastingService_isCurrencySupported tests currency support validation
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

// TestForecastingService_ClearCache tests cache clearing functionality
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

// TestForecastingService_generateCacheKey tests cache key generation
func TestForecastingService_generateCacheKey(t *testing.T) {
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
