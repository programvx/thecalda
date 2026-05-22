import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import { apiFetch } from "@/lib/api";
import { ItemGallery } from "@/components/ItemGallery";
import { AddToCartButton } from "@/components/AddToCartButton";
import { StockIndicator } from "@/components/StockIndicator";
import type { CatalogWithItems } from "@/lib/types";

export default async function ItemPage({
  params,
}: {
  params: Promise<{ slug: string; itemSlug: string }>;
}) {
  const { slug, itemSlug } = await params;
  const { data: catalog } = await apiFetch<CatalogWithItems>(
    `/api/catalogs/${slug}`,
  );

  const item = catalog?.items.find((i) => i.slug === itemSlug);
  if (!catalog || !item) {
    notFound();
  }

  // An item can't be added to the cart when it isn't available to buy.
  const unavailable =
    item.stock?.status === "out_of_stock" ||
    item.stock?.status === "discontinued";

  return (
    <div className="mx-auto w-full max-w-5xl px-6 py-16">
      <Link
        href={`/catalog/${catalog.slug}`}
        className="inline-flex items-center gap-1.5 text-sm text-foreground/60 transition hover:text-foreground"
      >
        <ArrowLeft className="size-4" aria-hidden />
        Back to {catalog.name}
      </Link>

      <div className="mt-6 grid gap-10 md:grid-cols-[3fr_2fr]">
        <ItemGallery
          medias={item.medias}
          itemName={item.name}
          discount={item.discount}
        />

        <div>
          <h1 className="text-3xl font-semibold tracking-tight">{item.name}</h1>
          {item.sku && (
            <p className="mt-1 text-xs uppercase tracking-wide text-foreground/40">
              SKU {item.sku}
            </p>
          )}

          <div className="mt-5 flex items-baseline gap-3">
            {item.discount > 0 ? (
              <>
                <span className="text-2xl font-semibold">
                  €{item.priceDiscounted.toFixed(2)}
                </span>
                <span className="text-foreground/50 line-through">
                  €{item.price.toFixed(2)}
                </span>
              </>
            ) : (
              <span className="text-2xl font-semibold">
                €{item.price.toFixed(2)}
              </span>
            )}
          </div>

          <div className="mt-3">
            <StockIndicator stock={item.stock} />
          </div>

          <div className="mt-6">
            <AddToCartButton
              itemUid={item.uid}
              itemName={item.name}
              disabled={unavailable}
            />
          </div>

          {item.description && (
            <p className="mt-8 max-w-prose text-foreground/70">
              {item.description}
            </p>
          )}

          {item.properties.length > 0 && (
            <div className="mt-8">
              <h2 className="text-sm font-semibold tracking-tight">Details</h2>
              <dl className="mt-3 border-t border-black/10 text-sm dark:border-white/10">
                {item.properties.map((prop) => (
                  <div
                    key={prop.uid}
                    className="flex justify-between gap-4 border-b border-black/10 py-2 dark:border-white/10"
                  >
                    <dt className="text-foreground/60">{prop.label}</dt>
                    <dd className="font-medium">{prop.value}</dd>
                  </div>
                ))}
              </dl>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
