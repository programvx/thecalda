"use client";

import { ShoppingCart } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/components/CartProvider";

/**
 * AddToCartButton adds a catalog item to the cart and opens the cart panel.
 * It is disabled when the item is unavailable to buy, or while its own add is
 * in flight (cart-line edits in the panel do not disable it).
 */
export function AddToCartButton({
  itemUid,
  itemName,
  disabled = false,
}: {
  itemUid: string;
  itemName: string;
  disabled?: boolean;
}) {
  const { addItem, addPending } = useCart();

  return (
    <Button
      type="button"
      size="lg"
      onClick={() => addItem(itemUid)}
      disabled={disabled || addPending}
      aria-label={`Add ${itemName} to cart`}
      className="w-full sm:w-auto"
    >
      <ShoppingCart />
      Add to cart
    </Button>
  );
}
