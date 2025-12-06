# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Machine Marketplace is a monorepo application for renting/buying computational resources. Users can list machines with SSH access and rent them through a marketplace interface. Purchase events are sent to Kafka for async processing.

**Tech Stack:**
- Frontend: React 18 + TypeScript + Vite + Material-UI
- Backend: Go 1.24 + Echo framework + PostgreSQL 15
- Authentication: JWT (HTTP-only cookies)
- Message Queue: Kafka (for purchase event processing)

## Development Commands

### Running Services Locally (Docker)
```bash
docker-compose up          # Start all services (client:80, server:3001, postgres:5432, kafka:9092, zookeeper:2181)
docker-compose down        # Stop all services
```

### Backend (Go) - from /server directory
```bash
# Build
make build                 # Build binary to build/machine-marketplace
make build-prod           # Production build (linux/amd64)

# Run individual services (local PostgreSQL)
make run-auth             # Auth service on port 3001
make run-order            # Order service on port 3002
make run-process          # Process service on port 3003 (stub)

# Development
make dev                  # Auto-reload with air (requires: go install github.com/cosmtrek/air@latest)

# Testing & Quality
make test                 # Run all tests
make test-coverage        # Generate coverage.html report
make fmt                  # Format Go code
make vet                  # Run go vet
make deps                 # Download and tidy dependencies

# Database configuration for local runs (edit these in Makefile if needed)
# DB_HOST=localhost, DB_USER=boazfrid, DB_PASSWORD="", DB_NAME=machine_market
```

### Frontend (React) - from /client directory
```bash
npm run dev               # Start Vite dev server
npm run build             # Build production bundle (TypeScript + Vite)
npm run lint              # Run ESLint
npm run preview           # Preview production build
```

## Architecture

### Multi-Service Backend

The backend runs as separate service processes, each started via command-line arguments:

```bash
./machine-marketplace auth-service --db-host=localhost --db-user=postgres --db-password=secret
./machine-marketplace order-service --db-host=localhost ...
./machine-marketplace process-service --db-host=localhost ...
```

**Service Separation:**
- **Auth Service** (port 3001): User signup, login, logout, JWT issuance
- **Order Service** (port 3002): Machine listing, creation, purchasing (sends to Kafka)
- **Process Service** (port 3003): Stub for background job processing (not implemented)
- **Kafka** (port 9092): Message queue for purchase events

Entry flow: `main.go` → `cmd/root.go` (parses service type) → `cmd/{auth,order,process}.go` → respective modules

### Module Pattern

Each service is a `Module` struct:
```go
type Module struct {
    P  string         // Port
    DB *db.Queries    // Database queries (sqlc-generated)
    E  *echo.Echo     // HTTP router
}
```

Services initialize with: `database.InitWithConfig()` → `module.SetupRoutes()` → `echo.Start()`

### Database (PostgreSQL)

**Type-safe queries via sqlc:**
- Schema: `internal/DB/schema.sql`
- Queries: `internal/DB/query.sql`
- Generated: `internal/DB/generated/` (models.go, query.sql.go)

**Key tables:**
- `users`: id, name (unique), email (unique), password (bcrypt)
- `machines`: id, name, owner_id, buyer_id (NULL = available), ram, cpu, memory, host, ssh_user, key

**Initialization:**
- Database schema and seed data run automatically on service startup
- Seed data includes test users: `boaz@test.com`, `noam@test.com`

### Authentication Flow

1. **Login** → Auth service validates credentials → generates JWT (24h expiry) → stores in HTTP-only cookie
2. **Protected routes** → Middleware extracts JWT → validates signature → adds claims to context
3. **Authorization** → User ID from `claims.Issuer` field → filter queries by owner_id/buyer_id

**Security notes:**
- JWT secret is hardcoded as `"secret"` in `pkg/auth/auth.go` (should use env var)
- SSH host key verification is disabled (InsecureIgnoreHostKey)
- No rate limiting on auth endpoints

### Kafka Purchase Events

**Topic:** `machine-purchases`

**Purchase Flow:**
1. Client sends POST to `/api/v1/order/buy` with machine_id and deal_duration_hours
2. Order service updates database (sets buyer_id on machine)
3. Order service sends purchase event to Kafka

**Event Schema:**
```go
{
    "machine_id": int32,
    "buyer_id": int32,
    "deal_expiration": timestamp,
    "purchase_time": timestamp,
    "machine_name": string
}
```

**Kafka Configuration:**
- Broker: `kafka:29092` (Docker) or `localhost:9092` (local)
- Environment variable: `KAFKA_BROKER`

### Frontend Structure

```
client/src/
├── auth/              # Auth page wrapper
├── identification/    # Login/Signup forms
├── landing/           # Landing page
├── marketplace/       # Core feature
│   ├── pages/         # MarketFeed, AddMachine, MachineDetails, PurchasedMachines
│   ├── components/    # MachineCard, Terminal, Layout
│   └── style/         # Component styles
├── dashboard/         # User dashboard
└── global/            # API clients (authApi, orderApi)
```

**API clients** (src/global/api.ts):
- `authApi` → http://localhost:3001 (withCredentials: true)
- `orderApi` → http://localhost:3002 (withCredentials: true)

**Routing:**
- `/` → Landing
- `/login`, `/signup` → Auth
- `/marketplace` → Machine feed
- `/machine/:id` → Machine details (with purchase button)
- `/add-machine` → Register machine
- `/my-machines` → Owned machines
- `/purchased-machines` → Rented machines

**Key API Endpoints:**
- POST `/api/v1/order/buy` - Purchase a machine (sends to Kafka)

## Important Patterns

### Adding API Endpoints

Backend (internal/order/order.go):
```go
func (m *Module) SetupOrderRoutes() {
    g := m.E.Group("/api/v1/order")
    g.Use(middleware.EchoAuth)  // Auth middleware for all routes
    g.GET("/new-endpoint", m.NewHandler)
}
```

Frontend (src/):
```typescript
const response = await orderApi.get('/api/v1/order/new-endpoint');
```

### Database Changes

1. Update `internal/DB/schema.sql` (add table/column)
2. Update `internal/DB/query.sql` (add query methods)
3. Run `sqlc generate` to regenerate code
4. Use generated methods in handlers: `m.DB.QueryMethod(ctx, params)`

### Error Handling

Backend returns HTTP status codes with JSON/string messages. Frontend uses try-catch with console.error.

## Known Issues

1. **Process service**: Only prints to stdout, no actual Kafka consumer implementation
2. **Duplicate logic**: `internal/machine/trade.go` has unused functions that overlap with order service
3. **Missing Kafka consumer**: No service consuming purchase events from Kafka yet
4. **No deal expiration handling**: Purchase events include expiration but no cleanup/notification system

## Code Organization

**pkg/** (public, reusable):
- `auth/`: JWT generation, password hashing (bcrypt)
- `database/`: PostgreSQL connection, schema setup, seed data
- `kafka/`: Kafka producer for purchase events

**internal/** (private, service-specific):
- `user/`: Auth endpoints (signup, login, logout, get user)
- `order/`: Machine marketplace endpoints including purchase (BuyMachine → Kafka)
- `machine/`: Machine creation, SSH validation
- `middleware/`: JWT validation, CORS
- `data/`: Data transfer objects
- `DB/`: sqlc-generated database code

**cmd/**:
- `root.go`: Command-line parser for service type
- `auth.go`, `order.go`, `process.go`: Service initialization commands
