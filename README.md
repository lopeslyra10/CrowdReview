# CrowdReview Backend (Go + Gin)

Backend reference implementation for CrowdReview, following a lightweight Clean Architecture layout with Gin, GORM (PostgreSQL), Redis (rate limit + background queue), and JWT auth.

## Stack
- Go 1.22, Gin, GORM (PostgreSQL)
- Redis for caching, rate limit, and async fraud validation queue
- JWT + Refresh tokens
- Swagger (via swaggo) entrypoint on `/swagger/*any` (docs generation not run here)

## Project Structure
```
/cmd/api              # main entrypoint & wiring
/config               # config loading
/internal/handlers    # HTTP handlers
/internal/services    # business services
/internal/repository  # data access (GORM)
/internal/models      # domain models
/internal/validation  # fraud engine & worker
/internal/rules       # fraud heuristics
/pkg/middleware       # shared middleware
/pkg/utils            # helpers (jwt, password, responses)
```

### Architectural Notes
- Gin handles transport-only concerns in `/internal/handlers`; no business logic.
- Services in `/internal/services` express use-cases; they depend on repository interfaces, not GORM specifics.
- Repositories in `/internal/repository` wrap GORM and expose interfaces injected into services (clean architecture dependency flow).
- Fraud validation is isolated in `/internal/validation` (engine + worker) and `/internal/rules` (heuristics).
- Middleware under `/pkg/middleware` stays framework-agnostic where possible (auth, admin, rate limiting, logging).

## Running locally
1) Set environment variables (or `.env`):
```
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/crowdreview?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=super-secret
REFRESH_SECRET=another-secret
TOKEN_TTL_MINUTES=30
REFRESH_TTL_HOURS=24
RATE_LIMIT_REQUESTS=20
RATE_LIMIT_WINDOW=60
```
2) Start deps with docker-compose:
```
docker-compose up -d
```
3) Run the API:
```
go run ./cmd/api
```

## Swagger
- `swag init -g cmd/api/main.go` to generate docs (requires swag CLI).
- Served at `/swagger/*any` when `docs` package is generated.

## Tests
```
go test ./...
```

## Notes
- Background fraud validation uses a Go channel queue (`fraud-validation-queue`) and persists `ReviewValidationResult`.
- Middleware: AuthRequired, AdminRequired, RateLimitMiddleware, RequestLogger.
- Admin routes under `/admin/*` require `role=admin`.
