# CrowdReview Backend (Go + Gin)

Implementação de backend para o CrowdReview, seguindo um layout leve de Clean Architecture com Gin, GORM (PostgreSQL), Redis (rate-limit + fila de tarefas), e autenticação JWT.

## Stack
- Go 1.22, Gin, GORM (PostgreSQL)
- Redis para cache, rate-limit e fila assíncrona de validação de fraude
- JWT + Tokens de Refresh
- Swagger (via swaggo) disponível em /swagger/*any (geração de docs não incluída aqui)

## Estrutura do Projeto
```
/cmd/api              # ponto de entrada e wiring
/config               # carregamento de configuração
/internal/handlers    # HTTP handlers
/internal/services    # business services
/internal/repository  # data access (GORM)
/internal/models      # domain models
/internal/validation  # fraud engine & worker
/internal/rules       # fraud heuristics
/pkg/middleware       # shared middleware
/pkg/utils            # helpers (jwt, password, responses)
```

### Notas de Arquitetura
- O Gin lida apenas com transporte em `/internal/handlers`; nenhuma regra de negócio fica aqui..
- Repositórios em  `/internal/repository` encapsulam o GORM e expõem interfaces que são injetadas nos services (fluxo de dependência do Clean Architecture).
- Os serviços em `/internal/services` representam os use-cases; eles dependem de interfaces de repositório, não de implementações específicas do GORM.
- A validação de fraude fica isolada em `/internal/validation` (engine + worker) e /internal/rules (heurísticas).
- Middleware em `/pkg/middleware` permanece o mais independente possível do framework (auth, admin, rate-limit, logging).

## Running locally
1) Configure as variáveis de ambiente (ou arquivo `.env`):
```bash
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
2) Suba as dependências com docker-compose:
```
docker-compose up -d
```
3) Execute a API:
```
go run ./cmd/api
```

## Swagger
- Use `swag init -g cmd/api/main.go` para gerar a documentação (requer o CLI do swag).
- Servido em `/swagger/*any` quando o pacote `docs` é gerado.

## Testes
```bash
go test ./...
```

## Notas
- A validação de fraude em background usa uma fila implementada com Go channels (`fraud-validation-queue`) e persiste `ReviewValidationResult`.
- Middleware disponíveis: AuthRequired, AdminRequired, RateLimitMiddleware, RequestLogger.
- Rotas de admin em `/admin/*` exigem `role=admin`.
