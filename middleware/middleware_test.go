package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check security headers
	expectedHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Content-Security-Policy":   "default-src 'self'",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actualValue)
		}
	}
}

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(200, gin.H{"request_id": requestID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check that X-Request-ID header is set
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Check that request_id is set in context
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, requestID) {
		t.Errorf("Expected response to contain request_id %s, got: %s", requestID, responseBody)
	}
}

func TestRequestID_ExistingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(200, gin.H{"request_id": requestID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "existing-request-id")
	router.ServeHTTP(w, req)

	// Check that existing X-Request-ID header is preserved
	requestID := w.Header().Get("X-Request-ID")
	if requestID != "existing-request-id" {
		t.Errorf("Expected existing request ID to be preserved, got: %s", requestID)
	}
}

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	loggerInstance := logger.New("debug")
	router := gin.New()
	router.Use(RequestLogger(loggerInstance))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// The RequestLogger should not cause any errors
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGenerateRequestID(t *testing.T) {
	requestID := generateRequestID()

	if requestID == "" {
		t.Error("Expected request ID to be generated")
	}

	// Check format: should be timestamp-8chars
	parts := strings.Split(requestID, "-")
	if len(parts) != 2 {
		t.Errorf("Expected request ID format 'timestamp-random', got: %s", requestID)
	}

	if len(parts[1]) != 8 {
		t.Errorf("Expected random part to be 8 characters, got: %d", len(parts[1]))
	}
}

func TestRandomString(t *testing.T) {
	tests := []int{1, 5, 10, 20}

	for _, length := range tests {
		result := randomString(length)

		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}

		// Check that it only contains valid characters
		validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for _, char := range result {
			if !strings.ContainsRune(validChars, char) {
				t.Errorf("Expected only valid characters, got: %c", char)
			}
		}
	}
}

func TestRandomString_ZeroLength(t *testing.T) {
	result := randomString(0)
	if result != "" {
		t.Errorf("Expected empty string for length 0, got: %s", result)
	}
}

func TestSecurityHeaders_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Make multiple requests
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Request %d: Expected status 200, got %d", i, w.Code)
		}

		// Check that security headers are present in each request
		if w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("Request %d: Expected X-Content-Type-Options header", i)
		}
	}
}

func TestRequestID_Uniqueness(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(200, gin.H{"request_id": requestID})
	})

	// Generate multiple request IDs with small delays to ensure different timestamps
	requestIDs := make(map[string]bool)
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		if requestIDs[requestID] {
			t.Errorf("Duplicate request ID generated: %s", requestID)
		}
		requestIDs[requestID] = true

		// Small delay to ensure different timestamps
		time.Sleep(1 * time.Millisecond)
	}
}
