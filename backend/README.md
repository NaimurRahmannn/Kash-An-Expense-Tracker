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

This API provides the required backend for a Personal Expense Tracker assignment. It supports user registration, login, CSV-backed expense CRUD, pagination, filtering, sorting, and spending summaries.

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

## Tech Stack

| Technology | Purpose |
| --- | --- |
| Go 1.22+ | Backend language |
| Beego v2 | Web framework and routing |
| CSV | Required assignment storage |
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
```

> CSV is the default storage for this assignment. The application creates required CSV files automatically when model functions need them.

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

![Test coverage result](docs/images/test-coverage1.png)
![Test coverage result](docs/images/test-coverage2.png)

## Notes

> This backend intentionally stays assignment-focused and does not include bonus features.

- Passwords are stored as plain text for assignment compatibility at this stage.
- Postgres, Docker, Swagger, frontend, voice input, budget features, and export features are not included.
