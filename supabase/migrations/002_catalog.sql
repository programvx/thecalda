-- 002_catalog.sql
-- Phase 2: product catalog — currencies, catalogs, items, and item media.
--
-- Standard column convention (see 001_users.sql):
--   id          bigint identity primary key   -- internal
--   uid         uuid unique, auto-generated    -- public identifier
--   created_at  timestamptz, set on insert
--   updated_at  timestamptz, kept current by public.set_updated_at()
--
-- public.set_updated_at() is defined in 001_users.sql and reused here.

-- ---------------------------------------------------------------------------
-- currencies — reference table of supported currencies.
-- ---------------------------------------------------------------------------
create table public.currencies (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  code       text        not null unique check (char_length(code) = 3), -- ISO 4217, e.g. 'EUR'
  name       text        not null,                                      -- e.g. 'Euro'
  symbol     text                                                       -- e.g. '€'
);

-- Default currency. Inserted in the migration (not seed.sql) so it is present
-- after `migration up` on any environment, since items.currency_id is NOT NULL.
insert into public.currencies (code, name, symbol)
values ('EUR', 'Euro', '€');

-- ---------------------------------------------------------------------------
-- catalogs — a named grouping/collection of items.
-- ---------------------------------------------------------------------------
create table public.catalogs (
  id          bigint generated always as identity primary key,
  uid         uuid        not null default gen_random_uuid() unique,
  created_at  timestamptz not null default now(),
  updated_at  timestamptz not null default now(),
  slug        text        not null unique,
  name        text        not null,
  description text,
  is_active   boolean     not null default true
);

-- ---------------------------------------------------------------------------
-- items — actual e-commerce products.
-- price_discounted is computed from price and discount on every write; it
-- cannot be inserted or updated directly.
-- ---------------------------------------------------------------------------
create table public.items (
  id               bigint generated always as identity primary key,
  uid              uuid        not null default gen_random_uuid() unique,
  created_at       timestamptz not null default now(),
  updated_at       timestamptz not null default now(),
  slug             text        not null unique,
  sku              text        unique,
  name             text        not null,
  description      text,
  price            numeric(12,2)    not null check (price >= 0),
  discount         double precision not null default 0
                                    check (discount >= 0 and discount <= 1), -- ratio: 0.2 = 20%
  price_discounted numeric(12,2)
                     generated always as (round(price * (1 - discount::numeric), 2)) stored
                     not null,
  currency_id      bigint      not null references public.currencies (id),
  stock_quantity   integer     not null default 0 check (stock_quantity >= 0),
  is_active        boolean     not null default true
);

-- ---------------------------------------------------------------------------
-- catalog_items — many-to-many between catalogs and items, with ordering.
-- ---------------------------------------------------------------------------
create table public.catalog_items (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  catalog_id bigint      not null references public.catalogs (id) on delete cascade,
  item_id    bigint      not null references public.items (id)    on delete cascade,
  position   integer     not null default 0,
  unique (catalog_id, item_id)
);

-- ---------------------------------------------------------------------------
-- item_images — media for an item; lowest position is the primary image.
-- ---------------------------------------------------------------------------
create table public.item_images (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  item_id    bigint      not null references public.items (id) on delete cascade,
  url        text        not null,
  alt        text,
  position   integer     not null default 0
);

-- ---------------------------------------------------------------------------
-- updated_at triggers.
-- ---------------------------------------------------------------------------
create trigger currencies_set_updated_at
  before update on public.currencies
  for each row execute function public.set_updated_at();

create trigger catalogs_set_updated_at
  before update on public.catalogs
  for each row execute function public.set_updated_at();

create trigger items_set_updated_at
  before update on public.items
  for each row execute function public.set_updated_at();

create trigger catalog_items_set_updated_at
  before update on public.catalog_items
  for each row execute function public.set_updated_at();

create trigger item_images_set_updated_at
  before update on public.item_images
  for each row execute function public.set_updated_at();

-- ---------------------------------------------------------------------------
-- Indexes.
--
-- Already covered automatically (primary keys + UNIQUE constraints):
--   every id, every uid, currencies.code, catalogs.slug, items.slug,
--   items.sku, and catalog_items (catalog_id, item_id) — the last also
--   serves lookups filtered by catalog_id alone (leftmost-prefix).
--
-- The indexes below cover the foreign keys not already indexed, plus the
-- ordered-retrieval access paths. Low-cardinality booleans (is_active) are
-- intentionally left unindexed; add partial indexes if/when query patterns
-- justify them.
-- ---------------------------------------------------------------------------

-- FK: items -> currencies (joins, and cheap currency updates/deletes).
create index items_currency_id_idx on public.items (currency_id);

-- FK: catalog_items -> items, reverse lookup ("which catalogs hold this item").
create index catalog_items_item_id_idx on public.catalog_items (item_id);

-- Render a catalog's items in order: WHERE catalog_id = ? ORDER BY position.
create index catalog_items_catalog_position_idx
  on public.catalog_items (catalog_id, position);

-- FK: item_images -> items; the item_id prefix also serves ordered retrieval
-- of an item's images: WHERE item_id = ? ORDER BY position.
create index item_images_item_position_idx
  on public.item_images (item_id, position);

-- ---------------------------------------------------------------------------
-- Row Level Security — public read of active content.
-- The Go API uses a direct (superuser) connection for which RLS is not
-- enforced, so all writes go through the API. These policies guard read
-- access via PostgREST. No insert/update/delete policies are defined.
-- ---------------------------------------------------------------------------
alter table public.currencies    enable row level security;
alter table public.catalogs      enable row level security;
alter table public.items         enable row level security;
alter table public.catalog_items enable row level security;
alter table public.item_images   enable row level security;

create policy currencies_select
  on public.currencies for select using (true);

create policy catalogs_select_active
  on public.catalogs for select using (is_active);

create policy items_select_active
  on public.items for select using (is_active);

create policy catalog_items_select
  on public.catalog_items for select using (true);

create policy item_images_select
  on public.item_images for select using (true);
