import type { ItemStock, StockStatus } from "@/lib/types";

const STOCK_META: Record<StockStatus, { label: string; dot: string }> = {
  in_stock: { label: "In stock", dot: "bg-green-500" },
  low_stock: { label: "Low stock", dot: "bg-amber-500" },
  out_of_stock: { label: "Out of stock", dot: "bg-red-500" },
  discontinued: { label: "Discontinued", dot: "bg-foreground/30" },
};

/**
 * StockIndicator renders a small colored dot plus a label for an item's stock
 * status. Renders nothing when the item has no stock record.
 */
export function StockIndicator({ stock }: { stock: ItemStock | null }) {
  if (!stock) {
    return null;
  }
  const meta = STOCK_META[stock.status] ?? STOCK_META.out_of_stock;
  return (
    <span className="inline-flex items-center gap-1.5 text-xs text-foreground/60">
      <span className={`size-2 rounded-full ${meta.dot}`} aria-hidden />
      {meta.label}
    </span>
  );
}
