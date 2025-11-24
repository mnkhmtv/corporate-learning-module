# Internal Training & Mentoring System

A corporate training system for organizing employee learning through mentorship and courses.

## Description

A simple and intuitive backend service for managing corporate training processes. Employees can create training requests, administrators assign mentors, and the system tracks progress and collects feedback.

### Key Features

- **Request Management** — employees create learning requests for specific topics
- **Mentor Assignment** — administrators match mentors considering their workload
- **Learning Process** — collaborative task planning and progress tracking
- **Feedback System** — ratings and comments after training completion
- **Personal Dashboard** — application history and current learning status

## Architecture

The project is built on **Clean Architecture** principles with clear layer separation:

```
Domain (Entities) → Service (Use Cases) → Repository (Data) ← Transport (HTTP/Gin)
```

**Benefits:**
- Business logic independent from frameworks and databases
- Easy testing through mock interfaces
- Flexibility to replace components (e.g., switching from Postgres to another DB)

## Tech Stack

| Component           | Technology              | Version |
|---------------------|-------------------------|---------|
| **Language**        | Go                      | 1.25+   |
| **Web Framework**   | Gin                     | latest  |
| **Database**        | PostgreSQL              | 15+     |
| **DB Driver**       | pgx/v5                  | latest  |
| **Configuration**   | cleanenv                | latest  |
| **Authentication**  | golang-jwt/jwt          | v5      |
| **Logging**         | log/slog (stdlib)       | -       |
| **Migrations**      | golang-migrate/migrate  | latest  |
| **Password Hash**   | bcrypt (stdlib)         | -       |


## Structure

```
backend/
├── cmd/app/main.go              # Entry point
├── config/                      # Configuration files
├── internal/
│   ├── domain/                  # Domain models (User, Mentor, TrainingRequest...)
│   ├── pkg/                     # Private inner utils
│   ├── repository/              # Database layer
│   ├── service/                 # Business logic
│   └── transport/http/          # Gin handlers + middleware
├── pkg/                         # Public reusable utilities
├── docker-compose.yml
├── Makefile
└── README.md
```

For more details, see comments in the source code.

## Quick Start

### Prerequisites

- Go 1.25
- Docker & Docker Compose
- Make (optional)

### Local Setup

1. **Clone the repository**
   ```
   git clone https://github.com/your-username/internal-training-system.git
   cd internal-training-system
   ```

2. **Configure environment variables**
   ```
   cp .env.example .env
   # Edit .env with your parameters
   ```

3. **Start PostgreSQL via Docker**
   ```
   docker-compose up -d postgres
   ```

4. **Run migrations**
   ```
   make migrate-up
   # or
   migrate -path ./internal/repository/migrations -database "postgres://user:pass@localhost:5432/training_db?sslmode=disable" up
   ```

5. **Start the application**
   ```
   make run
   # or
   go run cmd/app/main.go
   ```

Server will be available at: `http://localhost:8080`

### Docker (Full Stack)

```
docker-compose up --build
```

This will start PostgreSQL + application with migrations.

## API Endpoints

### Main endpoint groups:

**Authentication**
- `POST /api/auth/register` — user registration
- `POST /api/auth/login` — login and JWT token retrieval

**Training Requests**
- `GET /api/requests` — all requests (admin)
- `POST /api/requests` — create request
- `POST /api/requests/:id/assign` — assign mentor

**Learning Process**
- `GET /api/learnings` — my learnings
- `PUT /api/learnings/:id/plan` — update plan
- `POST /api/learnings/:id/complete` — complete with feedback

**Mentors**
- `GET /api/mentors` — list mentors
- `POST /api/mentors` — add mentor

## ⚙️ Configuration

Settings are located in `config/config.yaml` and can be overridden via `.env`:

```
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: training_db
  sslmode: disable

auth:
  jwt_secret: your-secret-key
  token_ttl: 24h
```

## 🧪 Testing

```
# Unit tests
go test ./internal/service/... -v

# Integration tests
make test-integration

# Code coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🔧 Development

### Creating a new migration

```
migrate create -ext sql -dir internal/repository/migrations -seq add_new_table
```

### Useful Makefile commands

```
make run              # Start application
make migrate-up       # Apply migrations
make migrate-down     # Rollback last migration
make test             # Run tests
make lint             # Code check (golangci-lint)
make docker-build     # Build Docker image
```

## 📊 Monitoring (Roadmap)

**MVP (current version):**
- [x\] Structured logging (log/slog)
- [x\] Health check endpoint `/health`
- [x\] Middleware for logging latency

**Planned:**
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Distributed tracing (OpenTelemetry)