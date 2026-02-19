# Exchange Travel Planner

Mobile-first web app for exchange trip planning with a Go backend API.

## Stack
- Frontend: Next.js 14 + TypeScript + Tailwind CSS
- Backend: Go (net/http), PostgreSQL (pgx/v5) with in-memory fallback
- Migrations: Flyway (via Docker)
- Local DB: Docker Compose (PostgreSQL 16)

## Local Development Setup

### 1. Start PostgreSQL + Run Migrations

```bash
# Install Docker runtime (Colima) if not already installed
brew install colima docker docker-compose
colima start

# Start Postgres and run Flyway migrations
docker compose up -d
```

This starts PostgreSQL on `localhost:5432` and automatically runs all migrations (schema + seed data).

### 2. Start Backend

```bash
cd backend

# Use local Postgres (default if DATABASE_URL not set)
DATABASE_URL="postgres://exchange:exchange_local@localhost:5432/exchange_dev?sslmode=disable" go run ./cmd/server

# Or without DATABASE_URL to use in-memory store (no DB needed)
go run ./cmd/server
```

Backend runs on `http://localhost:8080`.

### 3. Start Frontend

```bash
npm install
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

Frontend runs on `http://localhost:3000`.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Backend port | `8080` |
| `DATABASE_URL` | Postgres connection string | Falls back to in-memory store |
| `AUTH_DISABLED` | Set `true` to skip JWT auth (local dev) | `false` |
| `SUPABASE_JWT_SECRET` | Supabase JWT secret for token verification | Required when auth enabled |
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL | Skip auth if unset |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | Supabase anon key | Skip auth if unset |

### Authentication

The app uses **Supabase Auth** for authentication. For local development without Supabase:

1. Set `AUTH_DISABLED=true` when running the backend — all API requests use a default "demo-user"
2. Leave `NEXT_PUBLIC_SUPABASE_URL` and `NEXT_PUBLIC_SUPABASE_ANON_KEY` unset — the frontend skips auth

For production with Supabase:

1. Create a Supabase project and get your JWT secret, URL, and anon key
2. Set all env vars (see `.env.example`)
3. The backend validates JWT tokens on all `/api/*` routes (except `/health`)
4. The frontend provides login/signup UI and attaches Bearer tokens to API calls

## Database Migrations

Migrations live in `backend/migrations/` using Flyway naming (`V1__description.sql`).

**Local:** Migrations run automatically via `docker compose up`.

**Dev/Supabase:** Set `DATABASE_URL` as a JDBC string and run:
```bash
docker compose -f docker-compose.dev.yml run --rm flyway
```

**CI:** Pushing migration changes to `main` triggers the GitHub Action which runs Flyway against Supabase (requires `SUPABASE_DATABASE_URL` secret).

## Implemented MVP Areas
- Weekend Trip Optimizer (`POST /api/trips/optimize`)
- Academic Travel Windows (`GET /api/travel-windows`)
- Budget Tracker + Forecast (`/api/budget/entries`, `/api/budget/forecast`)
- Transport/Stay search adapters (`/api/search/transport`, `/api/search/stays`)
- Study-travel conflict checks (`POST /api/conflicts/evaluate`)
- Group trip share flow (`GET /api/trips/:id`, `POST /api/trips/:id/share`)
- Mobile-friendly screens for Home, Calendar, Discover, Budget, Group, Settings, Trip Detail
- PWA manifest and install metadata
