.PHONY: build run dev watch wire migrate migrate-down migrate-create test clean help lint fmt vet install-tools

# Application
APP_NAME=xonvera
MAIN_FILE=cmd/main.go

## help: Display this help message
help:
	@echo "Available commands:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/^## //'

## build: Build the application binary
build:
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

## run: Build and run the application
run: build
	./bin/$(APP_NAME)

## dev: Run in development mode
dev:
	go run $(MAIN_FILE)

## watch: Run with Air (live reload)
watch:
	air

## wire: Generate Wire dependencies
wire:
	cd internal/dependencies && wire

## migrate: Run database migrations
migrate:
	go run $(MAIN_FILE) -migrate

## migrate-down: Rollback database migrations by one step
migrate-down:
	go run $(MAIN_FILE) -migrate-down

## migrate-reset: Reset database migrations
migrate-reset:
	go run $(MAIN_FILE) -migrate-reset

## migrate-create name=<name>: Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@MIGRATION_PATH=internal/infrastructure/database/migrations; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch $$MIGRATION_PATH/$${TIMESTAMP}_$(name).up.sql; \
	touch $$MIGRATION_PATH/$${TIMESTAMP}_$(name).down.sql; \
	echo "Created migration files:"; \
	echo "  $$MIGRATION_PATH/$${TIMESTAMP}_$(name).up.sql"; \
	echo "  $$MIGRATION_PATH/$${TIMESTAMP}_$(name).down.sql"

## test: Run all tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts and cache
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

## deps: Install project dependencies and tools
deps:
	go mod download
	go mod tidy

## install-tools: Install development tools
install-tools:
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## fmt: Format code using gofmt
fmt:
	go fmt ./...
	gofmt -s -w .

## vet: Run go vet for static analysis
vet:
	go vet ./...

## lint: Run golangci-lint
lint:
	golangci-lint run --timeout=5m

## docker-up: Start Docker containers
docker-up:
	docker-compose up -d

## docker-down: Stop Docker containers
docker-down:
	docker-compose down

## docker-logs: View Docker container logs
docker-logs:
	docker-compose logs -f
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Install dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"
