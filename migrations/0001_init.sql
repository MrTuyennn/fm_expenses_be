-- Reference schema matching the GORM AutoMigrate models.
-- Run this manually if you prefer explicit migrations over AutoMigrate.
--
-- Auth is fully self-rolled in this backend (see internal/usecase/auth) —
-- Postgres here is just a Postgres host (e.g. Supabase's), not an identity
-- provider. There is deliberately no `auth.users` FK and no Row Level
-- Security: RLS via `auth.uid()` only makes sense when Postgres itself can
-- verify the caller's JWT (as Supabase Auth + PostgREST do). This backend's
-- own tokens are opaque to Postgres, so the Go backend is the sole
-- authorization gatekeeper — every query already scopes by user_id in
-- internal/repo/persistent/*.

create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text unique not null,
    password_hash text not null,
    full_name text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists refresh_tokens (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users (id) on delete cascade,
    token_hash text unique not null,
    expires_at timestamptz not null,
    revoked_at timestamptz,
    created_at timestamptz not null default now()
);
create index if not exists idx_refresh_tokens_user_id on refresh_tokens (user_id);

create table if not exists categories (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users (id) on delete cascade,
    name text not null,
    icon text,
    type text not null check (type in ('income', 'expense')),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index if not exists idx_categories_user_id on categories (user_id);

create table if not exists expenses (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users (id) on delete cascade,
    category_id uuid not null references categories (id) on delete restrict,
    amount numeric(14, 2) not null,
    currency varchar(3) not null default 'VND',
    note text,
    spent_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index if not exists idx_expenses_user_id on expenses (user_id);
create index if not exists idx_expenses_category_id on expenses (category_id);
