---
name: run-project
description: >-
  Setup and run fm_expenses_be (Go expense-tracking backend, self-rolled
  auth, Postgres hosted on Supabase). Use when the user asks to run, start,
  or set up this project, or hits DB/JWT connection errors on first run.
---

# Run fm_expenses_be

Go backend, entrypoint at `cmd/api`. Postgres is hosted on Supabase, but Supabase is used **only**
as a DB host — auth (register/login/refresh/logout) is fully self-rolled on this backend, see
`internal/usecase/auth`. See `CLAUDE.md` at repo root for the full list of gotchas this skill is
built from.

## Steps

1. **Check `.env` exists.** If not, copy `.env.example` → `.env`.
2. **Entrypoint is `cmd/api`, not root.** Always run `go run ./cmd/api` (or `make run`).
   `go run .` fails with "no go files listed".
3. **If it fails to connect to the DB** (`dial unix .../.s.PGSQL.5432: no such file or directory`
   or similar): `DATABASE_URL` is missing/wrong. It must be a real Postgres connection string from
   **Supabase Dashboard → Project Settings → Database → Connect → Direct connection / Session
   pooler**, with `[YOUR-PASSWORD]` replaced by the actual DB password (strip the brackets).
   It is NOT the `sb_publishable_...` / `sb_secret_...` API keys — those don't go in `DATABASE_URL`.
4. **`JWT_SECRET`** is a secret **this backend picks itself** (`openssl rand -base64 32`) — it has
   nothing to do with Supabase's dashboard. Do not go looking for it in Supabase's API/JWT
   settings; that was only relevant to an earlier version of this backend that verified Supabase
   Auth tokens directly. This backend now issues and verifies its own tokens.
5. **Startup logs showing `SLOW SQL >= 200ms`** are normal (GORM's `AutoMigrate` schema checks
   plus network latency to Supabase) — not an error. Success looks like a final line:
   `fm_expenses_be listening on :8080`.

## Verifying it worked

```bash
curl -s localhost:8080/health   # expect "ok", no auth required

# Register (also logs in — returns a token pair immediately)
curl -s -X POST localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234!","full_name":"Test"}'
# -> grab data.tokens.access_token from the response

curl -s localhost:8080/api/v1/users/me -H "Authorization: Bearer <access_token>"
```

Routes under `/api/v1/*` (except `/api/v1/auth/*`) require `Authorization: Bearer <access_token>`.
Full auth flow docs: `docs/features/01-authentication.md`.
