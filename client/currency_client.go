package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	currencymodels "github.com/dalfonso89/currency-exchange-service/models"
	"github.com/dalfonso89/financial-forecasting-service/config"
	"github.com/dalfonso89/financial-forecasting-service/logger"
)

// CurrencyClient handles communication with the currency exchange service
type CurrencyClient struct {
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
}

// NewCurrencyClient creates a new currency client
func NewCurrencyClient(cfg *config.Config, logger logger.Logger) *CurrencyClient {
	return &CurrencyClient{
		baseURL: cfg.CurrencyExchangeServiceURL,
		httpClient: &http.Client{
			Timeout: cfg.CurrencyExchangeTimeout,
		},
		logger: logger,
	}
}

// GetRates fetches exchange rates from the currency exchange service
func (c *CurrencyClient) GetRates(ctx context.Context, baseCurrency string) (*currencymodels.RatesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/rates/%s", c.baseURL, baseCurrency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.logger.Debugf("Fetching rates from: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("currency service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ratesResponse currencymodels.RatesResponse
	if err := json.Unmarshal(body, &ratesResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Debugf("Successfully fetched rates for base currency: %s", baseCurrency)
	return &ratesResponse, nil
}

// GetRatesWithQuery fetches exchange rates using query parameters
func (c *CurrencyClient) GetRatesWithQuery(ctx context.Context, baseCurrency string) (*currencymodels.RatesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/rates?base=%s", c.baseURL, baseCurrency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.logger.Debugf("Fetching rates from: %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("currency service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var ratesResponse currencymodels.RatesResponse
	if err := json.Unmarshal(body, &ratesResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Debugf("Successfully fetched rates for base currency: %s", baseCurrency)
	return &ratesResponse, nil
}

// HealthCheck checks if the currency exchange service is healthy
func (c *CurrencyClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("currency service health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
