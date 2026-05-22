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

/** Mirrors the backend `model.OrderItem` JSON shape — one line of an order.
 *  `item` is the linked catalog item; only its media is relied on here. */
export type OrderItem = {
  uid: string;
  itemName: string;
  quantity: number;
  unitPrice: number;
  unitDiscount: number;
  unitPriceDiscounted: number;
  lineTotal: number;
  item: { medias: ItemMedia[] } | null;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.Address` JSON shape. */
export type Address = {
  uid: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string | null;
  addressLine1: string;
  addressLine2: string | null;
  postalCode: string;
  city: string;
  country: string;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.PaymentMethod` JSON shape. */
export type PaymentMethod = {
  uid: string;
  code: string;
  name: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

/** Mirrors the backend `model.Order` JSON shape. A cart is an order with
 *  type "cart". The address / payment fields are null for carts. */
export type Order = {
  uid: string;
  type: "cart" | "order";
  status: string;
  totalAmount: number;
  orderNumber: string | null;
  notes: string | null;
  placedAt: string | null;
  billingAddress: Address | null;
  shippingAddress: Address | null;
  paymentMethod: PaymentMethod | null;
  items: OrderItem[];
  createdAt: string;
  updatedAt: string;
};

/** One postal address in a checkout request. */
export type CheckoutAddress = {
  firstName: string;
  lastName: string;
  phone: string | null;
  addressLine1: string;
  addressLine2: string | null;
  city: string;
  postalCode: string;
  country: string;
};

/** Request body for checking out a cart (POST /api/orders/:uid/checkout). */
export type CheckoutPayload = {
  email: string;
  note: string | null;
  paymentMethodCode: string;
  shippingAddress: CheckoutAddress;
  billingAddress: CheckoutAddress | null;
};
