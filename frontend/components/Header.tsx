import Link from "next/link";
import { LogIn } from "lucide-react";
import { createClient } from "@/lib/supabase/server";
import { Logo } from "@/components/Logo";
import { UserMenu } from "@/components/UserMenu";

/** Site header — shows auth state via a sign-in link or a user menu. */
export async function Header() {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  // Prefer the full name from sign-up metadata; fall back to the email.
  const displayName = user
    ? (user.user_metadata?.full_name as string | undefined) ||
      user.email ||
      ""
    : "";

  return (
    <header>
      <nav className="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
        <Link href="/" aria-label="thecalda home">
          <Logo />
        </Link>

        <div className="flex items-center gap-4 text-sm">
          {user ? (
            <UserMenu name={displayName} />
          ) : (
            <Link
              href="/auth/sign-in"
              className="flex items-center gap-1.5 rounded-md border border-black/15 px-3 py-1 hover:bg-black/5 dark:border-white/20 dark:hover:bg-white/10"
            >
              <LogIn className="size-4" aria-hidden />
              Sign in
            </Link>
          )}
        </div>
      </nav>
    </header>
  );
}
