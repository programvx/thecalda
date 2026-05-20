# thecalda

A simple e-commerce application built on **Supabase** (PostgreSQL + Auth), a
**Go / Gin** backend, and a **Next.js** (App Router, SSR) frontend.

The project is built in phases. **Phase 1 — authentication — is complete.**

## Phase 1 scope

- Local Supabase stack (Postgres + Auth) managed by the Supabase CLI.
- `public.users` table linked to `auth.users`, populated automatically on
  signup by a database trigger.
- Go/Gin backend with Supabase JWT verification and a `GET /api/me` endpoint.
- Next.js frontend: email/password sign-up, sign-in and sign-out, plus an
  SSR account page that renders the profile returned by the backend.

Items, orders and item-change tracking are **out of scope for Phase 1** and
are planned for later phases.

## Architecture

```
Browser ──signup / login──▶ Supabase Auth ──┐
   │                                        │ trigger: handle_new_user()
   │  session cookie (@supabase/ssr)         ▼
   ▼                                  public.users  ◀── auth.users
Next.js (SSR) ──Bearer access token──▶ Go / Gin API ──pgx──▶ Postgres
                                        verifies JWT via JWKS (ES256)
```

The frontend talks to **Supabase Auth directly** for sign-up/sign-in only;
everything else goes through the **Go API**. Change tracking (later phases)
happens in the database.

## Tech stack

| Layer    | Choice                                                          |
| -------- | --------------------------------------------------------------- |
| Database | PostgreSQL via Supabase (local stack, Supabase CLI)             |
| Auth     | Supabase Auth — backend verifies the JWT against the JWKS (ES256) |
| Backend  | Go + Gin, GORM (`gorm.io/gorm`), `zap`; Clean Architecture       |
| Frontend | Next.js 16 (App Router, SSR), React 19, Tailwind v4, `@supabase/ssr` |

## Repository layout

```
thecalda/
├── supabase/
│   ├── config.toml
│   └── migrations/001_users.sql      # public.users + triggers + RLS
├── backend/                          # Go / Gin API
│   ├── cmd/api/                       # entrypoint
│   └── internal/
│       ├── settings/ core/ db/ db/crud/
│       ├── model/ services/ handlers/
│       ├── middlewares/               # cors, logger, security, auth (JWKS)
│       └── routers/
└── frontend/                         # Next.js app
    ├── app/                           # routes (pages + route handlers)
    ├── components/
    ├── lib/                           # supabase clients, api client, types
    └── proxy.ts                       # session refresh (Next.js 16 "proxy")
```

## Prerequisites

- Docker (for the Supabase local stack)
- [Supabase CLI](https://supabase.com/docs/guides/cli)
- Go 1.26+
- Node.js 22+ and [pnpm](https://pnpm.io)

## Running locally

Start the three pieces in separate terminals.

### 1. Supabase

```bash
supabase start          # boots Postgres + Auth in Docker, applies migrations
supabase status         # prints URLs and keys
```

### 2. Backend

```bash
cd backend
cp .env.example .env    # defaults already match the local Supabase stack
make run                # or: go run ./cmd/api
```

API: <http://localhost:8080> — check `GET /health`.

### 3. Frontend

```bash
cd frontend
cp .env.local.example .env.local
# set NEXT_PUBLIC_SUPABASE_ANON_KEY to the "Publishable" key from `supabase status`
pnpm install
pnpm dev
```

App: <http://localhost:3000>

## Standard table columns

Every application table follows this convention:

| Column       | Definition                                          |
| ------------ | --------------------------------------------------- |
| `id`         | `bigint generated always as identity primary key`   |
| `uid`        | `uuid not null default gen_random_uuid() unique`    |
| `created_at` | `timestamptz not null default now()`                |
| `updated_at` | `timestamptz not null default now()` (trigger-kept) |

`id` is internal; `uid` is the public identifier used in API paths and URLs.

## API

| Method | Path       | Auth          | Description                       |
| ------ | ---------- | ------------- | --------------------------------- |
| GET    | `/health`  | public        | Liveness + database check         |
| GET    | `/api/me`  | Supabase JWT  | Authenticated caller's profile    |

Responses use a consistent envelope: `{ "data": ... }` or
`{ "error": { "code", "message", "details" } }`.

## Tests

```bash
cd backend && go test ./...   # handler unit tests
```

The full Phase 1 flow (sign up → SSR account page from `/api/me` → sign out →
auth gate → sign in) was verified end to end with a headless browser.

## Notes

- The backend listens on **8080** (configurable via `PORT` in `backend/.env`).
- Local Supabase issues asymmetric (ES256) JWTs; the backend verifies them
  against `<SUPABASE_URL>/auth/v1/.well-known/jwks.json`.
- Database migrations are owned by the Supabase CLI (`supabase/migrations/`).
