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

3. **Start via Docker**
   ```
   docker-compose up -d --build
   ```

This will start PostgreSQL + application with migrations.
Server will be available at: `http://localhost:8080`

# Data Models (DTOs)

## User

```json
{
  "id": "string",
  "name": "string",
  "role": "employee | admin",
  "email": "string (unique)",
  "department": "string (optional)",
  "jobTitle": "string (optional)",
  "telegram": "string",
  "createdAt": "ISO Date string",
  "updatedAt": "ISO Date string"
}
```

## Training Request (Request)

```json
{
  "id": "string",
  "user": {
    "id": "string",
    "name": "string",
    "jobTitle": "string",
    "telegram": "string", 
  },
  "topic": "string",
  "description": "string",
  "status": "pending | approved | rejected",
  "createdAt": "ISO Date string",
  "updatedAt": "ISO Date string"
}
```

## Mentor

```json
{
  "id": "string",
  "name": "string",
  "jobTitle": "string",
  "experience": "string",
  "workload": "number (0-5)",
  "email": "string",
  "telegram": "string"
}
```

## Learning Plan Item (Plan)

```json
{
  "id": "string",
  "text": "string",
  "completed": "boolean"
}
```

## Learning Process (Learning)

```json
{
  "id": "string",
  "request": {
    "id": "string",
    "topic": "string",
    "description": "string"
  },
  "user": {
    "id": "string",
    "name": "string"
  },
  "mentor": {
    "id": "string",
    "name": "string",
    "telegram": "string",
    "jobTitle": "string",
    "experience": "string"
  },
  "status": "active | completed",
  "startDate": "ISO Date string",
  "endDate": "ISO Date string (optional)",
  "plan": "LearningPlanItem[]",
  "feedback": {
    "rating": "number",
    "comment": "string"
  },
  "notes": "string (optional)"
}
```

# API Endpoints

## /health

| Path | Method | Description          | Access | Body | Response                              | AuthRequired |
|------|--------|----------------------|--------|------|---------------------------------------|--------------|
|      | GET    | Check backend health | All    |      | "service": string<br>"status": string | -            |

## /metrics

| Path | Method | Description                | Access | Body | Response           | AuthRequired |
|------|--------|----------------------------|--------|------|--------------------|--------------|
|      | GET    | Get metrics for Prometheus | All    |      | Raw text with data | -            |

## /auth

| Path      | Method | Description                | Access | Body                                                                                                                         | Response (JSON)                 | AuthRequired |
|-----------|--------|----------------------------|--------|------------------------------------------------------------------------------------------------------------------------------|---------------------------------|--------------|
| /register | POST   | Register new user          | All    | "name": string<br>"email": string<br>"password": string<br>"department": string<br>"jobTitile": string<br>"telegram": string | User                            | -            |
| /login    | POST   | Login                      | All    | "email": string<br>"password": string                                                                                        | "token": string<br>"user": User | -            |
| /me       | GET    | Get current user's info    | All    |                                                                                                                              | User                            | +            |
| /me       | PUT    | Change current user's info | All    | "name": string<br>"email": string<br>"password": string<br>"department": string<br>"jobTitile": string<br>"telegram": string | User                            | +            |

## /users

| Path | Method | Description                  | Access | Body                                                                                                                         | Response (JSON)   | AuthRequired |
|------|--------|------------------------------|--------|------------------------------------------------------------------------------------------------------------------------------|-------------------|--------------|
| /    | GET    | Get all users                | Admin  |                                                                                                                              | "users": User\[\] | +            |
| /:id | GET    | Get user info by its `id`    | Admin  |                                                                                                                              | User              | +            |
| /:id | PUT    | Change user info by its `id` | Admin  | "name": string<br>"email": string<br>"password": string<br>"department": string<br>"jobTitile": string<br>"telegram": string | User              | +            |


## /requests

| Path | Method | Description                 | Access                                   | Body                                     | Response (JSON)         | AuthRequired |
|------|--------|-----------------------------|------------------------------------------|------------------------------------------|-------------------------|--------------|
| /    | GET    | Get all requests            | Admin                                    |                                          | "requests": Request\[\] | +            |
| /    | POST   | Create new request          | All                                      | "topic": string<br>"description": string | Request                 | +            |
| /my  | GET    | Get current user's requests | All                                      |                                          | "requests": Request\[\] | +            |
| /:id | GET    | Get request by id           | All (if id in `/my`) \| Admin otherwise  |                                          | Request                 | +            |
| /:id | PUT    | Change request by id        | All (if id in `/my`) \| Admin  otherwise | "topic": string<br>"description": string | Request                 | +            |


## /mentors

| Path | Method | Description              | Access                                           | Body                                                                                                                                   | Response (JSON)       | AuthRequired |
|------|--------|--------------------------|--------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|-----------------------|--------------|
| /    | GET    | Get all mentors          | Admin                                            |                                                                                                                                        | "mentors": Mentor\[\] | +            |
| /    | POST   | Create new mentor        | Admin                                            | "name": string<br>"jobTitle": string<br>"experience": string<br>"workload": 0 <= integer <= 5<br>"email": string<br>"telegram": string | Mentor                | +            |
| /:id | GET    | Get mentor info by id    | All (if id in `/requests/my`) \| Admin otherwise |                                                                                                                                        | Mentor                | +            |
| /:id | PUT    | Change mentor info by id | Admin                                            | "name": string<br>"jobTitle": string<br>"experience": string<br>"workload": 0 <= integer <= 5<br>"email": string<br>"telegram": string | Mentor                | +            |

## /learnings

| Path          | Method | Description                 | Access                                  | Body                                                                                                                                                                                           | Response (JSON)           | AuthRequired |
|---------------|--------|-----------------------------|-----------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------|--------------|
| /             | GET    | Get all learnings           | Admin                                   |                                                                                                                                                                                                | "learnings": Learning\[\] | +            |
| /             | POST   | Create new learning         | All                                     | "topic": string<br>"description": string                                                                                                                                                       | Learning                  | +            |
| /my           | GET    | Get user learnings          | All                                     |                                                                                                                                                                                                | "learnings": Learning\[\] | +            |
| /:id          | GET    | Get learning by id          | All (if id in `/my`) \| Admin otherwise |                                                                                                                                                                                                | Learning                  | +            |
| /:id          | PUT    | Change learning info by id  | Admin                                   | "topic": string<br>"description": string<br>"status": active \| completed<br>"plan": Plan[]<br>"feedback": {<br>  "rating": 1 <= integer <= 5 <br>  "comment": string<br>},<br>"notes": string | Learning                  | +            |
| /:id/plan     | PUT    | Change learning plan by id  | All (if id in /my) \| Admin otherwise   | "plan": Plan[]                                                                                                                                                                                 | Learning                  | +            |
| /:id/notes    | PUT    | Change learning notes by id | All (if id in /my) \| Admin otherwise   | "notes": string                                                                                                                                                                                | Learning                  | +            |
| /:id/complete | POST   | Complete learning by id     | All (if id in /my) \| Admin otherwise   | "rating": 1 <= integer <= 5<br>"comment": string                                                                                                                                               | Learning                  | +            |



## Configuration

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

## Monitoring (Roadmap)

**MVP (current version):**
- [x\] Structured logging (log/slog)
- [x\] Health check endpoint `/health`
- [x\] Middleware for logging latency
- [x\] Prometheus metrics
- [x\] Grafana dashboards
