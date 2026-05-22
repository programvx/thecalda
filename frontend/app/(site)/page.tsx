import Link from "next/link";
import { createClient } from "@/lib/supabase/server";
import { apiFetch } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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
          Welcome to TheCalda
        </h1>
        <p className="mt-3 max-w-prose text-foreground/70">
          A small demo storefront. Browse the collections below.
        </p>
        {!user && (
          <div className="mt-6">
            <Button asChild>
              <Link href="/auth/sign-in">Get started</Link>
            </Button>
          </div>
        )}
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
                  className="block h-full"
                >
                  <Card className="h-full min-h-44 justify-center transition-colors hover:bg-muted">
                    <CardHeader>
                      <CardTitle className="text-lg">{category.name}</CardTitle>
                      {category.description && (
                        <CardDescription>
                          {category.description}
                        </CardDescription>
                      )}
                    </CardHeader>
                  </Card>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
