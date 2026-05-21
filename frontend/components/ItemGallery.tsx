"use client";

import { useState } from "react";
import Image from "next/image";
import type { ItemMedia } from "@/lib/types";

/**
 * ItemGallery shows an item's primary media large, with a thumbnail strip to
 * switch the active image. Falls back to a placeholder when the item has no
 * media. A discount badge is overlaid on the top-right of the main image.
 */
export function ItemGallery({
  medias,
  itemName,
  discount,
}: {
  medias: ItemMedia[];
  itemName: string;
  discount: number;
}) {
  const [selected, setSelected] = useState(0);

  if (medias.length === 0) {
    return (
      <div className="relative aspect-square rounded-xl bg-black/5 dark:bg-white/10">
        <DiscountBadge discount={discount} />
      </div>
    );
  }

  const active = medias[selected] ?? medias[0];

  return (
    <div>
      <div className="relative aspect-square overflow-hidden rounded-xl bg-black/5 dark:bg-white/10">
        <Image
          src={active.url}
          alt={active.alt ?? itemName}
          fill
          sizes="(min-width: 768px) 55vw, 90vw"
          className="object-cover"
          priority
        />
        <DiscountBadge discount={discount} />
      </div>

      {medias.length > 1 && (
        <div className="mt-3 grid grid-cols-4 gap-3">
          {medias.map((media, index) => (
            <button
              key={media.uid}
              type="button"
              onClick={() => setSelected(index)}
              aria-label={`View image ${index + 1}`}
              aria-current={index === selected}
              className={`relative aspect-square overflow-hidden rounded-lg bg-black/5 transition dark:bg-white/10 ${
                index === selected
                  ? "ring-2 ring-foreground"
                  : "opacity-70 hover:opacity-100"
              }`}
            >
              <Image
                src={media.url}
                alt={media.alt ?? itemName}
                fill
                sizes="120px"
                className="object-cover"
              />
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

/** DiscountBadge overlays a "-NN%" pill; renders nothing when there's no discount. */
function DiscountBadge({ discount }: { discount: number }) {
  if (discount <= 0) {
    return null;
  }
  return (
    <span className="absolute right-3 top-3 rounded-full bg-red-600 px-3 py-1 text-sm font-semibold text-white">
      -{Math.round(discount * 100)}%
    </span>
  );
}
