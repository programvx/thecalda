-- 002_catalog.sql
-- Phase 2: product catalog — currencies, catalogs, item stock, items, item media, and item properties.
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
-- item_stock — per-item stock level and availability status.
-- Linked from items.stock_id (defined below).
-- ---------------------------------------------------------------------------
create type public.item_stock_status as enum (
  'in_stock',
  'low_stock',
  'out_of_stock',
  'discontinued'
);

create table public.item_stock (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  quantity   integer     not null default 0 check (quantity >= 0),
  status     public.item_stock_status not null default 'out_of_stock'
);

-- ---------------------------------------------------------------------------
-- items — actual e-commerce products.
-- price_discounted is computed from price and discount on every write; it
-- cannot be inserted or updated directly.
-- ---------------------------------------------------------------------------
create table public.items (
  id               bigint generated always as identity primary key,
  uid              uuid          not null default gen_random_uuid() unique,
  created_at       timestamptz   not null default now(),
  updated_at       timestamptz   not null default now(),
  slug             text          not null unique,
  sku              text          unique,
  name             text          not null,
  description      text,
  price            numeric(12,2) not null check (price >= 0),
  price_discounted numeric(12,2) generated always as (round(price * (1 - discount::numeric), 2)) stored not null,
  discount         double        precision not null default 0 check (discount >= 0 and discount <= 1), -- ratio: 0.2 = 20%
  currency_id      bigint        not null references public.currencies (id),
  stock_id         bigint        references public.item_stock (id),
  is_active        boolean       not null default true
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
  position   integer     not null default 0,
  item_id    bigint      not null references public.items (id)    on delete cascade,
  unique (catalog_id, item_id)
);

-- ---------------------------------------------------------------------------
-- item_medias — media (image, video, 3D, ...) attached to an item; the
-- lowest position is the primary media item.
-- ---------------------------------------------------------------------------
create type public.item_media_type as enum (
  'image',
  'video',
  '3d',
  'document'
);

create table public.item_medias (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  media_type public.item_media_type not null default 'image',
  url        text        not null,
  alt        text,
  position   integer     not null default 0,
  item_id    bigint      not null references public.items (id) on delete cascade
);

-- ---------------------------------------------------------------------------
-- item_properties — arbitrary label/value attributes for an item.
-- One item can have many properties.
-- ---------------------------------------------------------------------------
create table public.item_properties (
  id         bigint generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  label      text        not null,
  value      text        not null,
  item_id    bigint      not null references public.items (id) on delete cascade
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

create trigger item_stock_set_updated_at
  before update on public.item_stock
  for each row execute function public.set_updated_at();

create trigger items_set_updated_at
  before update on public.items
  for each row execute function public.set_updated_at();

create trigger catalog_items_set_updated_at
  before update on public.catalog_items
  for each row execute function public.set_updated_at();

create trigger item_medias_set_updated_at
  before update on public.item_medias
  for each row execute function public.set_updated_at();

create trigger item_properties_set_updated_at
  before update on public.item_properties
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

-- FK: items -> item_stock (joins, and cheap item_stock updates/deletes).
create index items_stock_id_idx on public.items (stock_id);

-- FK: catalog_items -> items, reverse lookup ("which catalogs hold this item").
create index catalog_items_item_id_idx on public.catalog_items (item_id);

-- Render a catalog's items in order: WHERE catalog_id = ? ORDER BY position.
create index catalog_items_catalog_position_idx
  on public.catalog_items (catalog_id, position);

-- FK: item_medias -> items; the item_id prefix also serves ordered retrieval
-- of an item's images: WHERE item_id = ? ORDER BY position.
create index item_medias_item_position_idx
  on public.item_medias (item_id, position);

-- FK: item_properties -> items; fetch an item's properties (WHERE item_id = ?).
create index item_properties_item_id_idx
  on public.item_properties (item_id);

-- ---------------------------------------------------------------------------
-- Row Level Security — public read of active content.
-- The Go API uses a direct (superuser) connection for which RLS is not
-- enforced, so all writes go through the API. These policies guard read
-- access via PostgREST. No insert/update/delete policies are defined.
-- ---------------------------------------------------------------------------
alter table public.currencies      enable row level security;
alter table public.catalogs        enable row level security;
alter table public.item_stock      enable row level security;
alter table public.items           enable row level security;
alter table public.catalog_items   enable row level security;
alter table public.item_medias     enable row level security;
alter table public.item_properties enable row level security;

create policy currencies_select
  on public.currencies for select using (true);

create policy catalogs_select_active
  on public.catalogs for select using (is_active);

create policy item_stock_select
  on public.item_stock for select using (true);

create policy items_select_active
  on public.items for select using (is_active);

create policy catalog_items_select
  on public.catalog_items for select using (true);

create policy item_medias_select
  on public.item_medias for select using (true);

create policy item_properties_select
  on public.item_properties for select using (true);
