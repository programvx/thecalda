-- seed.sql — local development seed data.
-- Loaded after migrations on `supabase db reset`.

-- ---------------------------------------------------------------------------
-- catalogs — merchandising collections of items.
-- ---------------------------------------------------------------------------
insert into public.catalogs (slug, name, description, is_active)
values
  ('featured',     'Featured',     'Hand-picked products we love.',            true),
  ('new-arrivals', 'New Arrivals', 'The latest additions to the store.',       true),
  ('best-sellers', 'Best Sellers', 'Our most popular products.',               true),
  ('on-sale',      'On Sale',      'Discounted products and clearance deals.', true),
  ('holiday-2025', 'Holiday 2025', 'Seasonal holiday collection (archived).',  false)
on conflict (slug) do nothing;

-- ---------------------------------------------------------------------------
-- items — products. currency_id is resolved to the seeded EUR currency;
-- discount is a 0..1 ratio.
-- ---------------------------------------------------------------------------
insert into public.items (slug, sku, name, description, price, discount, currency_id)
select v.slug, v.sku, v.name, v.description, v.price, v.discount, eur.id
from (values
  ('wireless-headphones', 'WH-001', 'Wireless Headphones',
   'Over-ear Bluetooth headphones with active noise cancellation.', 89.99, 0.15),
  ('mechanical-keyboard', 'KB-001', 'Mechanical Keyboard',
   'Compact 75% mechanical keyboard with hot-swappable switches.', 129.00, 0),
  ('led-desk-lamp', 'DL-001', 'LED Desk Lamp',
   'Adjustable desk lamp with three brightness levels.', 34.50, 0),
  ('ceramic-mug', 'CM-001', 'Ceramic Coffee Mug',
   '350ml stoneware mug, microwave and dishwasher safe.', 12.99, 0),
  ('dotted-notebook', 'NB-001', 'Dotted Notebook',
   'A5 hardcover notebook with 192 dotted pages.', 8.50, 0),
  ('canvas-backpack', 'BP-001', 'Canvas Backpack',
   'Water-resistant canvas backpack with a padded laptop sleeve.', 64.00, 0.20),
  ('steel-water-bottle', 'WB-001', 'Stainless Steel Water Bottle',
   'Insulated 750ml bottle that keeps drinks cold for 24 hours.', 19.99, 0),
  ('aluminum-phone-stand', 'PS-001', 'Aluminum Phone Stand',
   'Foldable desk stand compatible with phones and tablets.', 15.00, 0.10)
) as v(slug, sku, name, description, price, discount)
cross join (select id from public.currencies where code = 'EUR') as eur
on conflict (slug) do nothing;

-- ---------------------------------------------------------------------------
-- catalog_items — places items into catalogs (items may appear in several).
-- catalog_id / item_id are resolved from slugs.
-- ---------------------------------------------------------------------------
insert into public.catalog_items (catalog_id, item_id, position)
select c.id, i.id, v.position
from (values
  ('featured',     'wireless-headphones',  0),
  ('featured',     'mechanical-keyboard',  1),
  ('featured',     'canvas-backpack',      2),
  ('featured',     'led-desk-lamp',        3),
  ('new-arrivals', 'dotted-notebook',      0),
  ('new-arrivals', 'aluminum-phone-stand', 1),
  ('new-arrivals', 'steel-water-bottle',   2),
  ('new-arrivals', 'ceramic-mug',          3),
  ('best-sellers', 'wireless-headphones',  0),
  ('best-sellers', 'ceramic-mug',          1),
  ('best-sellers', 'canvas-backpack',      2),
  ('on-sale',      'wireless-headphones',  0),
  ('on-sale',      'canvas-backpack',      1),
  ('on-sale',      'aluminum-phone-stand', 2)
) as v(catalog_slug, item_slug, position)
join public.catalogs c on c.slug = v.catalog_slug
join public.items i on i.slug = v.item_slug
on conflict (catalog_id, item_id) do nothing;

-- ---------------------------------------------------------------------------
-- item_stock — one stock record per item, linked via items.stock_id. Seeded
-- in a DO block because the link is set on items after each stock insert.
-- Runs once: skipped if item_stock already has rows.
-- ---------------------------------------------------------------------------
do $$
declare
  rec          record;
  new_stock_id bigint;
begin
  if exists (select 1 from public.item_stock) then
    return;
  end if;

  for rec in
    select * from (values
      ('wireless-headphones',  42,  'in_stock'),
      ('mechanical-keyboard',  18,  'in_stock'),
      ('led-desk-lamp',        7,   'low_stock'),
      ('ceramic-mug',          0,   'out_of_stock'),
      ('dotted-notebook',      120, 'in_stock'),
      ('canvas-backpack',      4,   'low_stock'),
      ('steel-water-bottle',   60,  'in_stock'),
      ('aluminum-phone-stand', 0,   'discontinued')
    ) as v(item_slug, quantity, status)
  loop
    insert into public.item_stock (quantity, status)
    values (rec.quantity, rec.status::public.item_stock_status)
    returning id into new_stock_id;

    update public.items
    set stock_id = new_stock_id
    where slug = rec.item_slug;
  end loop;
end $$;

-- ---------------------------------------------------------------------------
-- item_properties — label/value attributes per item (one item, many rows).
-- Runs once: skipped if item_properties already has rows.
-- ---------------------------------------------------------------------------
insert into public.item_properties (item_id, label, value)
select i.id, v.label, v.value
from (values
  ('wireless-headphones',  'Color',           'Black'),
  ('wireless-headphones',  'Connectivity',    'Bluetooth 5.3'),
  ('wireless-headphones',  'Battery life',    '30 hours'),
  ('mechanical-keyboard',  'Layout',          '75%'),
  ('mechanical-keyboard',  'Switch type',     'Tactile brown'),
  ('mechanical-keyboard',  'Connection',      'USB-C'),
  ('led-desk-lamp',        'Color',           'White'),
  ('led-desk-lamp',        'Material',        'Aluminum'),
  ('led-desk-lamp',        'Brightness levels', '3'),
  ('ceramic-mug',          'Color',           'White'),
  ('ceramic-mug',          'Capacity',        '350 ml'),
  ('ceramic-mug',          'Material',        'Stoneware'),
  ('dotted-notebook',      'Color',           'Navy'),
  ('dotted-notebook',      'Page count',      '192'),
  ('dotted-notebook',      'Size',            'A5'),
  ('canvas-backpack',      'Color',           'Khaki'),
  ('canvas-backpack',      'Capacity',        '22 L'),
  ('canvas-backpack',      'Material',        'Water-resistant canvas'),
  ('steel-water-bottle',   'Color',           'Brushed silver'),
  ('steel-water-bottle',   'Capacity',        '750 ml'),
  ('steel-water-bottle',   'Material',        'Stainless steel'),
  ('aluminum-phone-stand', 'Color',           'Space grey'),
  ('aluminum-phone-stand', 'Material',        'Aluminum'),
  ('aluminum-phone-stand', 'Foldable',        'Yes')
) as v(item_slug, label, value)
join public.items i on i.slug = v.item_slug
where not exists (select 1 from public.item_properties);

-- ---------------------------------------------------------------------------
-- item_medias — media attached to an item; the lowest position is the primary
-- media. Image-only here (media_type 'image'); placeholder images come from
-- picsum.photos with a stable per-image seed.
-- Runs once: skipped if item_medias already has rows.
-- ---------------------------------------------------------------------------
insert into public.item_medias (item_id, media_type, url, alt, position)
select i.id, 'image'::public.item_media_type, v.url, v.alt, v.position
from (values
  ('wireless-headphones',  'https://picsum.photos/seed/wireless-headphones-1/800/800',  'Wireless Headphones — front view',          0),
  ('wireless-headphones',  'https://picsum.photos/seed/wireless-headphones-2/800/800',  'Wireless Headphones — side profile',        1),
  ('wireless-headphones',  'https://picsum.photos/seed/wireless-headphones-3/800/800',  'Wireless Headphones — in carrying case',    2),
  ('mechanical-keyboard',  'https://picsum.photos/seed/mechanical-keyboard-1/800/800',  'Mechanical Keyboard — top-down view',       0),
  ('mechanical-keyboard',  'https://picsum.photos/seed/mechanical-keyboard-2/800/800',  'Mechanical Keyboard — side profile',        1),
  ('mechanical-keyboard',  'https://picsum.photos/seed/mechanical-keyboard-3/800/800',  'Mechanical Keyboard — keycap close-up',     2),
  ('led-desk-lamp',        'https://picsum.photos/seed/led-desk-lamp-1/800/800',        'LED Desk Lamp — lit on a desk',             0),
  ('led-desk-lamp',        'https://picsum.photos/seed/led-desk-lamp-2/800/800',        'LED Desk Lamp — folded flat',               1),
  ('ceramic-mug',          'https://picsum.photos/seed/ceramic-mug-1/800/800',          'Ceramic Coffee Mug — front view',           0),
  ('ceramic-mug',          'https://picsum.photos/seed/ceramic-mug-2/800/800',          'Ceramic Coffee Mug — filled with coffee',   1),
  ('dotted-notebook',      'https://picsum.photos/seed/dotted-notebook-1/800/800',      'Dotted Notebook — hardcover front',         0),
  ('dotted-notebook',      'https://picsum.photos/seed/dotted-notebook-2/800/800',      'Dotted Notebook — open dotted pages',       1),
  ('canvas-backpack',      'https://picsum.photos/seed/canvas-backpack-1/800/800',      'Canvas Backpack — front view',              0),
  ('canvas-backpack',      'https://picsum.photos/seed/canvas-backpack-2/800/800',      'Canvas Backpack — worn on the back',        1),
  ('canvas-backpack',      'https://picsum.photos/seed/canvas-backpack-3/800/800',      'Canvas Backpack — laptop sleeve interior',  2),
  ('steel-water-bottle',   'https://picsum.photos/seed/steel-water-bottle-1/800/800',   'Stainless Steel Water Bottle — upright',    0),
  ('steel-water-bottle',   'https://picsum.photos/seed/steel-water-bottle-2/800/800',   'Stainless Steel Water Bottle — cap detail', 1),
  ('aluminum-phone-stand', 'https://picsum.photos/seed/aluminum-phone-stand-1/800/800', 'Aluminum Phone Stand — holding a phone',    0),
  ('aluminum-phone-stand', 'https://picsum.photos/seed/aluminum-phone-stand-2/800/800', 'Aluminum Phone Stand — folded flat',        1)
) as v(item_slug, url, alt, position)
join public.items i on i.slug = v.item_slug
where not exists (select 1 from public.item_medias);
