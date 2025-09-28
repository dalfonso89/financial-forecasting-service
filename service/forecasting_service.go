package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dalfonso89/financial-forecasting-service/client"
	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
	"github.com/dalfonso89/financial-forecasting-service/models"
)

// ForecastingService handles financial forecasting operations
type ForecastingService struct {
	config         *config.Config
	logger         logger.Logger
	currencyClient *client.CurrencyClient

	// Cache for forecasts
	cacheMutex sync.RWMutex
	cache      map[string]models.ForecastResponse
}

// NewForecastingService creates a new forecasting service
func NewForecastingService(cfg *config.Config, logger logger.Logger) *ForecastingService {
	return &ForecastingService{
		config:         cfg,
		logger:         logger,
		currencyClient: client.NewCurrencyClient(cfg, logger),
		cache:          make(map[string]models.ForecastResponse),
	}
}

// GenerateForecast generates a financial forecast for a currency pair
func (fs *ForecastingService) GenerateForecast(ctx context.Context, req *models.ForecastRequest) (*models.ForecastResponse, error) {
	// Validate request
	if err := fs.validateForecastRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Set defaults
	if req.Periods == 0 {
		req.Periods = fs.config.DefaultForecastPeriods
	}
	if req.ForecastType == "" {
		req.ForecastType = "linear"
	}

	// Check cache first
	cacheKey := fs.generateCacheKey(req)
	fs.cacheMutex.RLock()
	if cached, exists := fs.cache[cacheKey]; exists {
		fs.cacheMutex.RUnlock()
		fs.logger.Debugf("Returning cached forecast for %s/%s", req.BaseCurrency, req.TargetCurrency)
		return &cached, nil
	}
	fs.cacheMutex.RUnlock()

	// Fetch current exchange rates
	rates, err := fs.currencyClient.GetRates(ctx, req.BaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}

	// Get current rate for target currency
	currentRate, exists := rates.Rates[req.TargetCurrency]
	if !exists {
		return nil, fmt.Errorf("target currency %s not found in exchange rates", req.TargetCurrency)
	}

	// Generate forecast based on type
	var forecasts []models.ForecastPeriod
	var confidenceScore float64

	switch req.ForecastType {
	case "linear":
		forecasts, confidenceScore = fs.generateLinearForecast(currentRate, req)
	case "exponential":
		forecasts, confidenceScore = fs.generateExponentialForecast(currentRate, req)
	case "moving_average":
		forecasts, confidenceScore = fs.generateMovingAverageForecast(currentRate, req)
	default:
		return nil, fmt.Errorf("unsupported forecast type: %s", req.ForecastType)
	}

	// Create response
	response := &models.ForecastResponse{
		BaseCurrency:    req.BaseCurrency,
		TargetCurrency:  req.TargetCurrency,
		CurrentRate:     currentRate,
		Amount:          req.Amount,
		ForecastType:    req.ForecastType,
		Periods:         req.Periods,
		Forecasts:       forecasts,
		GeneratedAt:     time.Now(),
		ConfidenceScore: confidenceScore,
	}

	// Cache the result
	fs.cacheMutex.Lock()
	fs.cache[cacheKey] = *response
	fs.cacheMutex.Unlock()

	fs.logger.Infof("Generated %s forecast for %s/%s with %d periods", req.ForecastType, req.BaseCurrency, req.TargetCurrency, req.Periods)
	return response, nil
}

// GenerateMultiCurrencyForecast generates forecasts for multiple currencies
func (fs *ForecastingService) GenerateMultiCurrencyForecast(ctx context.Context, req *models.MultiCurrencyForecastRequest) (*models.MultiCurrencyForecastResponse, error) {
	// Set defaults
	if req.Periods == 0 {
		req.Periods = fs.config.DefaultForecastPeriods
	}
	if req.ForecastType == "" {
		req.ForecastType = "linear"
	}

	// Fetch current exchange rates
	rates, err := fs.currencyClient.GetRates(ctx, req.BaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}

	// Generate forecasts for each currency
	currencyForecasts := make(map[string][]models.ForecastPeriod)

	for _, currency := range req.Currencies {
		rate, exists := rates.Rates[currency]
		if !exists {
			fs.logger.Warnf("Currency %s not found in exchange rates, skipping", currency)
			continue
		}

		forecastReq := &models.ForecastRequest{
			BaseCurrency:   req.BaseCurrency,
			TargetCurrency: currency,
			Amount:         req.Amount,
			Periods:        req.Periods,
			ForecastType:   req.ForecastType,
		}

		var forecasts []models.ForecastPeriod
		switch req.ForecastType {
		case "linear":
			forecasts, _ = fs.generateLinearForecast(rate, forecastReq)
		case "exponential":
			forecasts, _ = fs.generateExponentialForecast(rate, forecastReq)
		case "moving_average":
			forecasts, _ = fs.generateMovingAverageForecast(rate, forecastReq)
		}

		currencyForecasts[currency] = forecasts
	}

	response := &models.MultiCurrencyForecastResponse{
		BaseCurrency: req.BaseCurrency,
		Amount:       req.Amount,
		ForecastType: req.ForecastType,
		Periods:      req.Periods,
		Currencies:   currencyForecasts,
		GeneratedAt:  time.Now(),
	}

	fs.logger.Infof("Generated multi-currency forecast for %d currencies", len(currencyForecasts))
	return response, nil
}

// AnalyzeTrend analyzes the trend for a currency pair
func (fs *ForecastingService) AnalyzeTrend(ctx context.Context, baseCurrency, targetCurrency string, periods int) (*models.TrendAnalysis, error) {
	// For now, we'll use a simple analysis based on current rates
	// In a real implementation, you might want to fetch historical data
	rates, err := fs.currencyClient.GetRates(ctx, baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}

	rate, exists := rates.Rates[targetCurrency]
	if !exists {
		return nil, fmt.Errorf("target currency %s not found in exchange rates", targetCurrency)
	}

	// Simple trend analysis (in a real implementation, you'd use historical data)
	analysis := &models.TrendAnalysis{
		CurrencyPair:   fmt.Sprintf("%s/%s", baseCurrency, targetCurrency),
		Trend:          "sideways", // Placeholder
		Volatility:     0.05,       // Placeholder
		AverageRate:    rate,
		MinRate:        rate * 0.95,
		MaxRate:        rate * 1.05,
		AnalysisPeriod: periods,
		GeneratedAt:    time.Now(),
	}

	return analysis, nil
}

// validateForecastRequest validates the forecast request
func (fs *ForecastingService) validateForecastRequest(req *models.ForecastRequest) error {
	if req.BaseCurrency == "" {
		return fmt.Errorf("base currency is required")
	}
	if req.TargetCurrency == "" {
		return fmt.Errorf("target currency is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if req.Periods < 0 {
		return fmt.Errorf("periods cannot be negative")
	}
	if req.Periods > 365 {
		return fmt.Errorf("periods cannot exceed 365")
	}

	// Check if currencies are supported
	if !fs.isCurrencySupported(req.BaseCurrency) {
		return fmt.Errorf("base currency %s is not supported", req.BaseCurrency)
	}
	if !fs.isCurrencySupported(req.TargetCurrency) {
		return fmt.Errorf("target currency %s is not supported", req.TargetCurrency)
	}

	return nil
}

// isCurrencySupported checks if a currency is supported
func (fs *ForecastingService) isCurrencySupported(currency string) bool {
	for _, supported := range fs.config.SupportedCurrencies {
		if supported == currency {
			return true
		}
	}
	return false
}

// generateCacheKey generates a cache key for the request
func (fs *ForecastingService) generateCacheKey(req *models.ForecastRequest) string {
	return fmt.Sprintf("%s_%s_%s_%d_%d", req.BaseCurrency, req.TargetCurrency, req.ForecastType, int(req.Amount), req.Periods)
}

// generateLinearForecast generates a linear forecast
func (fs *ForecastingService) generateLinearForecast(currentRate float64, req *models.ForecastRequest) ([]models.ForecastPeriod, float64) {
	forecasts := make([]models.ForecastPeriod, req.Periods)

	// Simple linear trend (in a real implementation, you'd use more sophisticated algorithms)
	trend := 0.001 // 0.1% change per period

	for i := 0; i < req.Periods; i++ {
		period := i + 1
		rate := currentRate * (1 + trend*float64(period))
		amount := req.Amount * rate

		var change, changePercent float64
		if i > 0 {
			prevRate := currentRate * (1 + trend*float64(i))
			change = rate - prevRate
			changePercent = (change / prevRate) * 100
		}

		forecasts[i] = models.ForecastPeriod{
			Period:        period,
			Date:          time.Now().AddDate(0, 0, period).Format("2006-01-02"),
			Rate:          math.Round(rate*10000) / 10000, // Round to 4 decimal places
			Amount:        math.Round(amount*100) / 100,   // Round to 2 decimal places
			Change:        math.Round(change*10000) / 10000,
			ChangePercent: math.Round(changePercent*100) / 100,
		}
	}

	confidenceScore := 0.7 // Placeholder confidence score
	return forecasts, confidenceScore
}

// generateExponentialForecast generates an exponential forecast
func (fs *ForecastingService) generateExponentialForecast(currentRate float64, req *models.ForecastRequest) ([]models.ForecastPeriod, float64) {
	forecasts := make([]models.ForecastPeriod, req.Periods)

	// Simple exponential trend
	growthRate := 0.002 // 0.2% growth per period

	for i := 0; i < req.Periods; i++ {
		period := i + 1
		rate := currentRate * math.Pow(1+growthRate, float64(period))
		amount := req.Amount * rate

		var change, changePercent float64
		if i > 0 {
			prevRate := currentRate * math.Pow(1+growthRate, float64(i))
			change = rate - prevRate
			changePercent = (change / prevRate) * 100
		}

		forecasts[i] = models.ForecastPeriod{
			Period:        period,
			Date:          time.Now().AddDate(0, 0, period).Format("2006-01-02"),
			Rate:          math.Round(rate*10000) / 10000,
			Amount:        math.Round(amount*100) / 100,
			Change:        math.Round(change*10000) / 10000,
			ChangePercent: math.Round(changePercent*100) / 100,
		}
	}

	confidenceScore := 0.6 // Placeholder confidence score
	return forecasts, confidenceScore
}

// generateMovingAverageForecast generates a moving average forecast
func (fs *ForecastingService) generateMovingAverageForecast(currentRate float64, req *models.ForecastRequest) ([]models.ForecastPeriod, float64) {
	forecasts := make([]models.ForecastPeriod, req.Periods)

	// Simple moving average with some volatility
	baseRate := currentRate
	volatility := 0.01 // 1% volatility

	for i := 0; i < req.Periods; i++ {
		period := i + 1
		// Add some random-like variation based on period
		variation := math.Sin(float64(period)*0.1) * volatility
		rate := baseRate * (1 + variation)
		amount := req.Amount * rate

		var change, changePercent float64
		if i > 0 {
			prevVariation := math.Sin(float64(i)*0.1) * volatility
			prevRate := baseRate * (1 + prevVariation)
			change = rate - prevRate
			changePercent = (change / prevRate) * 100
		}

		forecasts[i] = models.ForecastPeriod{
			Period:        period,
			Date:          time.Now().AddDate(0, 0, period).Format("2006-01-02"),
			Rate:          math.Round(rate*10000) / 10000,
			Amount:        math.Round(amount*100) / 100,
			Change:        math.Round(change*10000) / 10000,
			ChangePercent: math.Round(changePercent*100) / 100,
		}
	}

	confidenceScore := 0.5 // Placeholder confidence score
	return forecasts, confidenceScore
}

// ClearCache clears the forecast cache
func (fs *ForecastingService) ClearCache() {
	fs.cacheMutex.Lock()
	defer fs.cacheMutex.Unlock()
	fs.cache = make(map[string]models.ForecastResponse)
	fs.logger.Info("Forecast cache cleared")
}
