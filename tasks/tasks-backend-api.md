## Relevant Files

- `main.go` - Entry point of the application.
- `go.mod` / `go.sum` - Go dependencies.
- `config/` - Configuration management (Viper).
- `internal/` - Core application logic (handlers, services, repositories).
- `pkg/` - Reusable packages (database connection, jwt, utils).
- `migrations/` - SQL migration files for golang-migrate.

### Notes

- The Backend API will be built using Golang, Gin Gonic, and PostgreSQL as specified in the PRD.
- Ensure the application runs on the port specified in your Docker configuration.
- Use `golang-migrate` for database schema management.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [x]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Create and checkout a new branch for this feature (e.g., `git checkout -b feature/backend-api`)
- [x] 1.0 Initialize Go Project and Core Dependencies
  - [x] 1.1 Run `go mod init` in the root directory
  - [x] 1.2 Install core dependencies (Gin, Viper, pq driver, golang-jwt, etc.) using `go get`
  - [x] 1.3 Create the basic folder structure (`cmd/`, `internal/`, `pkg/`, `config/`)
- [x] 2.0 Setup Configuration and Database Connection (PostgreSQL)
  - [x] 2.1 Setup `viper` to read configurations from `.env` or system environment
  - [x] 2.2 Implement a package in `pkg/database/` to connect to PostgreSQL using `database/sql`
  - [x] 2.3 Add connection string variables in `.env` matching `docker-compose.yml`
- [x] 3.0 Implement Database Migrations (golang-migrate)
  - [x] 3.1 Setup `golang-migrate` CLI or library to handle schema migrations
  - [x] 3.2 Create migration files (`.up.sql` and `.down.sql`) for `users`, `classes`, `students`, and `attendances` tables based on the PRD
  - [x] 3.3 Add logic in `main.go` or a separate script to run migrations on startup
- [x] 4.0 Implement Authentication System (JWT)
  - [x] 4.1 Create JWT utility functions (generate and validate tokens) in `pkg/jwt/`
  - [x] 4.2 Create Authentication Middleware in `internal/middleware/` to protect API routes
  - [x] 4.3 Implement `POST /auth/login` handler
  - [x] 4.4 Implement `POST /auth/logout` and `POST /auth/refresh-token` handlers
- [x] 5.0 Implement Users and Classes API
  - [x] 5.1 Implement `users` repository, service, and handlers (`GET /users/me`, `GET /users`, `POST /users`, `PUT /users/{user_id}`)
  - [x] 5.2 Implement `classes` repository, service, and handlers (`GET /classes`, `GET /classes/{class_id}`)
- [x] 6.0 Implement Students API and QR Code Generation
  - [x] 6.1 Implement `students` repository, service, and handlers (CRUD for students)
  - [x] 6.2 Implement QR Code generation data logic
- [x] 7.0 Implement Attendance API (Scanning & Manual)
  - [x] 7.1 Implement repository, service, and handlers for QR code scanning (validate student and mark attendance)
  - [x] 7.2 Implement repository, service, and handlers for Manual Attendance Input by teacher
  - [x] 7.3 Implement Audit Log insertion for any attendance changes
- [x] 8.0 Setup Router (Gin Gonic) and Middleware
  - [x] 8.1 Initialize Gin engine in `main.go` or `cmd/api/`
  - [x] 8.2 Add global middleware (CORS, Logger, Recovery)
  - [x] 8.3 Register all routes and group them under `/api/v1`
  - [x] 8.4 Apply Authentication Middleware to protected routes
- [x] 9.0 Ensure Compatibility with Docker Configuration
  - [x] 9.1 Verify that the application port (`8080`) matches the Dockerfile `EXPOSE` port
  - [x] 9.2 Verify that the database connection string matches the PostgreSQL service name in `docker-compose.yml`
  - [x] 9.3 Ensure the built Go binary correctly references `.env` or reads from environment variables
