"use client";

import { createContext, useContext, useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import { addToCart, removeCartItem, updateCartItem } from "@/lib/actions/cart";
import type { Order } from "@/lib/types";

type CartContextValue = {
  cart: Order | null;
  isOpen: boolean;
  addPending: boolean;
  cartPending: boolean;
  itemCount: number;
  openCart: () => void;
  closeCart: () => void;
  addItem: (itemUid: string) => void;
  updateQuantity: (orderItemUid: string, quantity: number) => void;
  removeItem: (orderItemUid: string) => void;
  clearCart: () => void;
};

const CartContext = createContext<CartContextValue | null>(null);

/** useCart exposes the shared cart state. Must be used within a CartProvider. */
export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return ctx;
}

/**
 * CartProvider holds the cart state shared by the header button, the add-to-
 * cart button, and the cart panel. Mutations run as server actions; the cart
 * is replaced with whatever the action returns.
 *
 * Add-to-cart and cart-line mutations use separate transitions, so updating a
 * quantity in the panel does not disable the add-to-cart button on the item
 * page (and vice versa).
 */
export function CartProvider({
  initialCart,
  signedIn,
  children,
}: {
  initialCart: Order | null;
  signedIn: boolean;
  children: React.ReactNode;
}) {
  const [cart, setCart] = useState<Order | null>(initialCart);
  const [isOpen, setIsOpen] = useState(false);
  const [addPending, startAddTransition] = useTransition();
  const [cartPending, startCartTransition] = useTransition();
  const router = useRouter();

  const itemCount = (cart?.items ?? []).reduce((n, it) => n + it.quantity, 0);

  function addItem(itemUid: string) {
    // The cart belongs to a signed-in user; send guests to sign in first.
    if (!signedIn) {
      router.push("/auth/sign-in");
      return;
    }
    startAddTransition(async () => {
      setCart(await addToCart(itemUid));
      setIsOpen(true);
    });
  }

  function updateQuantity(orderItemUid: string, quantity: number) {
    if (quantity < 1) {
      return;
    }
    startCartTransition(async () => {
      setCart(await updateCartItem(orderItemUid, quantity));
    });
  }

  function removeItem(orderItemUid: string) {
    startCartTransition(async () => {
      setCart(await removeCartItem(orderItemUid));
    });
  }

  return (
    <CartContext.Provider
      value={{
        cart,
        isOpen,
        addPending,
        cartPending,
        itemCount,
        openCart: () => setIsOpen(true),
        closeCart: () => setIsOpen(false),
        addItem,
        updateQuantity,
        removeItem,
        clearCart: () => setCart(null),
      }}
    >
      {children}
    </CartContext.Provider>
  );
}
