-- 001_users.sql
-- Phase 1: application users linked to Supabase Auth.
--
-- Standard column convention (every app table):
--   id          bigint identity primary key   -- internal
--   uid         uuid unique, auto-generated    -- public identifier
--   created_at  timestamptz, set on insert
--   updated_at  timestamptz, kept current by the set_updated_at() trigger

-- ---------------------------------------------------------------------------
-- set_updated_at() — shared BEFORE UPDATE trigger function.
-- Attach to every table that carries an updated_at column.
-- ---------------------------------------------------------------------------
create or replace function public.set_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

-- ---------------------------------------------------------------------------
-- public.users — one profile row per auth.users entry.
-- auth_user_id links to Supabase Auth; id/uid follow the project convention.
-- ---------------------------------------------------------------------------
create table public.users (
  id           bigint generated always as identity primary key,
  uid          uuid        not null default gen_random_uuid() unique,
  created_at   timestamptz not null default now(),
  updated_at   timestamptz not null default now(),
  email        text,
  full_name    text,
  auth_user_id uuid        not null unique references auth.users (id) on delete cascade
);

create trigger users_set_updated_at
  before update on public.users
  for each row
  execute function public.set_updated_at();

-- ---------------------------------------------------------------------------
-- handle_new_user() — on signup, create the matching public.users row.
-- This is the "link your users table" requirement: every Supabase Auth user
-- automatically gets a public.users profile.
-- ---------------------------------------------------------------------------
create or replace function public.handle_new_user()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
begin
  insert into public.users (auth_user_id, email, full_name)
  values (
    new.id,
    new.email,
    coalesce(
      new.raw_user_meta_data ->> 'full_name',
      new.raw_user_meta_data ->> 'name'
    )
  );
  return new;
end;
$$;

create trigger on_auth_user_created
  after insert on auth.users
  for each row
  execute function public.handle_new_user();

-- ---------------------------------------------------------------------------
-- Row Level Security — defense in depth.
-- The Go API uses a direct (superuser) pgx connection, for which RLS is not
-- enforced; these policies guard any access that goes through PostgREST with
-- an end-user JWT.
-- ---------------------------------------------------------------------------
alter table public.users enable row level security;

create policy users_select_own
  on public.users
  for select
  using (auth.uid() = auth_user_id);

create policy users_update_own
  on public.users
  for update
  using (auth.uid() = auth_user_id)
  with check (auth.uid() = auth_user_id);
