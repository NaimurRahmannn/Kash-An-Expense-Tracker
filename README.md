# Kash — Personal Expense Tracker

A personal expense tracker built for the internship assignment, with a Go + Beego backend, CSV-based local storage, optional Postgres production storage, and a Next.js frontend with dashboard, expense management, and voice input feature.

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Beego](https://img.shields.io/badge/Beego-v2-2D3748?style=for-the-badge)
![Next.js](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![CSV](https://img.shields.io/badge/Storage-CSV-0F766E?style=for-the-badge)
![Postgres](https://img.shields.io/badge/Postgres-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

## Live Demo

| Service | URL |
| --- | --- |
| Frontend | https://kash-an-expense-tracker.vercel.app |
| Backend | https://kash-wq9x.onrender.com |
| Backend Health Check | https://kash-wq9x.onrender.com/api/v1/health |

## Project Overview

- The backend is the core assignment API built with Go + Beego.
- The frontend is a bonus UI built with Next.js and TypeScript.
- CSV storage remains the default for local and trainer review.
- Postgres is used for production deployment (hosted on Neon).
- Voice input is included in the frontend.

## Repository Structure

```text
.
├── backend/          # Go + Beego REST API
├── frontend/         # Next.js frontend
├── docs/             # Project images/screenshots
└── README.md         # Root project guide
```

## Local Development

- Backend setup, storage modes, Swagger docs, Docker usage, and tests are documented in [backend/README.md](backend/README.md).
- Frontend setup and UI details are documented in [frontend/README.md](frontend/README.md).
- When running the frontend locally, set `NEXT_PUBLIC_API_BASE_URL` to your backend URL (local or deployed).

## Documentation

- Backend details: [backend/README.md](backend/README.md)
- Frontend details: [frontend/README.md](frontend/README.md)
