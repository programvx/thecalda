/** Mirrors the backend `model.User` JSON shape (GET /api/me). */
export type User = {
  uid: string;
  authUserId: string;
  email: string;
  fullName: string;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.Catalog` JSON shape. */
export type Catalog = {
  uid: string;
  slug: string;
  name: string;
  description: string | null;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.ItemMedia` JSON shape. */
export type ItemMedia = {
  uid: string;
  mediaType: string;
  url: string;
  alt: string | null;
  position: number;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.ItemProperty` JSON shape. */
export type ItemProperty = {
  uid: string;
  label: string;
  value: string;
  createdAt: string;
  updatedAt: string;
};

/** Stock availability statuses, mirroring the `item_stock_status` enum. */
export type StockStatus =
  | "in_stock"
  | "low_stock"
  | "out_of_stock"
  | "discontinued";

/** Mirrors the backend `model.ItemStock` JSON shape. */
export type ItemStock = {
  uid: string;
  quantity: number;
  status: StockStatus;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.Item` JSON shape. */
export type Item = {
  uid: string;
  slug: string;
  sku: string | null;
  name: string;
  description: string | null;
  price: number;
  discount: number;
  priceDiscounted: number;
  isActive: boolean;
  stock: ItemStock | null;
  medias: ItemMedia[];
  properties: ItemProperty[];
  createdAt: string;
  updatedAt: string;
};

/** A catalog together with its items (GET /api/catalogs/:slug). */
export type CatalogWithItems = Catalog & {
  items: Item[];
};
