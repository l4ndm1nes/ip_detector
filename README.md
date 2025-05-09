# IP Detector API

IP Detector API is a RESTful service for user registration, authentication, and country detection based on IP address. Built with Go and secured with JWT.

---

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/l4ndm1nes/ip_detector.git
cd ip_detector
```

### 2. Set Up Environment Variables
Create a .env file in the project root:
```bash
DB_NAME=ip_detector
DB_USER=ipd
DB_PASSWORD=ipdpass
POSTGRES_DSN=postgres://ipd:ipdpass@db:5432/ip_detector?sslmode=disable

JWT_SECRET=supersecretkey
JWT_EXPIRATION=24h
```

### 3. Run with Docker Compose
```bash
make run
```

API available at: http://localhost:8080

Swagger UI: http://localhost:8080/swagger/index.html

## Usage
### Registration
POST /register

Register a new user:
```bash
{
  "name": "John Doe",
  "email": "john@example.com",
  "ip": "8.8.8.8",
  "password": "secret123"
}
```
### Login
POST /login

Get a JWT token:
```bash
{
  "email": "john@example.com",
  "password": "secret123"
}
```

### Protected Endpoints (JWT Required)
GET /users - List all users

GET /users/{id} - Get user by ID

#### Example:

Use the JWT token in Authorization header:
```bash
Bearer your_jwt_token
```

## Commands
Run the service:
```bash
make run
```
Run tests:
```bash
make test
```
Build the binary:
```bash
make build
```
Apply database migrations:
```bash
make migrate
```
Rollback database migrations:
```bash
make migrate-down
```
Clean up Docker resources:
```bash
make clean
```

## Project Structure
```bash
.
├── cmd/
│   └── httpserver/       - Main server entry point
├── internal/
│   ├── adapter/          - External integrations (DB, GeoIP)
│   ├── app/              - Core application logic (services)
│   ├── config/           - Configuration loading
│   ├── domain/           - Domain models and interfaces
│   └── logger/           - Centralized logging setup
├── migrations/           - Database migrations (SQL)
├── Dockerfile            - Docker setup
├── docker-compose.yml    - Docker Compose setup
├── Makefile              - Simplified project management
└── README.md             - Project documentation
```