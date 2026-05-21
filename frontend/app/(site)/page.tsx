import Link from "next/link";
import { createClient } from "@/lib/supabase/server";
import { apiFetch } from "@/lib/api";
import type { Catalog } from "@/lib/types";

export default async function HomePage() {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  const { data: catalogs } = await apiFetch<Catalog[]>("/api/catalogs");
  const categories = (catalogs ?? []).filter((c) => c.isActive);

  return (
    <div className="mx-auto w-full max-w-5xl px-6 py-16">
      <section>
        <h1 className="text-4xl font-semibold tracking-tight">
          Welcome to thecalda
        </h1>
        <p className="mt-3 max-w-prose text-foreground/70">
          A small demo storefront. Browse the collections below.
        </p>
        <div className="mt-6">
          <Link
            href={user ? "/account" : "/auth/sign-in"}
            className="inline-block rounded-md bg-foreground px-4 py-2 text-sm text-background"
          >
            {user ? "Go to your account" : "Get started"}
          </Link>
        </div>
      </section>

      <section className="mt-14">
        <h2 className="text-xl font-semibold tracking-tight">
          Browse categories
        </h2>

        {categories.length === 0 ? (
          <p className="mt-3 text-sm text-foreground/60">
            No categories available yet.
          </p>
        ) : (
          <ul className="mt-5 grid gap-5 sm:grid-cols-2">
            {categories.map((category) => (
              <li key={category.uid}>
                <Link
                  href={`/catalog/${category.slug}`}
                  className="block min-h-44 rounded-xl bg-black/5 p-8 transition hover:bg-black/10 dark:bg-white/5 dark:hover:bg-white/10"
                >
                  <h3 className="text-lg font-medium">{category.name}</h3>
                  {category.description && (
                    <p className="mt-2 text-sm text-foreground/60">
                      {category.description}
                    </p>
                  )}
                </Link>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
