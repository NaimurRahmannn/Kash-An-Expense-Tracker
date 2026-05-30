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
- [Repository Integration](#repository-integration)
- [Storage Modes](#storage-modes)
- [Storage Driver Switching](#storage-driver-switching)
- [Postgres Production Preparation](#postgres-production-preparation)
- [Postgres User Repository](#postgres-user-repository)
- [Postgres Expense Repository](#postgres-expense-repository)
- [Storage Strategy Decision](#storage-strategy-decision)
- [CSV Storage](#csv-storage)
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

## Tech Stack

| Technology | Purpose |
| --- | --- |
| Go 1.22+ | Backend language |
| Beego v2 | Web framework and routing |
| CSV | Required assignment storage |
| Postgres | Optional production storage driver |
| Go testing package | Unit and integration-style tests |

## Project Structure

```text
backend/
|-- conf/
|   `-- app.conf
|-- controllers/
|-- data/
|   `-- .gitkeep
|-- models/
|-- repositories/
|   |-- csv/
|   `-- postgres/
|-- routers/
|-- utils/
|-- validators/
|-- main.go
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

storage_driver = csv
csv_user_file = data/users.csv
csv_expense_file = data/expenses.csv

postgres_dsn =
postgres_auto_migrate = false
```

> CSV is the default storage for this assignment. The application creates required CSV files automatically when model functions need them.

`.env.example` is included as production reference documentation only. Environment variable loading is not required for local assignment runs.

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

### CSV Storage - Default Local Mode

CSV is the default storage mode because the assignment requires CSV file storage.

- Trainers can run the project locally with `bee run`.
- No database setup is required.
- Runtime files are generated automatically:
  - `data/users.csv`
  - `data/expenses.csv`
- Generated CSV files are ignored by Git.

### Postgres Storage - Production Mode

Postgres support is available as an optional production storage mode.

- It is enabled with `storage_driver = postgres`.
- It uses `postgres_dsn` for the database connection string.
- It is useful for production because hosted environments may not persist local CSV files reliably.
- Postgres is not required for local assignment testing.

Current status:

| Storage Item | Status |
| --- | --- |
| CSV storage | Fully implemented |
| Postgres configuration | Implemented |
| Postgres connection helper | Implemented |
| Postgres schema SQL | Implemented |
| Optional auto-migration helper | Implemented |
| Postgres user repository | Implemented |
| Postgres expense repository | Implemented |
| Production driver switch | Implemented |

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

## Storage Strategy Decision

The assignment requires CSV storage, so CSV remains the source of truth for local development and trainer testing.

Production deployment benefits from Postgres because it provides durable persistence, safer concurrent writes, and better scalability than local CSV files. The project is being prepared so API behavior can remain the same regardless of the selected storage driver.

This is a deliberate tradeoff: a little more architecture complexity later, but better deployment reliability. CSV remains the default to keep local setup simple.

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

Total statement coverage: **92.1%**

The project includes unit and integration-style tests for controllers, models, validators, CSV utilities, and route registration.

### Coverage Screenshots

![Test coverage result](../docs/images/test-coverage1.png)
![Test coverage result](../docs/images/test-coverage2.png)

## Notes

> This backend intentionally stays assignment-focused and does not include bonus features.

- Passwords are stored as plain text for assignment compatibility at this stage.
- Postgres storage is wired for production mode, while CSV remains the default for local assignment runs.
- Docker, Swagger, frontend, voice input, budget features, and export features are not included.
