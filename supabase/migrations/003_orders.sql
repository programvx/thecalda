-- 003_orders.sql
-- Phase 3: orders and shopping carts.
--
-- Carts and orders share one set of tables, distinguished by orders.type:
--   'cart'  — a draft owned by a user; status is always 'not_applicable'
--   'order' — a placed order; status follows the lifecycle enum
-- Checkout flips a cart row to type='order' in place.
--
-- Standard column convention (see 001_users.sql):
--   id          bigint identity primary key   -- internal
--   uid         uuid unique, auto-generated    -- public identifier
--   created_at  timestamptz, set on insert
--   updated_at  timestamptz, kept current by public.set_updated_at()
--
-- public.set_updated_at() is defined in 001_users.sql and reused here.

-- ---------------------------------------------------------------------------
-- Enums: order type (cart vs order) and the order lifecycle status.
-- ---------------------------------------------------------------------------
create type public.order_type as enum (
  'cart',
  'order'
);

create type public.order_status as enum (
  'not_applicable',
  'pending',
  'paid',
  'shipped',
  'delivered',
  'cancelled'
);

-- ---------------------------------------------------------------------------
-- addresses — postal addresses. An order references one for billing and one
-- for shipping (orders.billing_address_id / shipping_address_id); the two may
-- point at the same row. A cart references neither until checkout.
-- ---------------------------------------------------------------------------
create table public.addresses (
  id            bigint      generated always as identity primary key,
  uid           uuid        not null default gen_random_uuid() unique,
  created_at    timestamptz not null default now(),
  updated_at    timestamptz not null default now(),
  first_name    text        not null,
  last_name     text        not null,
  email         text        not null,
  phone         text,
  address_line1 text        not null,
  address_line2 text,
  postal_code   text        not null,
  city          text        not null,
  country       text        not null
);

-- ---------------------------------------------------------------------------
-- payment_methods — reference table of supported payment methods. An order
-- references one via orders.payment_method_id (set at checkout).
-- ---------------------------------------------------------------------------
create table public.payment_methods (
  id         bigint      generated always as identity primary key,
  uid        uuid        not null default gen_random_uuid() unique,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  code       text        not null unique,      -- stable key, e.g. 'card'
  name       text        not null,             -- display name
  is_active  boolean     not null default true -- offered at checkout?
);

-- Supported payment methods. Inserted in the migration (not seed.sql) so the
-- reference set is present on every environment, the same as currencies.
insert into public.payment_methods (code, name)
values
  ('card',             'Credit / Debit Card'),
  ('paypal',           'PayPal'),
  ('bank_transfer',    'Bank Transfer'),
  ('cash_on_delivery', 'Cash on Delivery');

-- ---------------------------------------------------------------------------
-- orders — carts and placed orders share this table (see header).
-- total_amount is denormalized: kept in sync with order_items by the
-- recalc_order_total() trigger defined further below.
-- The address / payment_method links are null for carts, set at checkout.
-- ---------------------------------------------------------------------------
create table public.orders (
  id                  bigint        generated always as identity primary key,
  uid                 uuid          not null default gen_random_uuid() unique,
  created_at          timestamptz   not null default now(),
  updated_at          timestamptz   not null default now(),
  user_id             bigint        not null references public.users (id) on delete cascade,
  type                public.order_type not null default 'cart',
  status              public.order_status not null default 'not_applicable',
  currency_id         bigint        not null references public.currencies (id),
  billing_address_id  bigint        references public.addresses (id) on delete restrict,
  shipping_address_id bigint        references public.addresses (id) on delete restrict,
  payment_method_id   bigint        references public.payment_methods (id) on delete restrict,
  total_amount        numeric(12,2) not null default 0 check (total_amount >= 0),
  order_number        text          unique,
  notes               text,
  placed_at           timestamptz,
  -- A cart is always 'not_applicable'; an order is never 'not_applicable'.
  constraint orders_type_status_chk check (
    (type = 'cart'  and status =  'not_applicable') or
    (type = 'order' and status <> 'not_applicable')
  )
);

-- ---------------------------------------------------------------------------
-- order_items — the lines of a cart/order. item_name and the unit_* prices
-- are snapshotted from the catalog so a line stays stable when the catalog
-- changes. unit_price_discounted and line_total are computed on every write.
-- ---------------------------------------------------------------------------
create table public.order_items (
  id                    bigint           generated always as identity primary key,
  uid                   uuid             not null default gen_random_uuid() unique,
  created_at            timestamptz      not null default now(),
  updated_at            timestamptz      not null default now(),
  order_id              bigint           not null references public.orders (id) on delete cascade,
  item_id               bigint           not null references public.items (id) on delete restrict,
  item_name             text             not null,                       -- snapshot of items.name
  quantity              integer          not null check (quantity > 0),
  unit_price            numeric(12,2)    not null check (unit_price >= 0), -- snapshot of items.price
  unit_discount         double precision not null default 0 check (unit_discount >= 0 and unit_discount <= 1), -- ratio: 0.2 = 20%
  unit_price_discounted numeric(12,2)    generated always as (round(unit_price * (1 - unit_discount::numeric), 2)) stored not null,
  line_total            numeric(12,2)    generated always as (round(quantity * unit_price * (1 - unit_discount::numeric), 2)) stored not null,
  unique (order_id, item_id)
);

-- ---------------------------------------------------------------------------
-- recalc_order_total() — keeps orders.total_amount in sync with its
-- order_items. A generated column cannot aggregate a child table, so the
-- total is denormalized and refreshed here on every order_items change.
-- ---------------------------------------------------------------------------
create or replace function public.recalc_order_total()
returns trigger
language plpgsql
as $$
begin
  if (tg_op = 'DELETE') then
    update public.orders
      set total_amount = coalesce(
        (select sum(line_total) from public.order_items where order_id = old.order_id), 0)
      where id = old.order_id;
    return old;
  end if;

  update public.orders
    set total_amount = coalesce(
      (select sum(line_total) from public.order_items where order_id = new.order_id), 0)
    where id = new.order_id;

  -- An UPDATE that moves a line to a different order must refresh both.
  if (tg_op = 'UPDATE' and new.order_id <> old.order_id) then
    update public.orders
      set total_amount = coalesce(
        (select sum(line_total) from public.order_items where order_id = old.order_id), 0)
      where id = old.order_id;
  end if;

  return new;
end;
$$;

-- ---------------------------------------------------------------------------
-- Triggers.
-- ---------------------------------------------------------------------------
create trigger addresses_set_updated_at
  before update on public.addresses
  for each row execute function public.set_updated_at();

create trigger payment_methods_set_updated_at
  before update on public.payment_methods
  for each row execute function public.set_updated_at();

create trigger orders_set_updated_at
  before update on public.orders
  for each row execute function public.set_updated_at();

create trigger order_items_set_updated_at
  before update on public.order_items
  for each row execute function public.set_updated_at();

create trigger order_items_recalc_total
  after insert or update or delete on public.order_items
  for each row execute function public.recalc_order_total();

-- ---------------------------------------------------------------------------
-- Indexes.
--
-- Already covered automatically (primary keys + UNIQUE constraints):
--   every id, every uid, payment_methods.code, orders.order_number, and
--   order_items (order_id, item_id) — the last also serves lookups filtered
--   by order_id alone (leftmost-prefix).
--
-- The indexes below cover the foreign keys not already indexed, the
-- one-cart-per-user guard, and the lifecycle filter path.
-- ---------------------------------------------------------------------------

-- FK: orders -> users; also lists a user's orders (WHERE user_id = ?).
create index orders_user_id_idx on public.orders (user_id);

-- At most one cart per user: a single type='cart' row per user_id.
create unique index orders_one_cart_per_user_idx
  on public.orders (user_id) where (type = 'cart');

-- FK: orders -> currencies.
create index orders_currency_id_idx on public.orders (currency_id);

-- FK: orders -> addresses, for billing and shipping.
create index orders_billing_address_id_idx on public.orders (billing_address_id);
create index orders_shipping_address_id_idx on public.orders (shipping_address_id);

-- FK: orders -> payment_methods.
create index orders_payment_method_id_idx on public.orders (payment_method_id);

-- Lifecycle filters, e.g. WHERE type = 'order' AND status = ?.
create index orders_type_status_idx on public.orders (type, status);

-- FK: order_items -> items, reverse lookup ("which orders contain this item").
create index order_items_item_id_idx on public.order_items (item_id);

-- ---------------------------------------------------------------------------
-- Row Level Security.
-- The Go API uses a direct (superuser) connection for which RLS is not
-- enforced, so all writes go through the API. These policies guard read
-- access via PostgREST: orders, their addresses and their items are private
-- to the owning user; payment_methods is public reference data.
-- ---------------------------------------------------------------------------
alter table public.orders          enable row level security;
alter table public.addresses       enable row level security;
alter table public.payment_methods enable row level security;
alter table public.order_items     enable row level security;

create policy orders_select_own
  on public.orders for select
  using (
    user_id in (select id from public.users where auth_user_id = auth.uid())
  );

create policy addresses_select_own
  on public.addresses for select
  using (
    id in (
      select unnest(array[o.billing_address_id, o.shipping_address_id])
      from public.orders o
      join public.users u on u.id = o.user_id
      where u.auth_user_id = auth.uid()
    )
  );

create policy order_items_select_own
  on public.order_items for select
  using (
    order_id in (
      select o.id
      from public.orders o
      join public.users u on u.id = o.user_id
      where u.auth_user_id = auth.uid()
    )
  );

-- payment_methods is public reference data, like currencies in 002.
create policy payment_methods_select
  on public.payment_methods for select using (true);
