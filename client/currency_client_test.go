package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	currencymodels "github.com/dalfonso89/currency-exchange-service/models"
	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
)

func TestNewCurrencyClient(t *testing.T) {
	cfg := &config.Config{
		CurrencyExchangeServiceURL: "http://localhost:8081",
		CurrencyExchangeTimeout:    30 * time.Second,
	}
	loggerInstance := logger.New("debug")

	client := NewCurrencyClient(cfg, loggerInstance)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.baseURL != cfg.CurrencyExchangeServiceURL {
		t.Errorf("Expected baseURL %s, got %s", cfg.CurrencyExchangeServiceURL, client.baseURL)
	}

	if client.httpClient.Timeout != cfg.CurrencyExchangeTimeout {
		t.Errorf("Expected timeout %v, got %v", cfg.CurrencyExchangeTimeout, client.httpClient.Timeout)
	}

	if client.logger != loggerInstance {
		t.Error("Expected logger to be set correctly")
	}
}

func TestCurrencyClient_GetRates_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/rates/USD" {
			t.Errorf("Expected path /api/v1/rates/USD, got %s", r.URL.Path)
		}

		// Response data for testing
		_ = currencymodels.RatesResponse{
			Base:      "USD",
			Timestamp: 1640995200,
			Rates: map[string]float64{
				"EUR": 0.85,
				"GBP": 0.73,
			},
			Provider: "test",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","timestamp":1640995200,"rates":{"EUR":0.85,"GBP":0.73},"provider":"test"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	rates, err := client.GetRates(context.Background(), "USD")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rates == nil {
		t.Fatal("Expected rates response, got nil")
	}

	if rates.Base != "USD" {
		t.Errorf("Expected base USD, got %s", rates.Base)
	}

	if rates.Timestamp != 1640995200 {
		t.Errorf("Expected timestamp 1640995200, got %d", rates.Timestamp)
	}

	if len(rates.Rates) != 2 {
		t.Errorf("Expected 2 rates, got %d", len(rates.Rates))
	}

	if rates.Rates["EUR"] != 0.85 {
		t.Errorf("Expected EUR rate 0.85, got %f", rates.Rates["EUR"])
	}

	if rates.Provider != "test" {
		t.Errorf("Expected provider test, got %s", rates.Provider)
	}
}

func TestCurrencyClient_GetRates_HTTPError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	rates, err := client.GetRates(context.Background(), "USD")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if rates != nil {
		t.Error("Expected rates to be nil on error")
	}

	expectedError := "currency service returned status 500: Internal Server Error"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCurrencyClient_GetRates_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	rates, err := client.GetRates(context.Background(), "USD")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if rates != nil {
		t.Error("Expected rates to be nil on error")
	}

	expectedError := "failed to unmarshal response"
	if err.Error()[:len(expectedError)] != expectedError {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCurrencyClient_GetRatesWithQuery_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/rates" {
			t.Errorf("Expected path /api/v1/rates, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("base") != "USD" {
			t.Errorf("Expected query base=USD, got %s", r.URL.Query().Get("base"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","timestamp":1640995200,"rates":{"EUR":0.85},"provider":"test"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	rates, err := client.GetRatesWithQuery(context.Background(), "USD")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rates == nil {
		t.Fatal("Expected rates response, got nil")
	}

	if rates.Base != "USD" {
		t.Errorf("Expected base USD, got %s", rates.Base)
	}
}

func TestCurrencyClient_HealthCheck_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected path /health, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	err := client.HealthCheck(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCurrencyClient_HealthCheck_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.Config{
		CurrencyExchangeServiceURL: server.URL,
		CurrencyExchangeTimeout:    5 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	err := client.HealthCheck(context.Background())

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "currency service health check failed with status: 500"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCurrencyClient_ContextCancellation(t *testing.T) {
	cfg := &config.Config{
		CurrencyExchangeServiceURL: "http://localhost:9999", // Non-existent server
		CurrencyExchangeTimeout:    1 * time.Second,
	}
	loggerInstance := logger.New("debug")
	client := NewCurrencyClient(cfg, loggerInstance)

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	rates, err := client.GetRates(ctx, "USD")

	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}

	if rates != nil {
		t.Error("Expected rates to be nil on context cancellation")
	}
}
