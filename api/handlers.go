package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/middleware"
	"github.com/dalfonso89/financial-forecasting-service/models"
	"github.com/dalfonso89/financial-forecasting-service/service"
)

// HandlerConfig contains all dependencies for the Handlers
type HandlerConfig struct {
	Logger             logger.Logger
	ForecastingService *service.ForecastingService
}

// Handlers contains all HTTP handlers
type Handlers struct {
	logger             logger.Logger
	startTime          time.Time
	forecastingService *service.ForecastingService
}

// NewHandlers creates a new handlers instance with all dependencies
func NewHandlers(config HandlerConfig) *Handlers {
	return &Handlers{
		logger:             config.Logger,
		startTime:          time.Now(),
		forecastingService: config.ForecastingService,
	}
}

// SetupRoutes configures all the routes using Gin
func (handlers *Handlers) SetupRoutes() *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply middleware
	router.Use(middleware.RequestLogger(handlers.logger))
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RequestID())
	router.Use(handlers.corsMiddleware())

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Forecasting routes
		apiV1.POST("/forecast", handlers.GenerateForecast)
		apiV1.POST("/forecast/multi-currency", handlers.GenerateMultiCurrencyForecast)
		apiV1.GET("/forecast/trend/:base/:target", handlers.AnalyzeTrend)
		apiV1.DELETE("/forecast/cache", handlers.ClearCache)

		// Currency information routes
		apiV1.GET("/currencies", handlers.GetSupportedCurrencies)
		apiV1.GET("/currencies/rates/:base", handlers.GetCurrentRates)
	}

	return router
}

// HealthCheck handles health check requests
func (handlers *Handlers) HealthCheck(context *gin.Context) {
	healthCheckResponse := models.HealthCheck{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(handlers.startTime).String(),
	}

	context.JSON(http.StatusOK, healthCheckResponse)
}

// GenerateForecast handles single currency forecast requests
func (handlers *Handlers) GenerateForecast(context *gin.Context) {
	var req models.ForecastRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		handlers.writeErrorResponse(context, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	forecast, err := handlers.forecastingService.GenerateForecast(context.Request.Context(), &req)
	if err != nil {
		handlers.handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, forecast)
}

// GenerateMultiCurrencyForecast handles multi-currency forecast requests
func (handlers *Handlers) GenerateMultiCurrencyForecast(context *gin.Context) {
	var req models.MultiCurrencyForecastRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		handlers.writeErrorResponse(context, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	forecast, err := handlers.forecastingService.GenerateMultiCurrencyForecast(context.Request.Context(), &req)
	if err != nil {
		handlers.handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, forecast)
}

// AnalyzeTrend handles trend analysis requests
func (handlers *Handlers) AnalyzeTrend(context *gin.Context) {
	baseCurrency := context.Param("base")
	targetCurrency := context.Param("target")

	periodsStr := context.DefaultQuery("periods", "30")
	periods, err := strconv.Atoi(periodsStr)
	if err != nil {
		handlers.writeErrorResponse(context, http.StatusBadRequest, "invalid periods parameter", "periods must be a valid integer")
		return
	}

	analysis, err := handlers.forecastingService.AnalyzeTrend(context.Request.Context(), baseCurrency, targetCurrency, periods)
	if err != nil {
		handlers.handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, analysis)
}

// ClearCache handles cache clearing requests
func (handlers *Handlers) ClearCache(context *gin.Context) {
	handlers.forecastingService.ClearCache()
	context.JSON(http.StatusOK, gin.H{"message": "Cache cleared successfully"})
}

// GetSupportedCurrencies returns the list of supported currencies
func (handlers *Handlers) GetSupportedCurrencies(context *gin.Context) {
	// This would typically come from configuration
	supportedCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"}
	context.JSON(http.StatusOK, gin.H{"currencies": supportedCurrencies})
}

// GetCurrentRates fetches current exchange rates from the currency service
func (handlers *Handlers) GetCurrentRates(context *gin.Context) {
	baseCurrency := context.Param("base")

	// This would typically use the currency client directly
	// For now, we'll return a placeholder response
	context.JSON(http.StatusOK, gin.H{
		"message": "Current rates endpoint - would fetch from currency exchange service",
		"base":    baseCurrency,
	})
}

// writeErrorResponse writes an error response using Gin context
func (handlers *Handlers) writeErrorResponse(context *gin.Context, statusCode int, errorMessage, errorDetails string) {
	errorResponse := models.ErrorResponse{
		Error:   errorMessage,
		Message: errorDetails,
		Code:    statusCode,
	}

	context.JSON(statusCode, errorResponse)
}

// handleServiceError handles service errors
func (handlers *Handlers) handleServiceError(context *gin.Context, err error) {
	handlers.logger.Errorf("Service error: %v", err)
	handlers.writeErrorResponse(context, http.StatusInternalServerError, "service error", err.Error())
}

// corsMiddleware adds CORS headers using Gin middleware
func (handlers *Handlers) corsMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle HTTP method using type switch
		switch context.Request.Method {
		case "OPTIONS":
			context.AbortWithStatus(http.StatusOK)
			return
		case "GET", "POST", "PUT", "DELETE":
			// Continue processing
		default:
			context.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		context.Next()
	}
}

