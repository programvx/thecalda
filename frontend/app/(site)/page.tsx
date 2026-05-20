import Link from "next/link";
import { createClient } from "@/lib/supabase/server";

export default async function HomePage() {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  return (
    <div className="w-full mx-auto max-w-5xl px-6 py-20">
      <h1 className="text-4xl font-semibold tracking-tight">thecalda</h1>
      <p className="mt-4 max-w-prose text-foreground/70">
        Phase 1 — authentication. Create an account with email and password;
        a profile is created automatically and served by the Go backend.
      </p>

      <div className="mt-8">
        {user ? (
          <Link
            href="/account"
            className="inline-block rounded-md bg-foreground px-4 py-2 text-sm text-background"
          >
            Go to your account
          </Link>
        ) : (
          <Link
            href="/auth/sign-in"
            className="inline-block rounded-md bg-foreground px-4 py-2 text-sm text-background"
          >
            Get started
          </Link>
        )}
      </div>
    </div>
  );
}
