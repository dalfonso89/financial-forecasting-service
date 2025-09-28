package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the financial forecasting service
type Config struct {
	Port     string
	LogLevel string

	// Currency exchange service configuration
	CurrencyExchangeServiceURL string
	CurrencyExchangeTimeout    time.Duration

	// Forecasting configuration
	ForecastCacheTTL       time.Duration
	MaxConcurrentRequests  int
	DefaultForecastPeriods int
	SupportedCurrencies    []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Port:     getEnv("PORT", "8082"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		CurrencyExchangeServiceURL: getEnv("CURRENCY_EXCHANGE_SERVICE_URL", "http://localhost:8081"),
		CurrencyExchangeTimeout:    time.Duration(mustAtoi(getEnv("CURRENCY_EXCHANGE_TIMEOUT_SECONDS", "30"))) * time.Second,

		ForecastCacheTTL:       time.Duration(mustAtoi(getEnv("FORECAST_CACHE_TTL_SECONDS", "300"))) * time.Second, // 5 minutes
		MaxConcurrentRequests:  mustAtoi(getEnv("MAX_CONCURRENT_REQUESTS", "10")),
		DefaultForecastPeriods: mustAtoi(getEnv("DEFAULT_FORECAST_PERIODS", "30")),
		SupportedCurrencies:    getSupportedCurrencies(),
	}, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 60
	}
	return i
}

// getSupportedCurrencies parses supported currencies from environment variable
func getSupportedCurrencies() []string {
	currenciesEnv := getEnv("SUPPORTED_CURRENCIES", "USD,EUR,GBP,JPY,CAD,AUD,CHF,CNY,SEK,NZD")
	
	// Split by comma and clean up whitespace
	currencies := strings.Split(currenciesEnv, ",")
	for i, currency := range currencies {
		currencies[i] = strings.TrimSpace(strings.ToUpper(currency))
	}
	
	// Filter out empty strings
	var result []string
	for _, currency := range currencies {
		if currency != "" {
			result = append(result, currency)
		}
	}
	
	// If no valid currencies found, return default set
	if len(result) == 0 {
		return []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"}
	}
	
	return result
}

