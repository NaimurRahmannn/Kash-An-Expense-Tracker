# Expense Tracker Frontend

Next.js App Router frontend for the Personal Expense Tracker backend built with Go + Beego.

## Frontend Status

Frontend Part 1 completed:

- Next.js App Router setup
- Base layout
- Sidebar/topbar
- Placeholder dashboard
- Placeholder expenses page
- Placeholder add/edit expense pages
- Placeholder auth pages
- API client foundation
- Shared types

No real API calls, authentication, CRUD, dashboard charts, voice input, localStorage, or route protection are implemented yet.

## Getting Started

Install dependencies:

```bash
npm install
```

Run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

Build the app:

```bash
npm run build
```

Run linting:

```bash
npm run lint
```

## Environment

Create `.env.local` from `.env.example` when API wiring begins:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

## Planned API Routes

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/expenses`
- `POST /api/v1/expenses`
- `GET /api/v1/expenses/:id`
- `PUT /api/v1/expenses/:id`
- `DELETE /api/v1/expenses/:id`
- `GET /api/v1/expenses/summary`
