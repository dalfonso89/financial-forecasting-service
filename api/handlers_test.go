package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/models"
	"github.com/dalfonso89/financial-forecasting-service/service"
)

// Test helper function to create handlers
func createTestHandlers() *Handlers {
	gin.SetMode(gin.TestMode)
	loggerInstance := logger.New("debug")
	cfg := &config.Config{}
	forecastingService := service.NewForecastingService(cfg, loggerInstance)

	return &Handlers{
		logger:             loggerInstance,
		forecastingService: forecastingService,
	}
}

func TestHandlers_HealthCheck(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/health", handlers.HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.HealthCheck
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
	if response.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", response.Version)
	}
}

func TestHandlers_GetSupportedCurrencies(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/currencies", handlers.GetSupportedCurrencies)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/currencies", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	currencies, ok := response["currencies"].([]interface{})
	if !ok {
		t.Errorf("Expected currencies array in response")
	}
	if len(currencies) == 0 {
		t.Errorf("Expected non-empty currencies array")
	}
}

func TestHandlers_GenerateForecast_InvalidJSON(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.POST("/forecast", handlers.GenerateForecast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forecast", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse.Error != "invalid request" {
		t.Errorf("Expected error 'invalid request', got '%s'", errorResponse.Error)
	}
}

func TestHandlers_GenerateForecast_EmptyBody(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.POST("/forecast", handlers.GenerateForecast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forecast", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandlers_GenerateMultiCurrencyForecast_InvalidJSON(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.POST("/forecast/multi-currency", handlers.GenerateMultiCurrencyForecast)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forecast/multi-currency", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandlers_GenerateMultiCurrencyForecast_EmptyCurrencies(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.POST("/forecast/multi-currency", handlers.GenerateMultiCurrencyForecast)

	request := models.MultiCurrencyForecastRequest{
		BaseCurrency: "USD",
		Currencies:   []string{}, // Empty currencies
		Amount:       1000,
	}

	jsonData, _ := json.Marshal(request)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/forecast/multi-currency", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandlers_AnalyzeTrend_InvalidPeriods(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/forecast/trend/:base/:target", handlers.AnalyzeTrend)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/forecast/trend/USD/EUR?periods=invalid", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse.Error != "invalid periods parameter" {
		t.Errorf("Expected error 'invalid periods parameter', got '%s'", errorResponse.Error)
	}
}

func TestHandlers_AnalyzeTrend_DefaultPeriods(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/forecast/trend/:base/:target", handlers.AnalyzeTrend)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/forecast/trend/USD/EUR", nil)
	router.ServeHTTP(w, req)

	// This will fail due to currency service unavailability, but we can test the parameter parsing
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 due to service unavailability, got %d", w.Code)
	}
}

func TestHandlers_GetCurrentRates(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/currencies/rates/:base", handlers.GetCurrentRates)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/currencies/rates/USD", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["base"] != "USD" {
		t.Errorf("Expected base USD, got %v", response["base"])
	}
}

func TestHandlers_CORS_Middleware(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.Use(handlers.corsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test OPTIONS request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	// Test GET request with CORS headers
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check CORS headers
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actualValue)
		}
	}
}

func TestHandlers_CORS_InvalidMethod(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.Use(handlers.corsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for invalid method, got %d", w.Code)
	}
}

func TestHandlers_WriteErrorResponse(t *testing.T) {
	handlers := createTestHandlers()
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		handlers.writeErrorResponse(c, http.StatusBadRequest, "test error", "test details")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var errorResponse models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", errorResponse.Error)
	}

	if errorResponse.Message != "test details" {
		t.Errorf("Expected message 'test details', got '%s'", errorResponse.Message)
	}

	if errorResponse.Code != http.StatusBadRequest {
		t.Errorf("Expected code %d, got %d", http.StatusBadRequest, errorResponse.Code)
	}
}
