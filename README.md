# Financial Forecasting Service

A Golang microservice that provides financial forecasting capabilities by consuming the currency exchange service. This service offers various forecasting algorithms and trend analysis for currency pairs.

## Features

- **Currency Forecasting**: Generate forecasts for currency exchange rates using multiple algorithms
- **Multi-Currency Support**: Forecast multiple currencies simultaneously
- **Trend Analysis**: Analyze trends and volatility for currency pairs
- **Caching**: Built-in caching for improved performance
- **RESTful API**: Clean REST API with comprehensive endpoints
- **Health Checks**: Built-in health monitoring
- **Configurable**: Environment-based configuration

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Forecasting
- `POST /api/v1/forecast` - Generate single currency forecast
- `POST /api/v1/forecast/multi-currency` - Generate multi-currency forecast
- `GET /api/v1/forecast/trend/:base/:target` - Analyze currency trend
- `DELETE /api/v1/forecast/cache` - Clear forecast cache

### Currency Information
- `GET /api/v1/currencies` - Get supported currencies
- `GET /api/v1/currencies/rates/:base` - Get current exchange rates

## Configuration

The service can be configured using environment variables. Copy `env.example` to `.env` and modify as needed:

```bash
cp env.example .env
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8082 | Service port |
| `LOG_LEVEL` | info | Logging level |
| `CURRENCY_EXCHANGE_SERVICE_URL` | http://localhost:8081 | Currency exchange service URL |
| `CURRENCY_EXCHANGE_TIMEOUT_SECONDS` | 30 | Timeout for currency service calls |
| `FORECAST_CACHE_TTL_SECONDS` | 300 | Forecast cache TTL in seconds |
| `MAX_CONCURRENT_REQUESTS` | 10 | Maximum concurrent requests |
| `DEFAULT_FORECAST_PERIODS` | 30 | Default number of forecast periods |

## Usage

### Building and Running

```bash
# Install dependencies
make deps

# Build the service
make build

# Run the service
make run

# Run with environment file
make run-env
```

### Example API Calls

#### Generate Single Currency Forecast

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

#### Generate Multi-Currency Forecast

```bash
curl -X POST http://localhost:8082/api/v1/forecast/multi-currency \
  -H "Content-Type: application/json" \
  -d '{
    "base_currency": "USD",
    "currencies": ["EUR", "GBP", "JPY"],
    "amount": 1000,
    "periods": 30,
    "forecast_type": "exponential"
  }'
```

#### Analyze Trend

```bash
curl http://localhost:8082/api/v1/forecast/trend/USD/EUR?periods=30
```

## Forecasting Types

The service supports three forecasting algorithms:

1. **Linear**: Simple linear trend forecasting
2. **Exponential**: Exponential growth/decay forecasting
3. **Moving Average**: Moving average with volatility

## Architecture

```
┌─────────────────────┐    HTTP    ┌──────────────────────┐
│ Financial           │ ──────────► │ Currency Exchange    │
│ Forecasting Service │            │ Service              │
└─────────────────────┘            └──────────────────────┘
         │
         │ REST API
         ▼
┌─────────────────────┐
│ Client Applications │
└─────────────────────┘
```

## Dependencies

- **Gin**: HTTP web framework
- **Logrus**: Structured logging
- **Godotenv**: Environment variable loading

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Building for Different Platforms

```bash
# Build for Linux
make build-linux

# Build for Windows
make build-windows

# Build for macOS
make build-mac

# Build for all platforms
make build-all
```

## Integration with Currency Exchange Service

This service is designed to work with the currency exchange service and **imports the `RatesResponse` struct directly** from the currency-exchange-service module. Make sure the currency exchange service is running on the configured URL before starting this service.

The financial forecasting service will:
1. Fetch current exchange rates from the currency service using the shared `RatesResponse` struct
2. Apply forecasting algorithms to predict future rates
3. Return structured forecast data with confidence scores
4. Cache results for improved performance

### Shared Data Models

The service now uses the same data structures as the currency exchange service:
- `RatesResponse` struct is imported directly from `github.com/dalfonso89/currency-exchange-service/models`
- This ensures consistency between services and eliminates duplicate type definitions

## Error Handling

The service includes comprehensive error handling:
- Input validation
- Service communication errors
- Graceful degradation
- Structured error responses

## Monitoring

- Health check endpoint for service monitoring
- Request logging with correlation IDs
- Performance metrics through logging
- Cache statistics
