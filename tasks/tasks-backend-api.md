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

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [ ] 0.0 Create feature branch
  - [ ] 0.1 Create and checkout a new branch for this feature (e.g., `git checkout -b feature/backend-api`)
- [ ] 1.0 Initialize Go Project and Core Dependencies
  - [ ] 1.1 Run `go mod init` in the root directory
  - [ ] 1.2 Install core dependencies (Gin, Viper, pq driver, golang-jwt, etc.) using `go get`
  - [ ] 1.3 Create the basic folder structure (`cmd/`, `internal/`, `pkg/`, `config/`)
- [ ] 2.0 Setup Configuration and Database Connection (PostgreSQL)
  - [ ] 2.1 Setup `viper` to read configurations from `.env` or system environment
  - [ ] 2.2 Implement a package in `pkg/database/` to connect to PostgreSQL using `database/sql`
  - [ ] 2.3 Add connection string variables in `.env` matching `docker-compose.yml`
- [ ] 3.0 Implement Database Migrations (golang-migrate)
  - [ ] 3.1 Setup `golang-migrate` CLI or library to handle schema migrations
  - [ ] 3.2 Create migration files (`.up.sql` and `.down.sql`) for `users`, `classes`, `students`, and `attendances` tables based on the PRD
  - [ ] 3.3 Add logic in `main.go` or a separate script to run migrations on startup
- [ ] 4.0 Implement Authentication System (JWT)
  - [ ] 4.1 Create JWT utility functions (generate and validate tokens) in `pkg/jwt/`
  - [ ] 4.2 Create Authentication Middleware in `internal/middleware/` to protect API routes
  - [ ] 4.3 Implement `POST /auth/login` handler
  - [ ] 4.4 Implement `POST /auth/logout` and `POST /auth/refresh-token` handlers
- [ ] 5.0 Implement Users and Classes API
  - [ ] 5.1 Implement `users` repository, service, and handlers (`GET /users/me`, `GET /users`, `POST /users`, `PUT /users/{user_id}`)
  - [ ] 5.2 Implement `classes` repository, service, and handlers (`GET /classes`, `GET /classes/{class_id}`)
- [ ] 6.0 Implement Students API and QR Code Generation
  - [ ] 6.1 Implement `students` repository, service, and handlers (CRUD for students)
  - [ ] 6.2 Implement QR Code generation data logic
- [ ] 7.0 Implement Attendance API (Scanning & Manual)
  - [ ] 7.1 Implement repository, service, and handlers for QR code scanning (validate student and mark attendance)
  - [ ] 7.2 Implement repository, service, and handlers for Manual Attendance Input by teacher
  - [ ] 7.3 Implement Audit Log insertion for any attendance changes
- [ ] 8.0 Setup Router (Gin Gonic) and Middleware
  - [ ] 8.1 Initialize Gin engine in `main.go` or `cmd/api/`
  - [ ] 8.2 Add global middleware (CORS, Logger, Recovery)
  - [ ] 8.3 Register all routes and group them under `/api/v1`
  - [ ] 8.4 Apply Authentication Middleware to protected routes
- [ ] 9.0 Ensure Compatibility with Docker Configuration
  - [ ] 9.1 Verify that the application port (`8080`) matches the Dockerfile `EXPOSE` port
  - [ ] 9.2 Verify that the database connection string matches the PostgreSQL service name in `docker-compose.yml`
  - [ ] 9.3 Ensure the built Go binary correctly references `.env` or reads from environment variables
