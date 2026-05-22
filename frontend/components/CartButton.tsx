"use client";

import { usePathname } from "next/navigation";
import { ShoppingCart } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/components/CartProvider";

/**
 * CartButton is the header trigger that opens the cart side panel. It is
 * hidden on the checkout page so nothing pulls the user away from completing
 * the order.
 */
export function CartButton() {
  const pathname = usePathname();
  const { itemCount, openCart } = useCart();

  if (pathname === "/checkout") {
    return null;
  }

  return (
    <Button
      type="button"
      variant="outline"
      size="icon"
      onClick={openCart}
      aria-label={`Open cart, ${itemCount} item${itemCount === 1 ? "" : "s"}`}
      className="relative"
    >
      <ShoppingCart />
      {itemCount > 0 && (
        <span className="absolute -top-1.5 -right-1.5 flex min-w-4 items-center justify-center rounded-full bg-primary px-1 text-[10px] font-semibold text-primary-foreground">
          {itemCount}
        </span>
      )}
    </Button>
  );
}
