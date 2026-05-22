"use client";

import { useState } from "react";
import Image from "next/image";
import { Badge } from "@/components/ui/badge";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import type { Address, Order } from "@/lib/types";

const STATUS_LABELS: Record<string, string> = {
  not_applicable: "—",
  pending: "Pending",
  paid: "Paid",
  shipped: "Shipped",
  delivered: "Delivered",
  cancelled: "Cancelled",
};

type BadgeVariant = "default" | "secondary" | "destructive";

/** statusVariant picks a badge style for an order status. */
function statusVariant(status: string): BadgeVariant {
  if (status === "cancelled") return "destructive";
  if (status === "delivered") return "default";
  return "secondary";
}

function orderLabel(order: Order): string {
  return order.orderNumber ?? order.uid.slice(0, 8);
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

/** AddressBlock renders a formatted postal address. */
function AddressBlock({ title, address }: { title: string; address: Address }) {
  return (
    <div className="mt-6">
      <h3 className="text-sm font-semibold">{title}</h3>
      <address className="mt-1 text-sm leading-relaxed text-muted-foreground not-italic">
        {address.firstName} {address.lastName}
        <br />
        {address.addressLine1}
        {address.addressLine2 && (
          <>
            <br />
            {address.addressLine2}
          </>
        )}
        <br />
        {address.postalCode} {address.city}
        <br />
        {address.country}
        <br />
        {address.email}
        {address.phone ? ` · ${address.phone}` : ""}
      </address>
    </div>
  );
}

/** OrderDetails is the contents of the order drawer. */
function OrderDetails({ order }: { order: Order }) {
  const ship = order.shippingAddress;
  const bill = order.billingAddress;

  return (
    <>
      <SheetHeader className="border-b border-border">
        <SheetTitle className="font-mono text-sm">
          {orderLabel(order)}
        </SheetTitle>
        <SheetDescription className="sr-only">Order details</SheetDescription>
      </SheetHeader>

      <div className="flex-1 overflow-y-auto p-4">
        <dl className="grid gap-2 text-sm">
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">Status</dt>
            <dd>
              <Badge variant={statusVariant(order.status)}>
                {STATUS_LABELS[order.status] ?? order.status}
              </Badge>
            </dd>
          </div>
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">Placed</dt>
            <dd>{formatDate(order.placedAt ?? order.createdAt)}</dd>
          </div>
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">Payment</dt>
            <dd>{order.paymentMethod?.name ?? "—"}</dd>
          </div>
          <div className="flex items-center justify-between">
            <dt className="text-muted-foreground">Total</dt>
            <dd className="font-semibold">€{order.totalAmount.toFixed(2)}</dd>
          </div>
        </dl>

        <h3 className="mt-6 text-sm font-semibold">Items</h3>
        <ul className="mt-2 divide-y divide-border">
          {order.items.map((line) => {
            const image = line.item?.medias[0];
            return (
              <li key={line.uid} className="flex items-start gap-3 py-3">
                <div className="relative size-12 shrink-0 overflow-hidden rounded-lg bg-muted">
                  {image && (
                    <Image
                      src={image.url}
                      alt={image.alt ?? line.itemName}
                      fill
                      sizes="48px"
                      className="object-cover"
                    />
                  )}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium">
                    {line.itemName}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {line.quantity} × €{line.unitPriceDiscounted.toFixed(2)}
                  </p>
                </div>
                <span className="text-sm font-medium">
                  €{line.lineTotal.toFixed(2)}
                </span>
              </li>
            );
          })}
        </ul>

        {ship && bill && ship.uid === bill.uid ? (
          <AddressBlock title="Shipping & billing address" address={ship} />
        ) : (
          <>
            {ship && <AddressBlock title="Shipping address" address={ship} />}
            {bill && <AddressBlock title="Billing address" address={bill} />}
          </>
        )}

        {order.notes && (
          <div className="mt-6">
            <h3 className="text-sm font-semibold">Note</h3>
            <p className="mt-1 text-sm text-muted-foreground">{order.notes}</p>
          </div>
        )}
      </div>
    </>
  );
}

/**
 * OrdersTable renders the account orders table. Clicking a row opens a
 * right-side drawer with that order's details.
 */
export function OrdersTable({ orders }: { orders: Order[] }) {
  const [selected, setSelected] = useState<Order | null>(null);

  if (orders.length === 0) {
    return (
      <p className="px-6 py-12 text-center text-sm text-muted-foreground">
        No orders yet.
      </p>
    );
  }

  return (
    <>
      <div className="overflow-x-auto">
        <table className="w-full min-w-[640px] text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs uppercase tracking-wide text-muted-foreground">
              <th className="px-6 py-3 font-medium">Order</th>
              <th className="px-3 py-3 font-medium">Status</th>
              <th className="px-3 py-3 text-right font-medium">Items</th>
              <th className="px-3 py-3 text-right font-medium">Total</th>
              <th className="px-6 py-3 text-right font-medium">Date</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {orders.map((order) => (
              <tr
                key={order.uid}
                onClick={() => setSelected(order)}
                className="cursor-pointer transition-colors hover:bg-muted"
              >
                <td className="px-6 py-3 font-mono text-xs">
                  {orderLabel(order)}
                </td>
                <td className="px-3 py-3">
                  <Badge variant={statusVariant(order.status)}>
                    {STATUS_LABELS[order.status] ?? order.status}
                  </Badge>
                </td>
                <td className="px-3 py-3 text-right tabular-nums">
                  {order.items.reduce((n, it) => n + it.quantity, 0)}
                </td>
                <td className="px-3 py-3 text-right font-medium tabular-nums">
                  €{order.totalAmount.toFixed(2)}
                </td>
                <td className="px-6 py-3 text-right text-muted-foreground">
                  {formatDate(order.placedAt ?? order.createdAt)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Sheet
        open={selected !== null}
        onOpenChange={(open) => {
          if (!open) setSelected(null);
        }}
      >
        <SheetContent side="right" className="gap-0">
          {selected && <OrderDetails order={selected} />}
        </SheetContent>
      </Sheet>
    </>
  );
}
