"use client";

import Image from "next/image";
import Link from "next/link";
import { Minus, Plus, ShoppingCart, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { useCart } from "@/components/CartProvider";

/** CartPanel is the slide-out cart drawer, rendered once in the site layout. */
export function CartPanel() {
  const {
    cart,
    isOpen,
    cartPending,
    openCart,
    closeCart,
    updateQuantity,
    removeItem,
  } = useCart();

  const items = cart?.items ?? [];

  return (
    <Sheet
      open={isOpen}
      onOpenChange={(open) => (open ? openCart() : closeCart())}
    >
      <SheetContent side="right" className="gap-0">
        <SheetHeader className="border-b border-border">
          <SheetTitle className="flex items-center gap-2">
            <ShoppingCart className="size-4" aria-hidden />
            Your cart
          </SheetTitle>
          <SheetDescription className="sr-only">
            Items in your shopping cart
          </SheetDescription>
        </SheetHeader>

        {items.length === 0 ? (
          <div className="flex flex-1 items-center justify-center px-4">
            <p className="text-sm text-muted-foreground">Your cart is empty.</p>
          </div>
        ) : (
          <ul className="flex-1 divide-y divide-border overflow-y-auto px-4">
            {items.map((line) => {
              const image = line.item?.medias[0];
              return (
                <li key={line.uid} className="flex gap-3 py-4">
                  <div className="relative size-16 shrink-0 overflow-hidden rounded-lg bg-muted">
                    {image && (
                      <Image
                        src={image.url}
                        alt={image.alt ?? line.itemName}
                        fill
                        sizes="64px"
                        className="object-cover"
                      />
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium">
                      {line.itemName}
                    </p>
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      €{line.unitPriceDiscounted.toFixed(2)} each
                    </p>
                    <div className="mt-2 flex items-center gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        size="icon-xs"
                        disabled={cartPending || line.quantity <= 1}
                        onClick={() =>
                          updateQuantity(line.uid, line.quantity - 1)
                        }
                        aria-label="Decrease quantity"
                      >
                        <Minus />
                      </Button>
                      <span className="w-6 text-center text-sm tabular-nums">
                        {line.quantity}
                      </span>
                      <Button
                        type="button"
                        variant="outline"
                        size="icon-xs"
                        disabled={cartPending}
                        onClick={() =>
                          updateQuantity(line.uid, line.quantity + 1)
                        }
                        aria-label="Increase quantity"
                      >
                        <Plus />
                      </Button>
                    </div>
                  </div>

                  <div className="flex flex-col items-end justify-between">
                    <span className="text-sm font-semibold">
                      €{line.lineTotal.toFixed(2)}
                    </span>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      disabled={cartPending}
                      onClick={() => removeItem(line.uid)}
                      aria-label={`Remove ${line.itemName}`}
                      className="text-muted-foreground hover:text-destructive"
                    >
                      <Trash2 />
                    </Button>
                  </div>
                </li>
              );
            })}
          </ul>
        )}

        <SheetFooter className="border-t border-border">
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground">Total</span>
            <span className="text-base font-semibold">
              €{(cart?.totalAmount ?? 0).toFixed(2)}
            </span>
          </div>
          {items.length > 0 ? (
            <Button asChild>
              <Link href="/checkout" onClick={closeCart}>
                Checkout
              </Link>
            </Button>
          ) : (
            <Button disabled>Checkout</Button>
          )}
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
