# Expense Tracker API

Personal Expense Tracker API built with Go and Beego. This repository currently contains the Part 1 backend foundation only.

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

CSV is the default required storage for this assignment. The CSV user and expense files are configured in `conf/app.conf` and will be implemented in later parts.
