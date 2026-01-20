# Xonvera Core API Documentation

This directory contains API documentation for the Xonvera Core service.

## Files

- **openapi.yaml** - OpenAPI 3.0 specification for the API
- **postman_collection.json** - Postman collection for testing the API

## OpenAPI Specification

The `openapi.yaml` file contains the complete API specification in OpenAPI 3.0 format. You can:

- View it in [Swagger Editor](https://editor.swagger.io/)
- Generate client SDKs using [OpenAPI Generator](https://openapi-generator.tech/)
- Import it into API documentation tools like [Redoc](https://redocly.com/) or [Stoplight](https://stoplight.io/)

### Viewing the OpenAPI Spec

1. Visit [Swagger Editor](https://editor.swagger.io/)
2. File → Import File
3. Select `openapi.yaml`

## Postman Collection

The `postman_collection.json` file contains a complete Postman collection for testing the API.

### Importing into Postman

1. Open Postman
2. Click "Import" button
3. Select `postman_collection.json`
4. The collection will be imported with all endpoints configured

### Setting up Environment Variables

Create a new environment in Postman with the following variables:

| Variable | Description | Example Value |
|----------|-------------|---------------|
| `baseUrl` | Base URL of the API | `http://localhost:3000` |
| `accessToken` | JWT access token (auto-set after login) | - |
| `refreshToken` | JWT refresh token (auto-set after login) | - |

The collection includes test scripts that automatically save tokens to environment variables after successful login/register/refresh operations.

## API Endpoints

### Health
- `GET /health` - Health check

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout user (requires authentication)

### Protected Routes
- `GET /protected/profile` - Get user profile (requires authentication)

## Authentication

Most endpoints require JWT authentication. Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Token Flow

1. **Register** or **Login** to receive an access token and refresh token
2. Use the **access token** for authenticated requests
3. When the access token expires, use the **refresh token** to get a new access token
4. **Logout** invalidates the current access token

## Request Examples

### Register

```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "password": "SecurePass123"
  }'
```

### Login

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john.doe@example.com",
    "password": "SecurePass123"
  }'
```

### Get Profile (Protected)

```bash
curl -X GET http://localhost:3000/protected/profile \
  -H "Authorization: Bearer <access_token>"
```

## Response Format

All responses follow a standard format:

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "data": null
}
```

## Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required or failed
- `408 Request Timeout` - Request took too long
- `503 Service Unavailable` - Service is unhealthy

## Validation Rules

### Register
- **name**: Required, 2-100 characters
- **email**: Required, valid email format
- **phone**: Required, 10-15 characters
- **password**: Required, minimum 6 characters

### Login
- **username**: Required (email or phone)
- **password**: Required

### Refresh Token
- **refresh_token**: Required

## Development

To start the API server locally:

```bash
# From the project root
make run
```

The API will be available at `http://localhost:3000`

## Support

For API support, please contact the development team.
