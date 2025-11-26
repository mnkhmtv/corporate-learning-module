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

Server will be available at: `http://localhost:8080`

### Docker (Full Stack)

```
docker-compose up --build
```

This will start PostgreSQL + application with migrations.

## API Endpoints (may change by some time)

### Authentication

| Method | Endpoint             | Description             | Auth Required |
|--------|----------------------|-------------------------|---------------|
| POST   | `/api/auth/register` | Register a new user     | ❌             |
| POST   | `/api/auth/login`    | Login and get JWT token | ❌             |

### Training Requests

| Method | Endpoint                   | Description                                    | Auth Required | Role  |
|--------|----------------------------|------------------------------------------------|---------------|-------|
| POST   | `/api/requests`            | Create a training request                      | ✅             | All   |
| GET    | `/api/requests/my`         | Get my training requests                       | ✅             | All   |
| GET    | `/api/requests`            | Get all requests (with optional status filter) | ✅             | Admin |
| POST   | `/api/requests/:id/assign` | Assign mentor to request                       | ✅             | Admin |

### Mentors

| Method | Endpoint           | Description         | Auth Required | Role  |
|--------|--------------------|---------------------|---------------|-------|
| GET    | `/api/mentors`     | Get all mentors     | ✅             | All   |
| GET    | `/api/mentors/:id` | Get mentor by ID    | ✅             | All   |
| POST   | `/api/mentors`     | Create a new mentor | ✅             | Admin |

### Learning Processes

| Method | Endpoint                      | Description                     | Auth Required | Role |
|--------|-------------------------------|---------------------------------|---------------|------|
| GET    | `/api/learnings`              | Get my learning processes       | ✅             | All  |
| GET    | `/api/learnings/:id`          | Get learning process by ID      | ✅             | All  |
| GET    | `/api/learnings/:id/progress` | Get completion progress (%)     | ✅             | All  |
| PUT    | `/api/learnings/:id/plan`     | Update entire learning plan     | ✅             | All  |
| PATCH  | `/api/learnings/:id/notes`    | Update learning notes           | ✅             | All  |
| POST   | `/api/learnings/:id/complete` | Complete learning with feedback | ✅             | All  |

### Health Check

| Method | Endpoint  | Description      | Auth Required |
|--------|-----------|------------------|---------------|
| GET    | `/health` | API health check | ❌             |


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

## 📊 Monitoring (Roadmap)

**MVP (current version):**
- [x\] Structured logging (log/slog)
- [x\] Health check endpoint `/health`
- [x\] Middleware for logging latency
- [x\] Prometheus metrics
- [x\] Grafana dashboards
