# Financial Forecasting Service Makefile

.PHONY: build run test clean deps lint

# Build the service
build:
	go build -o financial-forecasting-service main.go

# Run the service
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f financial-forecasting-service
	rm -f financial-forecasting-service.exe

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Run the service with environment file
run-env:
	cp env.example .env
	go run main.go

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o financial-forecasting-service-linux main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o financial-forecasting-service.exe main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o financial-forecasting-service-mac main.go

# Run all builds
build-all: build-linux build-windows build-mac

