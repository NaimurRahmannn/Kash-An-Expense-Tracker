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

Frontend Part 2 completed:

- Login API integration
- Register API integration
- Auth storage helper
- Local user persistence
- Basic logout helper

Frontend Part 3 completed:

- Auth hook added
- Protected routes added
- AppShell pages require login
- Sidebar/Topbar use stored user
- Logout added
- Login/register redirect authenticated users to dashboard

Frontend Part 4 completed:

- Dashboard connected to backend summary API
- Recent expenses loaded from backend
- Date range summary support
- Dashboard loading and error states
- Category breakdown from real API data

Frontend Part 5 completed:

- Expenses page connected to backend list API
- Category/date filtering
- Sorting
- Simple pagination
- Loading/error/empty states
- Edit/delete action placeholders

Create/edit/delete expense API calls, chart library integration, and voice input are not implemented yet.

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

Create `.env.local` from `.env.example`:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

## Auth Flow Notes

- Register calls `POST /api/v1/auth/register`.
- Login calls `POST /api/v1/auth/login`.
- Login stores the backend `user_id`, `name`, and `email` in localStorage under `expense_tracker_user`.
- Protected pages read localStorage using `useAuth`.
- Dashboard uses `X-User-ID` from the stored login user.
- Dashboard and expenses list requests use the stored `user_id` with the `X-User-ID` header.
- User must login before dashboard data can load.

## Expenses Pagination Note

The backend list response currently returns only the current page data array. It does not return total count or total pages, so the frontend enables the Next button only when the returned data length matches the requested page limit.

## Planned API Routes

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/expenses`
- `POST /api/v1/expenses`
- `GET /api/v1/expenses/:id`
- `PUT /api/v1/expenses/:id`
- `DELETE /api/v1/expenses/:id`
- `GET /api/v1/expenses/summary`
