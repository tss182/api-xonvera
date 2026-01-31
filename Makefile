.PHONY: build run dev watch wire migrate migrate-down migrate-create test clean help lint fmt vet install-tools swag docker-up docker-down docker-logs update

# Application
APP_NAME=xonvera
MAIN_FILE=cmd/main.go

# Colors
COLOR_GREEN=\033[0;32m
COLOR_YELLOW=\033[1;33m
COLOR_CYAN=\033[0;36m
COLOR_RESET=\033[0m

# === 📦 BUILD COMMANDS ===

## help: Display all available commands with colors and groups
help:
	@printf "\n"
	@awk -v green="$(COLOR_GREEN)" -v yellow="$(COLOR_YELLOW)" -v cyan="$(COLOR_CYAN)" -v reset="$(COLOR_RESET)" '\
	/^# === / { \
		if (group != "") print ""; \
		group=$$0; \
		gsub(/^# === |===.*/, "", group); \
		printf "%s%s%s\n", green, group, reset; \
	} \
	/^## / && !/help:/ { \
		cmd=$$2; \
		sub(/: /, "", cmd); \
		desc=$$0; \
		gsub(/^## [^ ]* ?/, "", desc); \
		gsub(/^[^:]*: /, "", desc); \
		printf "  %s%-25s%s %s%s%s\n", yellow, cmd, reset, cyan, desc, reset; \
	}' $(MAKEFILE_LIST)

## build: Build the application binary
build:
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

# === 🚀 DEVELOPMENT COMMANDS ===

## run: Run the application
run:
	go run $(MAIN_FILE)

## dev: Run with reload
dev:
	air

## wire: Generate Wire dependencies
wire:
	cd internal/dependencies && wire

# === 🗄️ DATABASE COMMANDS ===

## migration: Run database migrations
migration:
	go run $(MAIN_FILE) -migrate

## migration-down: Rollback database migrations by one step
migration-down:
	go run $(MAIN_FILE) -migrate-down

## migration-reset: Reset database migrations
migration-reset:
	go run $(MAIN_FILE) -migrate-reset

## migration-create name=<name>: Create a new migration
migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration-create name=migration_name"; \
		exit 1; \
	fi
	@MIGRATION_PATH=internal/infrastructure/database/migrations; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch $$MIGRATION_PATH/$${TIMESTAMP}_$(name).up.sql; \
	touch $$MIGRATION_PATH/$${TIMESTAMP}_$(name).down.sql; \
	echo "Created migration files:"; \
	echo "  $$MIGRATION_PATH/$${TIMESTAMP}_$(name).up.sql"; \
	echo "  $$MIGRATION_PATH/$${TIMESTAMP}_$(name).down.sql"

# === ✅ TESTING & QUALITY ===

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

# === 🔧 UTILITY COMMANDS ===

## deps: Install project dependencies and tools
deps:
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## install-tools: Install development tools (alias for deps)
install-tools: deps

## update: Update all dependencies to latest versions
update:
	go get -u  ./...
	go mod tidy

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

# === 📚 DOCUMENTATION & API ===

## swag: Generate API documentation using swag
swag:
	swag init -g cmd/main.go --output docs/swagger

# === 🐳 DOCKER COMMANDS ===

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
