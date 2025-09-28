# Financial Forecasting Service - Usage Examples

This document provides examples of how to use the Financial Forecasting Service API.

## Prerequisites

1. Start the Currency Exchange Service on port 8081
2. Start the Financial Forecasting Service on port 8082

## Example API Calls

### 1. Health Check

```bash
curl http://localhost:8082/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "uptime": "5m30s"
}
```

### 2. Get Supported Currencies

```bash
curl http://localhost:8082/api/v1/currencies
```

Response:
```json
{
  "currencies": ["USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD"]
}
```

### 3. Generate Single Currency Forecast

```bash
curl -X POST http://localhost:8082/api/v1/forecast \
  -H "Content-Type: application/json" \
  -d '{
    "base_currency": "USD",
    "target_currency": "EUR",
    "amount": 1000,
    "periods": 30,
    "forecast_type": "linear"
  }'
```

Response:
```json
{
  "base_currency": "USD",
  "target_currency": "EUR",
  "current_rate": 0.85,
  "amount": 1000,
  "forecast_type": "linear",
  "periods": 30,
  "forecasts": [
    {
      "period": 1,
      "date": "2024-01-16",
      "rate": 0.851,
      "amount": 851.0,
      "change": 0.001,
      "change_percent": 0.12
    },
    {
      "period": 2,
      "date": "2024-01-17",
      "rate": 0.852,
      "amount": 852.0,
      "change": 0.001,
      "change_percent": 0.12
    }
  ],
  "generated_at": "2024-01-15T10:30:00Z",
  "confidence_score": 0.7
}
```

### 4. Generate Multi-Currency Forecast

```bash
curl -X POST http://localhost:8082/api/v1/forecast/multi-currency \
  -H "Content-Type: application/json" \
  -d '{
    "base_currency": "USD",
    "currencies": ["EUR", "GBP", "JPY"],
    "amount": 1000,
    "periods": 7,
    "forecast_type": "exponential"
  }'
```

### 5. Analyze Currency Trend

```bash
curl "http://localhost:8082/api/v1/forecast/trend/USD/EUR?periods=30"
```

Response:
```json
{
  "currency_pair": "USD/EUR",
  "trend": "sideways",
  "volatility": 0.05,
  "average_rate": 0.85,
  "min_rate": 0.8075,
  "max_rate": 0.8925,
  "analysis_period": 30,
  "generated_at": "2024-01-15T10:30:00Z"
}
```

### 6. Clear Forecast Cache

```bash
curl -X DELETE http://localhost:8082/api/v1/forecast/cache
```

Response:
```json
{
  "message": "Cache cleared successfully"
}
```

## Forecasting Types

### Linear Forecasting
- Simple linear trend forecasting
- Assumes constant rate of change
- Good for short-term predictions

### Exponential Forecasting
- Exponential growth/decay forecasting
- Assumes compounding rate of change
- Good for long-term trend analysis

### Moving Average Forecasting
- Moving average with volatility
- Includes random-like variations
- Good for volatile markets

## Error Handling

The service returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request parameters
- `500 Internal Server Error`: Service errors
- `503 Service Unavailable`: Currency exchange service unavailable

Example error response:
```json
{
  "error": "invalid request",
  "message": "base currency is required",
  "code": 400
}
```

## Integration Notes

- The service automatically fetches current exchange rates from the currency exchange service
- Results are cached for improved performance
- The service supports graceful degradation if the currency service is unavailable
- All timestamps are in UTC format

