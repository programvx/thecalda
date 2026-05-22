"use server";

import { apiFetch } from "@/lib/api";
import type { CheckoutPayload, Order } from "@/lib/types";

const JSON_HEADERS = { "Content-Type": "application/json" };

/** getCart returns the authenticated user's cart, or null if they have none. */
export async function getCart(): Promise<Order | null> {
  const { data } = await apiFetch<Order[]>("/api/orders?type=cart&pageSize=1");
  return data?.[0] ?? null;
}

/**
 * addToCart adds a catalog item to the user's cart, creating the cart on the
 * first add. A re-add of an item already in the cart is left as-is (the
 * backend rejects the duplicate). Returns the updated cart.
 */
export async function addToCart(
  itemUid: string,
  quantity = 1,
): Promise<Order | null> {
  const cart = await getCart();

  if (!cart) {
    await apiFetch<Order>("/api/orders", {
      method: "POST",
      headers: JSON_HEADERS,
      body: JSON.stringify({ type: "cart", items: [{ itemUid, quantity }] }),
    });
  } else {
    await apiFetch("/api/order-items", {
      method: "POST",
      headers: JSON_HEADERS,
      body: JSON.stringify({ orderUid: cart.uid, itemUid, quantity }),
    });
  }

  return getCart();
}

/** updateCartItem sets the quantity of a cart line. */
export async function updateCartItem(
  orderItemUid: string,
  quantity: number,
): Promise<Order | null> {
  await apiFetch(`/api/order-items/${orderItemUid}`, {
    method: "PUT",
    headers: JSON_HEADERS,
    body: JSON.stringify({ quantity }),
  });
  return getCart();
}

/** removeCartItem removes a line from the cart. */
export async function removeCartItem(
  orderItemUid: string,
): Promise<Order | null> {
  await apiFetch(`/api/order-items/${orderItemUid}`, { method: "DELETE" });
  return getCart();
}

/** checkoutCart places the cart as an order. Returns the new order number,
 *  or an error message the form can show. */
export async function checkoutCart(
  orderUid: string,
  payload: CheckoutPayload,
): Promise<{ orderNumber?: string; error?: string }> {
  const { data, error } = await apiFetch<Order>(
    `/api/orders/${orderUid}/checkout`,
    {
      method: "POST",
      headers: JSON_HEADERS,
      body: JSON.stringify(payload),
    },
  );
  if (error || !data) {
    return { error: error?.message ?? "Checkout failed. Please try again." };
  }
  return { orderNumber: data.orderNumber ?? data.uid };
}
