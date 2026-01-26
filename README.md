# Hexa - Go Hexagonal Architecture API

A production-ready RESTful API built with Go using Hexagonal Architecture (Ports and Adapters) pattern.

## ✨ Features

- 🏗️ **Clean Architecture** - Hexagonal/Ports & Adapters pattern
- 🔐 **Secure Authentication** - PASETO tokens with refresh mechanism
- ✅ **Input Validation** - Struct-based validation with detailed error messages
- 🛡️ **Rate Limiting** - Protection against API abuse
- 📊 **Health Checks** - Database connectivity monitoring
- 🧪 **Testable** - Mock infrastructure for unit testing
- 🔄 **Graceful Shutdown** - Production-ready signal handling
- 📝 **Request Tracing** - Request ID for debugging
- 🐳 **Docker Support** - One-command development environment

## Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) - Fast HTTP framework
- **ORM**: [GORM](https://gorm.io/) - ORM for Go
- **Database**: PostgreSQL
- **Authentication**: [PASETO](https://github.com/o1egl/paseto) - Platform-Agnostic Security Tokens
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management
- **DI**: [Wire](https://github.com/google/wire) - Compile-time dependency injection
- **Validation**: [validator](https://github.com/go-playground/validator) - Struct validation
- **Logging**: [Zap](https://github.com/uber-go/zap) - Structured logging

## Project Structure

```
hexa/
├── cmd/
│   └── main.go                     # Application entry point
├── internal/
│   ├── adapters/
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   ├── auth_dto.go
│   │   │   └── mapper.go
│   │   ├── handlers/               # HTTP handlers (input adapters)
│   │   │   ├── auth_handler.go
│   │   │   ├── health_handler.go
│   │   │   └── response.go
│   │   ├── middleware/             # HTTP middleware
│   │   │   ├── auth_middleware.go
│   │   │   ├── rate_limiter.go
│   │   │   └── request_id.go
│   │   ├── repositories/           # Database repositories (output adapters)
│   │   │   ├── token_repository.go
│   │   │   └── user_repository.go
│   │   ├── routes/                 # Route definitions
│   │   │   └── routes.go
│   │   └── validator/              # Input validation
│   │       └── validator.go
│   ├── core/
│   │   ├── domain/                 # Domain entities
│   │   │   ├── errors.go
│   │   │   ├── token.go
│   │   │   └── user.go
│   │   ├── ports/                  # Interfaces (input/output ports)
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── services/               # Business logic
│   │       ├── auth_service.go
│   │       └── token_service.go
│   ├── dependencies/               # Wire DI
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── infrastructure/
│   │   ├── config/                 # Configuration
│   │   │   └── config.go
│   │   ├── database/               # Database connection & migrations
│   │   │   ├── database.go
│   │   │   ├── migrate.go
│   │   │   └── migrations/
│   │   ├── graceful/               # Graceful shutdown
│   │   │   └── shutdown.go
│   │   └── logger/                 # Structured logging
│   │       └── logger.go
│   └── testutil/                   # Testing utilities
│       └── mocks.go
├── docs/
│   ├── API.md                      # API documentation
│   └── IMPROVEMENTS.md             # Architecture improvements
├── .env                            # Environment variables
├── .env.example                    # Environment variables example
├── .gitignore
├── .golangci.yml                   # Linter configuration
├── docker-compose.yml              # Development environment
├── go.mod                          # Go modules
├── Makefile                        # Development commands
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24+ or Go 1.21+
- PostgreSQL 13+
- Make (optional, but recommended)

### Quick Start with Docker

The fastest way to get started:

```bash
# 1. Clone the repository
git clone <repository-url>
cd hexa

# 2. Copy environment file
cp .env.example .env

# 3. Start PostgreSQL
make docker-up

# 4. Run migrations
make migrate

# 5. Start the application
make dev
```

The API will be available at `http://localhost:8080`

### Manual Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd hexa
   ```

2. **Install dependencies**
   ```bash
   make deps
   make install-tools
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL**
   ```bash
   # Option 1: Using Docker
   make docker-up
   
   # Option 2: Use your local PostgreSQL
   # Update .env with your PostgreSQL credentials
   ```

5. **Run database migrations**
   ```bash
   make migrate
   ```

6. **Start the application**
   ```bash
   # Development mode (with auto-reload)
   make watch
   
   # Or standard development mode
   make dev
   
   # Or build and run
   make build
   make run
   ```

### Development Commands

```bash
make help              # Show all available commands
make dev              # Run in development mode
make watch            # Run with hot reload (Air)
make build            # Build the application
make test             # Run tests
make test-coverage    # Run tests with coverage report
make migrate          # Run database migrations
make migrate-down     # Rollback migrations
make migrate-create name=<name>  # Create new migration
make wire             # Generate dependency injection code
make fmt              # Format code
make lint             # Run linter
make vet              # Run go vet
make clean            # Clean build artifacts
make docker-up        # Start Docker services
make docker-down      # Stop Docker services
```

## API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication
Most endpoints require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

### Rate Limiting
- **Auth endpoints**: 5 requests per 15 minutes
- **API endpoints**: 100 requests per minute

### Common Response Format

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error description",
  "data": null
}
```

---

## API Endpoints

### Health Check

#### `GET /health`
Check API health and database connectivity

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "database": {
    "status": "healthy"
  },
  "uptime": "1h23m45s"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### Authentication Endpoints

#### `POST /auth/register`
Register a new user

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "password": "password123"
}
```

**Validation Rules:**
- `name`: Required, 2-100 characters
- `email`: Required, valid email format
- `phone`: Required, 10-15 characters
- `password`: Required, minimum 6 characters

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "created_at": "2026-01-15T10:00:00Z"
    },
    "access_token": "v2.local...",
    "refresh_token": "v2.local...",
    "expires_at": 1705320000
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "password123"
  }'
```

---

#### `POST /auth/login`
Login with email or phone and password

**Request Body:**
```json
{
  "username": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "created_at": "2026-01-15T10:00:00Z"
    },
    "access_token": "v2.local...",
    "refresh_token": "v2.local...",
    "expires_at": 1705320000
  }
}
```

**Example (with email):**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john@example.com",
    "password": "password123"
  }'
```

**Example (with phone):**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "+1234567890",
    "password": "password123"
  }'
```

---

#### `POST /auth/refresh`
Refresh access token using refresh token

**Request Body:**
```json
{
  "refresh_token": "v2.local..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "created_at": "2026-01-15T10:00:00Z"
    },
    "access_token": "v2.local...",
    "refresh_token": "v2.local...",
    "expires_at": 1705320000
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "v2.local.your_refresh_token_here"
  }'
```

---

#### `POST /auth/logout`
Logout and invalidate tokens (requires authentication)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### Error Codes

| Status Code | Description |
|------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input or validation error |
| 401 | Unauthorized - Invalid or missing authentication token |
| 404 | Not Found - Resource doesn't exist |
| 408 | Request Timeout - Request took too long |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Something went wrong |
| 503 | Service Unavailable - Service is down or degraded |

---

## Security

- ✅ **Password Hashing**: bcrypt with default cost (10)
- ✅ **Token Security**: PASETO v2 (Platform-Agnostic Security Tokens)
- ✅ **Token Expiration**: Access tokens (1 hour), Refresh tokens (7 days)
- ✅ **Rate Limiting**: Prevents brute-force attacks
- ✅ **Input Validation**: Comprehensive validation on all inputs
- ✅ **SQL Injection Protection**: GORM parameterized queries
- ✅ **CORS**: Configurable cross-origin resource sharing
- ✅ **Request Timeout**: 30 seconds default
- ✅ **Graceful Shutdown**: Prevents data corruption on restart

---

## Configuration

Configuration is managed via environment variables. Copy `.env.example` to `.env` and update:

```bash
# Application Configuration
APP_NAME=hexa-api
APP_PORT=8080
APP_ENV=development
APP_REQUEST_TIMEOUT=30s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hexa_db
DB_SSLMODE=disable
DB_OPT_MAX_IDLE_CONN=10
DB_OPT_MAX_OPEN_CONN=100
DB_OPT_CONN_MAX_IDLE_TIME=10m
DB_OPT_CONN_MAX_LIFE_TIME=1h

# Token Configuration
TOKEN_SECRET_KEY=your-secret-key-min-32-characters-long
TOKEN_EXPIRE_HOURS=1
TOKEN_REFRESH_EXPIRE_DAYS=7
```

---

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/core/services/...

# Run with race detector
go test -race ./...
```

---

## Database Migrations

### Create a new migration
```bash
make migrate-create name=add_users_table
```

This creates two files:
- `000001_add_users_table.up.sql` - Migration
- `000001_add_users_table.down.sql` - Rollback

### Run migrations
```bash
make migrate
```

### Rollback migrations
```bash
make migrate-down
```

---

## Hexagonal Architecture

This project follows the Hexagonal Architecture (Ports and Adapters) pattern:

### Core Concepts

```
┌─────────────────────────────────────────────────────┐
│                  External World                      │
│  (HTTP, CLI, Message Queue, Database, etc.)         │
└────────────┬──────────────────────────┬─────────────┘
             │                          │
        ┌────▼────┐              ┌─────▼──────┐
        │ Input   │              │  Output    │
        │ Adapters│              │  Adapters  │
        │(Handler)│              │(Repository)│
        └────┬────┘              └─────▲──────┘
             │                          │
        ┌────▼────┐              ┌─────┴──────┐
        │ Input   │              │  Output    │
        │  Ports  │              │   Ports    │
        │(Service │              │(Repository │
        │ Interface)              │ Interface) │
        └────┬────┘              └─────▲──────┘
             │                          │
             │    ┌──────────────┐     │
             └────►   CORE/      ◄─────┘
                  │   DOMAIN     │
                  │  (Business   │
                  │    Logic)    │
                  └──────────────┘
```

### Layer Responsibilities

- **Domain** (`internal/core/domain/`): Business entities, rules, and errors
- **Ports** (`internal/core/ports/`): Interfaces defining contracts
- **Services** (`internal/core/services/`): Business logic implementation
- **Adapters** (`internal/adapters/`): External interface implementations
  - **Handlers**: HTTP request/response (input adapters)
  - **Repositories**: Database operations (output adapters)
  - **DTOs**: Data transfer objects for API
  - **Middleware**: Cross-cutting concerns
- **Infrastructure** (`internal/infrastructure/`): Technical capabilities
  - Configuration, database connection, logging, etc.

### Benefits

✅ **Testability** - Core logic independent of frameworks
✅ **Maintainability** - Clear separation of concerns
✅ **Flexibility** - Easy to swap implementations
✅ **Independence** - Core doesn't depend on external tools

---

## Project Principles

### 1. Dependency Rule
Dependencies point inward. Core domain has no dependencies.

```
Infrastructure → Adapters → Ports → Domain
```

### 2. DTO Pattern
Never expose domain models directly in API responses.

**Domain (internal)**:
```go
type User struct {
    Password string `gorm:"not null"`  // Has sensitive data
}
```

**DTO (external)**:
```go
type UserResponse struct {
    // No password field - safe for API
}
```

### 3. Port/Adapter Pattern
Define contracts in ports, implement in adapters.

**Port**:
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
}
```

**Adapter**:
```go
type userRepository struct {
    db *gorm.DB
}
func (r *userRepository) Create(...) error {
    // GORM implementation
}
```

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards

- Run `make fmt` before committing
- Run `make lint` to check code quality
- Run `make test` to ensure tests pass
- Add tests for new features
- Follow existing code structure

---

## Documentation

- **[API Documentation](docs/API.md)** - Detailed API reference
- **[Architecture Improvements](docs/IMPROVEMENTS.md)** - Design decisions and enhancements

---

## Troubleshooting

### Database connection errors
```bash
# Check if PostgreSQL is running
make docker-up

# Check logs
make docker-logs

# Verify connection in .env file
```

### Migration errors
```bash
# Check migration files
ls -la internal/infrastructure/database/migrations/

# Force migration version
# (use with caution)
make migrate-down
make migrate
```

### Port already in use
```bash
# Change APP_PORT in .env
# Or kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

---

## Production Deployment

### Build

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/hexa cmd/main.go
```

### Environment Variables

Ensure these are set in production:
- `APP_ENV=production`
- `TOKEN_SECRET_KEY` - Strong, unique secret (32+ characters)
- `DB_*` - Production database credentials
- Enable `DB_SSLMODE=require` for PostgreSQL

### Recommended Setup

- Use a process manager (systemd, supervisord)
- Set up reverse proxy (Nginx, Traefik)
- Enable HTTPS
- Configure monitoring (Prometheus, Grafana)
- Set up log aggregation
- Use managed database service
- Enable database backups
- Set up health check endpoints for load balancer

---

## Performance Considerations

- **Connection Pooling**: Configured via `DB_OPT_MAX_OPEN_CONN` and `DB_OPT_MAX_IDLE_CONN`
- **Request Timeout**: Default 30s, configurable via `APP_REQUEST_TIMEOUT`
- **Rate Limiting**: In-memory, consider Redis for distributed systems
- **Database Indexes**: Review `internal/infrastructure/database/migrations/`

---

## License

MIT License - see LICENSE file for details

---

## Support

For issues and questions:
- Create an issue in the repository
- Check existing documentation in `docs/`
- Review the codebase examples

---

**Built with ❤️ using Hexagonal Architecture principles**
