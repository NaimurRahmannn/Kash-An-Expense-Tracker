<div align="center">

# Expense Tracker Frontend

Modern Next.js App Router frontend for a Personal Expense Tracker backend built with Go + Beego.

![Next.js](https://img.shields.io/badge/Next.js-App%20Router-black)
![TypeScript](https://img.shields.io/badge/TypeScript-Strict-blue)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-UI-38bdf8)
![Auth](https://img.shields.io/badge/Auth-localStorage-orange)
![Status](https://img.shields.io/badge/Status-Assignment%20Ready-22c55e)

</div>

## Overview

This frontend provides the full user-facing experience for the Personal Expense Tracker assignment backend. It supports authentication, protected application pages, dashboard summaries, expense management, and a guided Voice Input bonus feature.

The backend uses Go + Beego and stores data in CSV as part of the assignment requirement. The frontend is designed around that backend contract instead of pretending unsupported production features exist.

## Table Of Contents

- [Tech Stack](#tech-stack)
- [Feature Matrix](#feature-matrix)
- [Backend Contract](#backend-contract)
- [Getting Started](#getting-started)
- [Environment](#environment)
- [Engineering Decisions](#engineering-decisions)
- [Voice Input Format](#voice-input-format)
- [Feature Design Rationale](#feature-design-rationale)
- [Current Status](#current-status)
- [Limitations](#limitations)
- [Future Improvements](#future-improvements)
- [API Routes Used](#api-routes-used)

## Tech Stack

| Area | Choice |
| --- | --- |
| Framework | Next.js App Router |
| Language | TypeScript |
| Styling | Tailwind CSS |
| Icons | lucide-react |
| Auth persistence | `localStorage` |
| Voice input | Browser Web Speech API |
| Backend | Go + Beego |

## Feature Matrix

| Feature | Status | Notes |
| --- | --- | --- |
| Register | Complete | Calls backend auth API |
| Login | Complete | Stores returned user in `localStorage` |
| Protected routes | Complete | App pages require login |
| Dashboard | Complete | Uses backend summary API |
| Expenses list | Complete | Backend filtering, sorting, simple pagination |
| Add expense | Complete | Uses `POST /expenses` |
| Edit expense | Complete | Uses `GET` and `PUT /expenses/:id` |
| Delete expense | Complete | Confirmation modal before DELETE |
| Profile | Present | Displays session/account context |
| Settings | Present | Shows API/session configuration context |
| Voice Input | Complete | Strict guided parser with editable review |
| Charts | Not implemented | Future improvement |

## Backend Contract

The backend returns a simple user object after login:

```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

The frontend stores that object in `localStorage` under:

```txt
expense_tracker_user
```

Protected expense routes require:

```txt
X-User-ID: <user_id>
```

Expense create/update payloads only include backend-supported fields:

```json
{
  "title": "Lunch",
  "amount": 350.5,
  "category": "Food",
  "note": "Team lunch",
  "expense_date": "2025-06-10"
}
```

## Getting Started

Install dependencies:

```bash
npm install
```

Run the development server:

```bash
npm run dev
```

Open:

```txt
http://localhost:3000
```

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

## Engineering Decisions

### 1. Strict Voice Input Format Instead Of Free-Form NLP

The Voice Input feature intentionally uses a strict guided format:

```txt
Title Amount Category Date
```

Examples:

- Lunch 350 food today
- Bus fare 80 transport yesterday
- Medicine 500 healthcare 2025-06-10
- Electricity bill 1200 utilities today

This is an intentional engineering tradeoff.

| Concern | Decision |
| --- | --- |
| Free-form speech parsing | Avoided because it is unreliable and edge-case heavy |
| Backend payload | Kept strict because the backend expects structured expense fields |
| Parser behavior | Deterministic, testable, and easier to maintain |
| Dependencies | No heavy NLP libraries |
| User safety | Voice results are never auto-saved |

The app always shows an editable confirmation form before sending voice-created data to the backend.

Parser safety rules:

- The first valid positive number becomes the amount.
- Words before the amount become the title.
- Category and date are parsed from words after the amount.
- Date must be `today`, `yesterday`, or `YYYY-MM-DD`.
- Category must match the allowed backend categories.
- Invalid or missing fields show clear errors.
- Empty transcript returns `No speech detected`.

### 2. LocalStorage Auth Because Backend Returns user_id

The backend assignment uses simple login and returns `user_id`, `name`, and `email`. Protected expense APIs require `X-User-ID`, so the frontend stores the login result in `localStorage`.

This keeps the frontend aligned with the backend assignment. JWT or server-managed session authentication would be better for production, but it is outside the current backend requirement.

This is suitable for an assignment or demo app. It is not production-grade authentication.

### 3. Dashboard Uses Summary API

The dashboard does not calculate all totals manually from frontend state. It uses the backend summary endpoint.

This keeps spending calculations and business rules centralized in the backend, while keeping the dashboard lightweight and consistent with backend results.

### 4. Expenses Page Uses Backend Filtering And Sorting

Category filter, date range filter, sorting, and pagination are sent to the backend.

Frontend search is limited to the currently loaded page because the backend does not expose full-text search. This avoids overclaiming frontend behavior that the API does not support.

### 5. Simple Pagination Because Backend Has No Total Count

The backend expense list returns only the current page array. It does not return `total_count` or `total_pages`.

The frontend therefore uses simple pagination:

| Control | Rule |
| --- | --- |
| Previous | Enabled when `page > 1` |
| Next | Enabled when returned item count equals `limit` |

This avoids displaying fake total pages.

### 6. Confirmation Before Delete

Delete is destructive. The UI asks for confirmation before calling the DELETE API, then refreshes the current list after successful deletion.

### 7. Form Validation On Frontend And Backend

Frontend validation provides faster feedback and better UX. The backend remains the source of truth for validation and data integrity.

### 8. AppShell And Protected Routes

Dashboard, expenses, voice, profile, and settings use a shared `AppShell`.

`AppShell` keeps sidebar/topbar layout consistent. `ProtectedRoute` prevents unauthenticated users from accessing app pages. Login and register remain public.

## Voice Input Format

Voice Input uses this strict guided format:

```txt
Title Amount Category Date
```

Examples:

- Lunch 350 food today
- Bus fare 80 transport yesterday
- Medicine 500 healthcare 2025-06-10
- Electricity bill 1200 utilities today
- Health care bill 900 health care today

Allowed categories:

```txt
Food, Transport, Housing, Entertainment, Shopping, Healthcare, Education, Utilities, Other
```

Allowed dates:

```txt
today, yesterday, YYYY-MM-DD
```

This design guides the user toward a predictable structure instead of trying to understand every possible natural sentence. That makes the parser reliable, reduces bugs, and keeps the feature easy to maintain.

### Parser Examples

| Spoken Input | Parsed Title | Amount | Category | Date |
| --- | --- | ---: | --- | --- |
| Lunch 350 food today | Lunch | 350 | Food | Today |
| Bus fare 80 transport yesterday | Bus fare | 80 | Transport | Yesterday |
| Medicine 500 healthcare 2025-06-10 | Medicine | 500 | Healthcare | 2025-06-10 |
| Health care bill 900 health care today | Health care bill | 900 | Healthcare | Today |

### Invalid Examples

| Input | Result |
| --- | --- |
| Lunch 350 food | Date missing or invalid |
| Lunch 350 today | Category not recognized |
| Lunch food today | Amount not found |
| Lunch 350 food tomorrow | Invalid date |
| Empty speech | No speech detected |

Parsed data should always be reviewed before saving.

## Feature Design Rationale

| Feature | Rationale |
| --- | --- |
| Authentication | Kept simple to match the backend assignment; stores `user_id` for `X-User-ID` |
| Dashboard | Uses summary API for accurate totals and keeps overview separate from management |
| Expenses page | Dedicated place for list, filters, sorting, pagination, edit, and delete |
| Add/Edit forms | Sends only `title`, `amount`, `category`, `note`, and `expense_date` |
| Delete | Uses confirmation modal to prevent accidental data loss |
| Voice Input | Bonus feature built on existing create API with strict parsing and manual review |
| Settings/Profile | Shows session/API context without faking unsupported backend profile updates |

## Current Status

Completed frontend parts:

- App Router foundation and layout shell
- Login/register API integration
- Local user persistence
- Protected app routes
- Dashboard summary API integration
- Expenses list API integration
- Add expense
- Edit expense
- Delete expense with confirmation
- Voice Input bonus feature
- Profile and settings pages

Chart library integration is not implemented yet.

## Limitations

| Limitation | Reason |
| --- | --- |
| `localStorage` auth | Matches assignment backend, not production-grade auth |
| No JWT/session token | Backend does not provide token/session auth |
| No profile update | Backend does not expose profile update endpoint |
| No total page count | Backend list API returns only an array |
| Strict voice parsing | Chosen intentionally for reliability |
| Browser-dependent voice input | Depends on Web Speech API support |
| Limited frontend search | Backend full-text search is not implemented |

## Future Improvements

- JWT authentication
- Backend pagination metadata
- Full-text search
- Better analytics charts
- Dark mode
- CSV export
- Improved voice command support
- Toast notification system
- Profile update support if backend adds an endpoint
- Deployment guide

## API Routes Used

| Method | Route | Purpose |
| --- | --- | --- |
| POST | `/api/v1/auth/register` | Register user |
| POST | `/api/v1/auth/login` | Login user |
| GET | `/api/v1/expenses` | List expenses |
| POST | `/api/v1/expenses` | Create expense |
| GET | `/api/v1/expenses/:id` | Get one expense |
| PUT | `/api/v1/expenses/:id` | Update expense |
| DELETE | `/api/v1/expenses/:id` | Delete expense |
| GET | `/api/v1/expenses/summary` | Dashboard summary |
