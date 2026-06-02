<div align="center">

# Expense Tracker API

Personal Expense Tracker backend API built with Go, Beego v2, and CSV storage.

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Beego](https://img.shields.io/badge/Beego-v2-2D3748?style=for-the-badge)
![Storage](https://img.shields.io/badge/Storage-CSV-0F766E?style=for-the-badge)
![Coverage](https://img.shields.io/badge/Coverage-92.1%25-16A34A?style=for-the-badge)

</div>

---

## Overview

This API provides the required backend for a Personal Expense Tracker assignment. It supports user registration, login, CSV-backed expense CRUD, pagination, filtering, sorting, and spending summaries. CSV remains the default local storage mode, with an optional Postgres driver wired for production-style deployments.

> The backend is API-only and keeps configuration in `conf/app.conf`.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Installation](#installation)
- [Running the API](#running-the-api)
- [Response Format](#response-format)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
- [Request Examples](#request-examples)
- [Postman Testing](#postman-testing)
- [Swagger / API Documentation](#swagger--api-documentation)
- [Docker Usage](#docker-usage)
- [Deployment Guide](#deployment-guide)
- [Repository Integration](#repository-integration)
- [Storage Modes](#storage-modes)
- [Storage Driver Switching](#storage-driver-switching)
- [Postgres Production Preparation](#postgres-production-preparation)
- [Postgres User Repository](#postgres-user-repository)
- [Postgres Expense Repository](#postgres-expense-repository)
- [Engineering Decision: Configurable Storage Driver](#engineering-decision-configurable-storage-driver)
- [CSV Storage](#csv-storage)
- [Production Checklist](#production-checklist)
- [Trainer / Local Review Checklist](#trainer--local-review-checklist)
- [Deployment Platform Notes](#deployment-platform-notes)
- [Frontend Connection](#frontend-connection)
- [Testing](#testing)
- [Test Coverage](#test-coverage)
- [Notes](#notes)

## Features

| Area | Supported |
| --- | --- |
| Health check | Yes |
| User registration | Yes |
| User login | Yes |
| Expense CRUD | Yes |
| Pagination | Yes |
| Category filtering | Yes |
| Date range filtering | Yes |
| Sorting | Yes |
| Spending summary | Yes |
| CSV storage | Yes |
| Storage driver switching | Yes |
| Postgres configuration | Implemented |
| Postgres connection foundation | Implemented |
| Postgres user repository | Implemented |
| Postgres expense repository | Implemented |
| Repository integration | Controllers use repository factory |
| Swagger API documentation | Yes |
| Docker backend packaging | Yes |

## Tech Stack

| Technology | Purpose |
| --- | --- |
| Go 1.22+ | Backend language |
| Beego v2 | Web framework and routing |
| CSV | Required assignment storage |
| Postgres | Optional production storage driver |
| Docker | Backend containerization |
| Go testing package | Unit and integration-style tests |

## Project Structure

```text
backend/
|-- .dockerignore
|-- .env.example
|-- Dockerfile
|-- conf/
|   |-- app.conf
|   `-- app.prod.example.conf
|-- controllers/
|-- data/
|   `-- .gitkeep
|-- docs/
|   |-- docs.go
|   `-- swagger.json
|-- models/
|-- repositories/
|   |-- csv/
|   `-- postgres/
|-- routers/
|-- storage/
|-- utils/
|-- validators/
|-- main.go
|-- docker-compose.yml
|-- docker-compose.postgres.yml
|-- go.mod
|-- go.sum
`-- README.md
```

## Configuration

Configuration is stored in `conf/app.conf`.

```ini
appname = expense-tracker-api
httpport = 8080
runmode = dev
copyrequestbody = true
enable_swagger = true

storage_driver = csv
csv_user_file = data/users.csv
csv_expense_file = data/expenses.csv

postgres_dsn =
postgres_auto_migrate = false
```

> CSV is the default storage for this assignment. The application creates required CSV files automatically when model functions need them.

`conf/app.prod.example.conf` shows a safe Postgres production example with placeholder credentials only.

Deployment environment variables can override the storage settings in `conf/app.conf`. Environment variable loading from a local `.env` file is not required for local assignment runs.

## Installation

Install or tidy dependencies:

```bash
go mod tidy
```

## Running the API

Start the Beego application:

```bash
bee run
```

The API runs on the port configured in `conf/app.conf`.

Default local base URL:

```text
http://localhost:8080
```

## Response Format

All endpoints use a consistent JSON response shape.

### Success

```json
{
  "success": true,
  "message": "..."
}
```

### Success with Data

```json
{
  "success": true,
  "message": "...",
  "data": {}
}
```

### Error

```json
{
  "success": false,
  "message": "..."
}
```

## Authentication

Registration and login do not require headers.

Expense endpoints require the authenticated user ID in the request header:

```http
X-User-ID: 1
```

Unauthorized response:

```json
{
  "success": false,
  "message": "Unauthorized"
}
```

## API Endpoints

| Method | Endpoint | Description | Auth |
| --- | --- | --- | --- |
| GET | `/api/v1/health` | Check server status | No |
| POST | `/api/v1/auth/register` | Register a user | No |
| POST | `/api/v1/auth/login` | Login a user | No |
| POST | `/api/v1/expenses` | Create an expense | `X-User-ID` |
| GET | `/api/v1/expenses` | List expenses | `X-User-ID` |
| GET | `/api/v1/expenses/summary` | Generate spending summary | `X-User-ID` |
| GET | `/api/v1/expenses/:id` | Get one expense | `X-User-ID` |
| PUT | `/api/v1/expenses/:id` | Update an expense | `X-User-ID` |
| DELETE | `/api/v1/expenses/:id` | Delete an expense | `X-User-ID` |

### List Expense Query Parameters

| Parameter | Format | Description |
| --- | --- | --- |
| `page` | Positive integer | Page number. Default: `1` |
| `limit` | Positive integer | Items per page. Default: `10` |
| `category` | Allowed category | Filter by expense category |
| `date_from` | `YYYY-MM-DD` | Include expenses on or after this date |
| `date_to` | `YYYY-MM-DD` | Include expenses on or before this date |
| `sort_by` | `amount` or `expense_date` | Sort field |
| `sort_order` | `asc` or `desc` | Sort direction. Default: `desc` |

## Request Examples

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:

```json
{
  "success": true,
  "message": "Server is running"
}
```

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```

### Create Expense

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}'
```

### List Expenses

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?page=1&limit=10" \
  -H "X-User-ID: 1"
```

### Filter by Category

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?category=Food" \
  -H "X-User-ID: 1"
```

### Filter by Date Range

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"
```

### Sort by Amount

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?sort_by=amount&sort_order=desc" \
  -H "X-User-ID: 1"
```

### Get One Expense

```bash
curl -X GET http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

### Update Expense

```bash
curl -X PUT http://localhost:8080/api/v1/expenses/1 \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"title":"Dinner","amount":500.00,"category":"Food","note":"Family dinner","expense_date":"2025-06-11"}'
```

### Delete Expense

```bash
curl -X DELETE http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

### Spending Summary

```bash
curl -X GET "http://localhost:8080/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"
```

## Postman Testing

Use this section to test the API manually in Postman.

### Environment Variables

Create a Postman environment with these variables:

| Variable | Value |
| --- | --- |
| `base_url` | `http://localhost:8080` |
| `user_id` | `1` |

### Suggested Test Order

| Step | Method | URL |
| --- | --- | --- |
| 1 | GET | `{{base_url}}/api/v1/health` |
| 2 | POST | `{{base_url}}/api/v1/auth/register` |
| 3 | POST | `{{base_url}}/api/v1/auth/login` |
| 4 | POST | `{{base_url}}/api/v1/expenses` |
| 5 | GET | `{{base_url}}/api/v1/expenses?page=1&limit=10` |
| 6 | GET | `{{base_url}}/api/v1/expenses/1` |
| 7 | PUT | `{{base_url}}/api/v1/expenses/1` |
| 8 | GET | `{{base_url}}/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30` |
| 9 | DELETE | `{{base_url}}/api/v1/expenses/1` |

### Auth Requests

Register request body:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}
```

Login request body:

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

### Expense Request Headers

For all expense endpoints, add this header in Postman:

| Key | Value |
| --- | --- |
| `X-User-ID` | `{{user_id}}` |

For requests with a JSON body, also add:

| Key | Value |
| --- | --- |
| `Content-Type` | `application/json` |

### Expense Request Bodies

Create expense body:

```json
{
  "title": "Lunch",
  "amount": 350.50,
  "category": "Food",
  "note": "Team lunch",
  "expense_date": "2025-06-10"
}
```

Update expense body:

```json
{
  "title": "Dinner",
  "amount": 500.00,
  "category": "Food",
  "note": "Family dinner",
  "expense_date": "2025-06-11"
}
```

Expected successful health response:

```json
{
  "success": true,
  "message": "Server is running"
}
```

Expected unauthorized response when `X-User-ID` is missing or invalid:

```json
{
  "success": false,
  "message": "Unauthorized"
}
```

## Swagger / API Documentation

Swagger is enabled for local API exploration when this config value is true:

```ini
enable_swagger = true
```

The Swagger docs reflect the current REST API under `/api/v1`. Protected expense and summary routes use the `X-User-ID` header. After login, copy the returned `user_id` value into `X-User-ID` before calling expense endpoints.

Passwords are hashed internally and never returned in API responses.

Generate docs:

```bash
bee generate docs
```

Run the API:

```bash
bee run
```

Open Swagger UI:

```text
http://localhost:8080/swagger/
```

The served JSON spec is available at:

```text
http://localhost:8080/swagger/doc.json
```

## Docker Usage

Docker is configured for the backend only. Run these commands from the `backend/` folder after opening Docker Desktop.

### Build Image Manually

```bash
cd backend
docker build -t expense-tracker-api .
```

### Run With CSV Storage

```bash
cd backend
docker compose up --build
```

Test the running container:

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:

```json
{
  "success": true,
  "message": "Server is running"
}
```

CSV mode is the default Docker mode. It stores generated CSV files in a named Docker volume and does not require Postgres.

Stop CSV mode:

```bash
docker compose down
```

Remove the CSV Docker volume:

```bash
docker compose down -v
```

### Run With Local Postgres

```bash
cd backend
docker compose -f docker-compose.postgres.yml up --build
```

This starts a local Postgres container, sets `STORAGE_DRIVER=postgres`, and enables `POSTGRES_AUTO_MIGRATE=true` so the backend can create required tables.

Stop local Postgres mode:

```bash
docker compose -f docker-compose.postgres.yml down
```

Remove the local Postgres volume:

```bash
docker compose -f docker-compose.postgres.yml down -v
```

### Run With Neon Postgres

Bash:

```bash
docker run --name expense-tracker-api \
  -p 8080:8080 \
  -e PORT=8080 \
  -e STORAGE_DRIVER=postgres \
  -e POSTGRES_DSN="YOUR_NEON_CONNECTION_STRING" \
  -e POSTGRES_AUTO_MIGRATE=true \
  expense-tracker-api
```

PowerShell:

```powershell
docker run --name expense-tracker-api `
  -p 8080:8080 `
  -e PORT=8080 `
  -e STORAGE_DRIVER=postgres `
  -e POSTGRES_DSN="YOUR_NEON_CONNECTION_STRING" `
  -e POSTGRES_AUTO_MIGRATE=true `
  expense-tracker-api
```

Replace `YOUR_NEON_CONNECTION_STRING` with the real Neon DSN and keep it out of Git. Neon usually requires `sslmode=require`.

### Docker Desktop Notes

- Open Docker Desktop before running Docker commands.
- After `docker compose up`, check Docker Desktop -> Containers for logs.
- API base URL: `http://localhost:8080`.
- Swagger UI, when enabled: `http://localhost:8080/swagger/`.
- The frontend remains deployed separately on Vercel and is not Dockerized here.

### Render And Vercel Alignment

Render backend environment variables:

```bash
STORAGE_DRIVER=postgres
POSTGRES_DSN=<Neon connection string>
POSTGRES_AUTO_MIGRATE=true
```

Vercel frontend environment variable:

```bash
NEXT_PUBLIC_API_BASE_URL=https://YOUR_RENDER_BACKEND_URL/api/v1
```

## Deployment Guide

The backend has two storage modes with the same API contract: CSV for local assignment review and Postgres for production deployment.

### Local Assignment Mode: CSV

CSV is the default mode. It satisfies the assignment requirement, does not require a database, and automatically creates runtime CSV files when the app needs them. This is the mode a trainer can use directly after cloning the project.

Run locally:

```bash
go mod tidy
bee run
```

Local config:

```ini
storage_driver = csv
csv_user_file = data/users.csv
csv_expense_file = data/expenses.csv
```

### Production Mode: Postgres

Use Postgres for live deployment. This avoids depending on runtime CSV files in hosted environments and keeps API behavior the same as CSV mode.

Production config:

```ini
storage_driver = postgres
postgres_dsn = postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=require
postgres_auto_migrate = true
```

Equivalent Render environment variables:

```bash
STORAGE_DRIVER=postgres
POSTGRES_DSN=postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=require
POSTGRES_AUTO_MIGRATE=true
```

When `postgres_auto_migrate = true`, the app runs `repositories/postgres/schema.sql` at startup. When `postgres_auto_migrate = false`, the database schema must already exist.

## Repository Integration

Controllers now depend on repository interfaces through the repository factory while keeping the current API behavior unchanged.

- `repositories.UserRepository` defines user persistence behavior.
- `repositories.ExpenseRepository` defines expense persistence behavior.
- `repositories.NewUserRepository()` provides the active user repository.
- `repositories.NewExpenseRepository()` provides the active expense repository.
- CSV repository adapters wrap the existing model functions.
- CSV remains the active implementation by default.
- CSV behavior remains unchanged.
- Postgres repositories are used when `storage_driver = postgres`.
- CSV remains the default for local assignment mode.

Current storage architecture status:

| Component | Status |
| --- | --- |
| CSV storage | Implemented and default |
| Storage config helper | Implemented |
| Repository interfaces | Implemented |
| CSV repository adapters | Implemented |
| Controllers use repository factory | Implemented |
| Postgres connection helper | Implemented |
| Postgres schema and migration helper | Implemented |
| Postgres user repository | Implemented |
| Postgres expense repository | Implemented |
| Production driver switch | Implemented |

## Storage Modes

### CSV Storage - Default Assignment Mode

CSV is the default mode because the assignment requires CSV file storage.

- It creates `data/users.csv` and `data/expenses.csv` automatically.
- It needs no database setup.
- It is the recommended mode for local review.

### Postgres Storage - Production Mode

Postgres is the optional production mode.

- It is enabled with `storage_driver = postgres`.
- It requires a valid `postgres_dsn`.
- It keeps the same API behavior as CSV mode.
- It is not required for local assignment testing.

## Storage Driver Switching

The backend supports runtime storage selection through `conf/app.conf`.

### Local / Assignment Mode

```ini
storage_driver = csv
csv_user_file = data/users.csv
csv_expense_file = data/expenses.csv
```

CSV remains the default mode and requires no database. This is the recommended mode for local assignment review with `bee run`.

### Postgres Production Mode

```ini
storage_driver = postgres
postgres_dsn = postgres://USER:PASSWORD@HOST:PORT/DB?sslmode=require
postgres_auto_migrate = true
```

When `storage_driver = postgres`, application startup opens the shared Postgres connection before the Beego server starts. If `postgres_auto_migrate = true`, embedded schema migrations run during storage initialization.

If Postgres mode is selected and the DSN is missing or invalid, startup fails fast with a clear error. CSV mode does not open Postgres and does not require Postgres to be installed or running.

## Postgres Production Preparation

Postgres connection and schema foundations are available for production deployment. CSV remains the default local storage mode, and Postgres is not required to run or test the assignment locally.

Prepared Postgres pieces:

- `repositories/postgres.OpenDB`
- `repositories/postgres.ConfigurePool`
- `repositories/postgres.RunMigrations`
- `repositories/postgres/schema.sql`
- pgx stdlib driver registration for `database/sql`

Auto-migration support is implemented through `schema.sql` and `RunMigrations`. Migrations run at startup only when `storage_driver = postgres` and `postgres_auto_migrate = true`.

### Schema Overview

`users` table:

| Column |
| --- |
| `id` |
| `name` |
| `email` |
| `password` |
| `created_at` |

`expenses` table:

| Column |
| --- |
| `id` |
| `user_id` |
| `title` |
| `amount` |
| `category` |
| `note` |
| `expense_date` |
| `created_at` |

### Production Config

Production mode uses:

```ini
storage_driver = postgres
postgres_dsn = postgres://USER:PASSWORD@HOST:PORT/DB?sslmode=require
postgres_auto_migrate = true
```

Default local assignment mode remains:

```ini
storage_driver = csv
```

Postgres repositories are selected by the repository factory when `storage_driver = postgres`.

## Postgres User Repository

The Postgres user repository has been implemented for production storage support.

Supported user operations:

- Create user
- Get all users
- Get user by email
- Get user by ID
- Get next ID for interface compatibility

The repository is tested with SQL mocks, so local tests do not require a live Postgres database. CSV remains the default storage driver, and the repository factory switches to Postgres only when `storage_driver = postgres`.

## Postgres Expense Repository

The Postgres expense repository has been implemented for production storage support.

Supported expense operations:

- Create expense
- Get expenses by user ID
- Get expense by ID with ownership check
- Update expense with ownership check
- Delete expense with ownership check
- Get next expense ID for interface compatibility

The repository is tested with SQL mocks, so local tests do not require a live Postgres database. CSV remains the default storage driver, and the repository factory switches to Postgres only when `storage_driver = postgres`.

## Configurable Storage Driver

This project intentionally uses a configurable storage driver instead of hardcoding one storage system.

The assignment requires CSV file storage, so CSV remains the default local driver. This keeps the project fully aligned with the assignment and allows a trainer to clone the repository and run the backend locally without installing a database.

For production deployment, the same API can switch to Postgres by changing configuration:

```ini
storage_driver = postgres
```

This design keeps the local assignment experience simple while making the deployed version more reliable.

Architecture:

```text
HTTP Request
   |
Controller
   |
Repository Interface
   |
Storage Driver Factory
   |
CSV Repository        Postgres Repository
(local assignment)    (production deployment)
```

Controllers do not know whether data comes from CSV or Postgres. The API response format stays the same, the same endpoints work in both storage modes, and the repository factory selects the correct implementation based on `storage_driver`.

`storage.Init()` opens a Postgres connection only when `storage_driver = postgres`. CSV mode never opens or requires Postgres. Postgres migrations can run automatically when `postgres_auto_migrate = true`.

## Reasoning Behind the Storage Architecture:

- It follows separation of concerns.
- It keeps business and API logic independent from storage details.
- It protects the assignment requirement by keeping CSV as default.
- It makes production deployment practical without rewriting controllers.
- It allows future storage implementations to be added with minimal API changes.
- It makes the codebase easier to test because repositories can be tested independently.
- It avoids mixing CSV and Postgres logic inside controllers.
- It keeps local development lightweight and production deployment reliable.

This means the project has two modes with the same API contract:

Local review:

```ini
storage_driver = csv
```

Uses `data/users.csv` and `data/expenses.csv`.

Production:

```ini
storage_driver = postgres
```

Uses Postgres tables with the same API behavior.

This approach demonstrates that the project satisfies the assignment requirement while also considering how the same backend would run in a real hosted environment.

## CSV Storage

### User CSV Format

```text
id,name,email,password,created_at
```

### Expense CSV Format

```text
id,user_id,title,amount,category,note,expense_date,created_at
```

### Allowed Expense Categories

| Category |
| --- |
| Food |
| Transport |
| Housing |
| Entertainment |
| Shopping |
| Healthcare |
| Education |
| Utilities |
| Other |

Generated CSV files are ignored by Git so local test and runtime data are not committed.

## Production Checklist

- Set `storage_driver = postgres`.
- Set `postgres_dsn` with the production database URL.
- Set `postgres_auto_migrate = true` for first deployment, or run the schema manually before startup.
- Never commit real database credentials.
- Confirm the health endpoint works.
- Test register and login.
- Test expense create, list, update, and delete.
- Test filtering and sorting.
- Test the summary endpoint.
- Set frontend `NEXT_PUBLIC_API_BASE_URL` to the deployed backend URL.

## Trainer / Local Review Checklist

- Clone the repository.
- Keep `storage_driver = csv`.
- Run `go mod tidy`.
- Run `go test ./...`.
- Run `bee run`.
- Test the health endpoint.
- Register and login.
- Create, list, update, and delete expenses.
- Test filtering, sorting, and summary.

No Postgres setup is needed for local review.

## Frontend Connection

The Next.js frontend should point to the deployed backend API with:

```ini
NEXT_PUBLIC_API_BASE_URL=https://YOUR_BACKEND_URL/api/v1
```

For local development:

```ini
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

## Testing

| Command | Purpose |
| --- | --- |
| `go test ./...` | Run all tests |
| `go test ./... -cover` | Run tests with package coverage |
| `go test ./... -coverprofile=coverage.out` | Generate coverage profile |
| `go tool cover -func coverage.out` | Show function-level coverage |
| `go vet ./...` | Run static checks |
| `gofmt -w .` | Format Go files |

## Test Coverage

Total statement coverage: **91.3%**

The project includes unit and integration-style tests for controllers, models, validators, CSV utilities, and route registration.

### Coverage Screenshots

![Test coverage result](../docs/images/test-coverage1.png)
![Test coverage result](../docs/images/test-coverage2.png)
![Test coverage result](../docs/images/test-coverage3.png)

## Password Security

- Passwords are hashed with bcrypt before storage.
- CSV mode stores bcrypt hashes in `data/users.csv`.
- Postgres mode stores bcrypt hashes in `users.password`.
- Plain-text passwords are never returned in API responses.
- Local users created before bcrypt was added may need to be recreated. Delete `data/users.csv` and register again.

## Notes

> This backend intentionally stays assignment-focused and does not include bonus features.

- Passwords are hashed with bcrypt before being stored in CSV or Postgres.
- Postgres storage is wired for production mode, while CSV remains the default for local assignment runs.
- Frontend containerization, voice input, budget features, and export features are not included.
