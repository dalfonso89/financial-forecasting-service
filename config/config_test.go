package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear any existing environment variables
	os.Clearenv()

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be created, got nil")
	}

	// Test default values
	if config.Port != "8082" {
		t.Errorf("Expected default port 8082, got %s", config.Port)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level info, got %s", config.LogLevel)
	}

	if config.CurrencyExchangeServiceURL != "http://localhost:8081" {
		t.Errorf("Expected default currency service URL http://localhost:8081, got %s", config.CurrencyExchangeServiceURL)
	}

	if config.CurrencyExchangeTimeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.CurrencyExchangeTimeout)
	}

	if config.ForecastCacheTTL != 300*time.Second {
		t.Errorf("Expected default cache TTL 300s, got %v", config.ForecastCacheTTL)
	}

	if config.MaxConcurrentRequests != 10 {
		t.Errorf("Expected default max concurrent requests 10, got %d", config.MaxConcurrentRequests)
	}

	if config.DefaultForecastPeriods != 30 {
		t.Errorf("Expected default forecast periods 30, got %d", config.DefaultForecastPeriods)
	}

	// Test supported currencies
	expectedCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"}
	if len(config.SupportedCurrencies) != len(expectedCurrencies) {
		t.Errorf("Expected %d supported currencies, got %d", len(expectedCurrencies), len(config.SupportedCurrencies))
	}

	for i, expected := range expectedCurrencies {
		if config.SupportedCurrencies[i] != expected {
			t.Errorf("Expected currency %s at index %d, got %s", expected, i, config.SupportedCurrencies[i])
		}
	}
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("CURRENCY_EXCHANGE_SERVICE_URL", "http://test:8081")
	os.Setenv("CURRENCY_EXCHANGE_TIMEOUT_SECONDS", "60")
	os.Setenv("FORECAST_CACHE_TTL_SECONDS", "600")
	os.Setenv("MAX_CONCURRENT_REQUESTS", "20")
	os.Setenv("DEFAULT_FORECAST_PERIODS", "60")

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test environment variable values
	if config.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Port)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", config.LogLevel)
	}

	if config.CurrencyExchangeServiceURL != "http://test:8081" {
		t.Errorf("Expected currency service URL http://test:8081, got %s", config.CurrencyExchangeServiceURL)
	}

	if config.CurrencyExchangeTimeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", config.CurrencyExchangeTimeout)
	}

	if config.ForecastCacheTTL != 600*time.Second {
		t.Errorf("Expected cache TTL 600s, got %v", config.ForecastCacheTTL)
	}

	if config.MaxConcurrentRequests != 20 {
		t.Errorf("Expected max concurrent requests 20, got %d", config.MaxConcurrentRequests)
	}

	if config.DefaultForecastPeriods != 60 {
		t.Errorf("Expected forecast periods 60, got %d", config.DefaultForecastPeriods)
	}

	// Clean up
	os.Clearenv()
}

func TestLoad_InvalidNumericValues(t *testing.T) {
	// Set invalid numeric environment variables
	os.Setenv("CURRENCY_EXCHANGE_TIMEOUT_SECONDS", "invalid")
	os.Setenv("FORECAST_CACHE_TTL_SECONDS", "not_a_number")
	os.Setenv("MAX_CONCURRENT_REQUESTS", "abc")
	os.Setenv("DEFAULT_FORECAST_PERIODS", "xyz")

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fall back to default values
	if config.CurrencyExchangeTimeout != 60*time.Second { // Default fallback
		t.Errorf("Expected fallback timeout 60s, got %v", config.CurrencyExchangeTimeout)
	}

	if config.ForecastCacheTTL != 60*time.Second { // Default fallback
		t.Errorf("Expected fallback cache TTL 60s, got %v", config.ForecastCacheTTL)
	}

	if config.MaxConcurrentRequests != 60 { // Default fallback
		t.Errorf("Expected fallback max concurrent requests 60, got %d", config.MaxConcurrentRequests)
	}

	if config.DefaultForecastPeriods != 60 { // Default fallback
		t.Errorf("Expected fallback forecast periods 60, got %d", config.DefaultForecastPeriods)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected test_value, got %s", result)
	}

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTING_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected default_value, got %s", result)
	}

	// Test with empty environment variable
	os.Setenv("EMPTY_VAR", "")
	result = getEnv("EMPTY_VAR", "default")
	if result != "default" {
		t.Errorf("Expected default, got %s", result)
	}

	// Clean up
	os.Clearenv()
}

func TestMustAtoi(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"-456", -456},
		{"invalid", 60}, // Default fallback
		{"", 60},        // Default fallback
		{"abc", 60},     // Default fallback
	}

	for _, test := range tests {
		result := mustAtoi(test.input)
		if result != test.expected {
			t.Errorf("mustAtoi(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

