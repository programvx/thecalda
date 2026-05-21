"use client";

import { useState } from "react";
import { Check, ShoppingCart } from "lucide-react";

/**
 * AddToCartButton is a placeholder cart action. The cart itself is not built
 * yet (Phase 2 covers the catalog only); the button gives local click feedback
 * so it reads as interactive. It is disabled when the item is unavailable.
 */
export function AddToCartButton({
  itemName,
  disabled = false,
}: {
  itemName: string;
  disabled?: boolean;
}) {
  const [added, setAdded] = useState(false);

  function handleClick() {
    setAdded(true);
    window.setTimeout(() => setAdded(false), 2000);
  }

  return (
    <button
      type="button"
      onClick={handleClick}
      disabled={disabled}
      aria-label={`Add ${itemName} to cart`}
      className="flex w-full items-center justify-center gap-2 rounded-lg bg-foreground px-5 py-3 text-sm font-medium text-background transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50 sm:w-auto"
    >
      {added ? (
        <>
          <Check className="size-4" aria-hidden />
          Added to cart
        </>
      ) : (
        <>
          <ShoppingCart className="size-4" aria-hidden />
          Add to cart
        </>
      )}
    </button>
  );
}
