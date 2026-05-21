import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { StockIndicator } from "@/components/StockIndicator";
import type { CatalogWithItems } from "@/lib/types";

export default async function CatalogPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const { data: catalog } = await apiFetch<CatalogWithItems>(
    `/api/catalogs/${slug}`,
  );

  if (!catalog) {
    notFound();
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-6 py-16">
      <h1 className="text-3xl font-semibold tracking-tight">{catalog.name}</h1>
      {catalog.description && (
        <p className="mt-2 max-w-prose text-foreground/70">
          {catalog.description}
        </p>
      )}

      <h2 className="mt-12 text-lg font-semibold tracking-tight">Items</h2>
      {catalog.items.length === 0 ? (
        <p className="mt-3 text-sm text-foreground/60">
          No items in this collection yet.
        </p>
      ) : (
        <ul className="mt-5 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {catalog.items.map((item) => {
            const primary = item.medias[0];
            return (
              <li key={item.uid}>
                <Link
                  href={`/catalog/${catalog.slug}/${item.slug}`}
                  className="block overflow-hidden rounded-xl bg-black/5 transition hover:bg-black/10 dark:bg-white/5 dark:hover:bg-white/10"
                >
                  <div className="relative aspect-square bg-black/5 dark:bg-white/10">
                    {primary && (
                      <Image
                        src={primary.url}
                        alt={primary.alt ?? item.name}
                        fill
                        sizes="(min-width: 1024px) 30vw, (min-width: 640px) 45vw, 90vw"
                        className="object-cover"
                      />
                    )}
                    {item.discount > 0 && (
                      <span className="absolute left-3 top-3 rounded-full bg-red-600 px-2.5 py-1 text-xs font-semibold text-white">
                        -{Math.round(item.discount * 100)}%
                      </span>
                    )}
                  </div>
                  <div className="p-6">
                    <h3 className="font-medium">{item.name}</h3>
                    {item.description && (
                      <p className="mt-1 line-clamp-2 text-sm text-foreground/60">
                        {item.description}
                      </p>
                    )}
                    <p className="mt-3 text-sm">
                      {item.discount > 0 ? (
                        <>
                          <span className="font-semibold">
                            €{item.priceDiscounted.toFixed(2)}
                          </span>{" "}
                          <span className="text-foreground/50 line-through">
                            €{item.price.toFixed(2)}
                          </span>
                        </>
                      ) : (
                        <span className="font-semibold">
                          €{item.price.toFixed(2)}
                        </span>
                      )}
                    </p>
                    <div className="mt-2">
                      <StockIndicator stock={item.stock} />
                    </div>
                  </div>
                </Link>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
