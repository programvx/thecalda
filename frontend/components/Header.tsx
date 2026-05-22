import Link from "next/link";
import { LogIn } from "lucide-react";
import { createClient } from "@/lib/supabase/server";
import { Logo } from "@/components/Logo";
import { UserMenu } from "@/components/UserMenu";
import { CartButton } from "@/components/CartButton";
import { Button } from "@/components/ui/button";

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
    <header className="px-4 py-4">
      <nav className="mx-auto flex max-w-5xl items-center justify-between gap-4 rounded-full border bg-card px-5 py-2.5 shadow-sm">
        <Link href="/" aria-label="TheCalda home">
          <Logo />
        </Link>

        <div className="flex items-center gap-3 text-sm">
          {user ? (
            <>
              <CartButton />
              <UserMenu name={displayName} />
            </>
          ) : (
            <Button asChild variant="outline">
              <Link href="/auth/sign-in">
                <LogIn aria-hidden />
                Sign in
              </Link>
            </Button>
          )}
        </div>
      </nav>
    </header>
  );
}
