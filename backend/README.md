# Expense Tracker API

Personal Expense Tracker API built with Go and Beego. This repository currently contains the backend foundation for the assignment.

## Tech Stack

- Go 1.22+
- Beego v2
- CSV storage configured as the default required storage

## Install Dependencies

```bash
go mod tidy
```

## Run

```bash
bee run
```

The API starts on the port configured in `conf/app.conf`.

## Test

```bash
go test ./...
```

## Storage

CSV is the default required storage for this assignment. The CSV user and expense files are configured in `conf/app.conf`.

## Part 2 Completed

- User model added with CSV-backed functions for loading, email lookup, creation, and next ID calculation.
- Reusable CSV utilities added in `utils`.
- `users.csv` is created automatically with the required header when user model functions need it.
- Tests are placed beside the source files in `models` and `utils`.

```bash
go test ./...
```

## Part 3 Completed

- Register endpoint added at `POST /api/v1/auth/register`.
- Login endpoint added at `POST /api/v1/auth/login`.
- Request validation added for registration and login.

Register:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```

## Part 4 Completed

- Expense model added with CSV-backed storage functions.
- Allowed expense categories added.
- Expense create, read, update, delete storage functions added in `models`.
- Expense validation helpers added in `validators`.
- Expense model and validator tests added beside source files.
- Expense API endpoints are not added yet; they will be added in the next part.

## Part 5 Completed

- Expense authentication via `X-User-ID` added.
- Create expense endpoint added at `POST /api/v1/expenses`.
- Expense validation connected to the API.
- List, get, update, delete, and summary expense endpoints are not added yet.

Create expense:

```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}'
```

## Part 6 Completed

- List expenses endpoint added at `GET /api/v1/expenses`.
- Get single expense endpoint added at `GET /api/v1/expenses/:id`.
- Basic pagination added with `page` and `limit` query parameters.
- Ownership protection added so users can only retrieve their own expenses.
- Update, delete, filtering, sorting, and summary endpoints are not added yet.

List expenses:

```bash
curl -X GET "http://localhost:8080/api/v1/expenses?page=1&limit=10" \
  -H "X-User-ID: 1"
```

Get one expense:

```bash
curl -X GET http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```
