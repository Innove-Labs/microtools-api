# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

microtools-go is a Go-based HTTP API service providing utility microservices including email validation, IP geolocation, IBAN validation, QR code generation, and barcode generation. The service uses MongoDB for persistence, CounterAPI.dev for API hit tracking, and serves on port 8000. The project follows the **Go Standard Project Layout** with clean architecture principles.

## Build & Run Commands

```bash
# Navigate to source directory
cd /home/steinsgate/main/innovelabs/projects/microtools-go/microtools

# Run the server (requires .env file)
go run cmd/api/main.go

# Build binary
go build -o bin/api cmd/api/main.go

# Install dependencies
go mod download

# Tidy dependencies
go mod tidy

# Docker build
docker build -t fawazsullialabs/innovelabs-micro-apis:0.0.1 .
```

## Environment Setup

Create a `.env` file in the `microtools/` directory with:
- `MONGO_URI` - MongoDB connection string
- `REDIS_URI` - Redis connection string
- `JWT_SECRET` - Secret key for JWT signing
- `COUNTER_API_KEY` - API key for CounterAPI.dev hit tracking

The GeoIP2 database file `geolite-2-city.mmdb` is located in the `assets/` directory for IP geolocation functionality.

## Architecture

The project follows **Go Standard Project Layout** with clean separation of concerns:

```
microtools/
├── cmd/api/              # Application entry point
│   └── main.go          # Main application setup and initialization
├── internal/            # Private application code
│   ├── config/         # Configuration management
│   ├── models/         # Data models and DTOs
│   ├── services/       # Business logic layer
│   │   ├── validation/ # Validation services (email, IP, IBAN)
│   │   └── generator/  # Generation services (QR, barcode)
│   ├── handlers/       # HTTP handlers (presentation layer)
│   ├── middleware/     # HTTP middleware
│   ├── router/         # Route configuration
│   ├── database/       # Database connections
│   └── utils/          # Utility functions
├── web/                # Web assets
│   └── templates/      # HTML templates
│       ├── layout/     # Layout templates
│       └── pages/      # Page templates
├── assets/             # Static assets (GeoIP database, etc.)
└── go.mod              # Module definition
```

### Layer Responsibilities

**cmd/api/main.go**: Application entry point
- Loads configuration
- Initializes database connections
- Sets up router and middleware
- Starts HTTP server

**internal/config**: Configuration management
- Environment variable loading via godotenv
- Config struct with application settings

**internal/models**: Data Transfer Objects (DTOs)
- Request models (EmailRequest, IPRequest, IBANRequest, QRRequest, etc.)
- Response models (EmailValidation, GeoIPResponse, IBANValidation, etc.)
- IBAN country specifications (60+ countries)

**internal/services**: Business logic layer
- `validation/email.go` - Email validation with syntax, domain, MX record checks, and disposable email detection
- `validation/ip.go` - IP geolocation using MaxMind GeoIP2 database
- `validation/iban.go` - IBAN validation with mod-97 checksum verification for 60+ countries
- `generator/qr.go` - QR code generation supporting 10 types (text, URL, email, WiFi, vCard, etc.)
- `generator/barcode.go` - 1D barcode generation (UPC-A, EAN-13, Code128) with PNG/SVG output

**internal/handlers**: HTTP layer
- Decodes JSON requests
- Calls appropriate service functions
- Encodes JSON responses
- Handles errors and status codes

**internal/middleware**: HTTP middleware
- `auth.go` - JWT authentication middleware
- `counter.go` - API counter middleware using CounterAPI.dev

**internal/router**: Route configuration
- Sets up gorilla/mux router
- Registers all endpoints with handlers
- Applies middleware

**internal/database**: Database connections
- `mongo.go` - MongoDB client initialization

**internal/utils**: Utility functions
- JWT token generation and validation

### HTTP Router
Uses gorilla/mux with these endpoints:
- `POST /api/v1/validate/email` - Email validation
- `POST /api/v1/validate/ip` - IP geolocation lookup
- `POST /api/v1/validate/iban` - IBAN validation
- `POST /api/v1/generate/qr` - QR code generation (returns PNG)
- `POST /api/v1/generate/barcode` - 1D barcode generation (returns PNG or SVG)
- `GET /api/v1/live` - Health check
- `GET /` - Home page with API documentation
- `GET /email-validation-api` - Email validation API page
- `GET /ip-geolocation-api` - IP geolocation API page
- `GET /iban-validation-api` - IBAN validation API page
- `GET /qr-code-generator-api` - QR code generator API page
- `GET /barcode-generator-api` - Barcode generator API page

### Active Middleware
- **APICounterMiddleware**: Applied globally via `router.Use()`. Fires a background HTTP call to CounterAPI.dev to increment per-endpoint counters. Non-blocking — the response is served before the counter call completes.

### Deployment
- **Dockerfile**: Multi-stage build using `golang:1.23` builder and `alpine:latest` runtime. Builds a static binary (`CGO_ENABLED=0`) and exposes port 8000.

## Key Implementation Details

### Email Validation (`internal/services/validation/email.go`)
The `ValidateEmail()` function returns a structured result with four checks:
- Syntax validation (regex-based)
- Domain validity (MX or A records)
- MX records presence
- Disposable email detection (against hardcoded list of 14 providers)

SMTP verification code exists but is commented out due to anti-spam policies blocking verification attempts.

### IP Geolocation (`internal/services/validation/ip.go`)
Uses the MaxMind GeoIP2 City database file located in `assets/geolite-2-city.mmdb`. Returns country, region, city, coordinates, and timezone for valid IPs.

### IBAN Validation (`internal/services/validation/iban.go`)
Comprehensive International Bank Account Number validation supporting 60+ countries:
- Country code validation
- Length validation per country
- BBAN format validation using regex patterns
- Mod-97 checksum verification (ISO 13616 standard)
- Returns detailed breakdown: country, bank code, account number, check digits, formatted IBAN
- Supports SEPA countries, Middle East, Latin America, and other regions

Country specifications are defined in `internal/models/iban.go`.

### QR Code Generation (`internal/services/generator/qr.go`)
Supports 10 types: text, url, email, tel, sms, wifi, vcard, geo, event, json
- Type-specific payload formatting (e.g., WIFI:, VCARD:, VEVENT:)
- Configurable size (128-1024px) and error correction (low/medium/high/highest)
- JSON input for structured types (wifi, vcard, event)

### Barcode Generation (`internal/services/generator/barcode.go`)
1D barcode generation with interface-based dependency injection:
- Supports UPC-A, EAN-13, Code128
- PNG and SVG output formats
- Customizable colors, dimensions, text placement
- Clean architecture with BarcodeService interface

## Working with This Codebase

### Code Organization
- All application code is in the `internal/` package (not accessible outside this module)
- Business logic is isolated in `services/` layer - completely independent of HTTP
- Handlers in `handlers/` deal only with HTTP concerns (request/response marshaling)
- Models in `models/` define all data structures shared across layers
- Configuration is centralized in `internal/config/`

### Adding New Features
1. Define request/response models in `internal/models/`
2. Implement business logic in `internal/services/`
3. Create HTTP handler in `internal/handlers/`
4. Register route in `internal/router/`

### Template System
- Templates are in `web/templates/` directory
- Layout files in `web/templates/layout/`
- Page-specific templates in `web/templates/pages/`
- The `web/` directory must be accessible relative to the executable

### Error Handling
- All handler functions follow the pattern: decode JSON → validate → call service → encode response
- Error responses use standard HTTP status codes with JSON error messages
- Service layer returns errors, handlers translate them to HTTP responses

### Module Information
- Module name: `github.com/innovelabs/microtools-go`
- Import paths must use full module path (e.g., `github.com/innovelabs/microtools-go/internal/models`)
- The `.git` directory lives inside `microtools/`, not at the project root

### Legacy Files
Old flat-structure files may still exist in the root directory (e.g., `email-handling.go`, `geolocation-handling.go`, etc.). These should be ignored as all functionality has been migrated to the new structure under `internal/`.
