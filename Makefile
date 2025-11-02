run:
	@echo "Starting server..."
	go run ./cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v ./...

test-handler:
	@echo "Running handler tests..."
	go test -v ./internal/handler/

test-service:
	@echo "Running service tests..."
	go test -v ./internal/service/
	
mock:
	@echo "Generating mocks..."
	mockery

build:
	@echo "Building binary..."
	go build -o build/expense-tracker ./cmd/api
