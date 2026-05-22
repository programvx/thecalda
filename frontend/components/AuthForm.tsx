"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff, LoaderCircle } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { Logo } from "@/components/Logo";
import { AppleIcon, FacebookIcon, GoogleIcon } from "@/components/BrandIcons";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export type AuthMode = "sign-in" | "sign-up";

const copy: Record<
  AuthMode,
  {
    title: string;
    subtitle: string;
    submit: string;
    switchText: string;
    switchCta: string;
    switchHref: string;
  }
> = {
  "sign-in": {
    title: "Sign in",
    subtitle: "Welcome back. Sign in to continue.",
    submit: "Sign in",
    switchText: "Don't have an account?",
    switchCta: "Sign up",
    switchHref: "/auth/sign-up",
  },
  "sign-up": {
    title: "Create your account",
    subtitle: "Get started with TheCalda.",
    submit: "Create account",
    switchText: "Already have an account?",
    switchCta: "Sign in",
    switchHref: "/auth/sign-in",
  },
};

const socialProviders = [
  { name: "Google", Icon: GoogleIcon },
  { name: "Facebook", Icon: FacebookIcon },
  { name: "Apple", Icon: AppleIcon },
];

/** Email/password auth form, shared by the sign-in and sign-up pages. */
export function AuthForm({ mode }: { mode: AuthMode }) {
  const router = useRouter();
  const supabase = createClient();
  const t = copy[mode];

  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const { error } =
      mode === "sign-in"
        ? await supabase.auth.signInWithPassword({ email, password })
        : await supabase.auth.signUp({
            email,
            password,
            options: { data: { full_name: fullName } },
          });

    if (error) {
      setError(error.message);
      setLoading(false);
      return;
    }

    // On success, keep `loading` true so the button stays disabled and the
    // spinner stays visible right up until the redirect navigates away.
    router.push("/account");
    router.refresh();
  }

  return (
    <div className="flex flex-1 items-center justify-center px-6 py-12">
      <div className="w-full max-w-md">
        <Card className="gap-0 p-8">
          <div className="flex flex-col items-center text-center">
            <Link href="/" aria-label="TheCalda home">
              <Logo orientation="vertical" />
            </Link>
            <h1 className="mt-6 text-xl font-semibold tracking-tight">
              {t.title}
            </h1>
            <p className="mt-1 text-sm text-muted-foreground">{t.subtitle}</p>
          </div>

          {/* Social sign-in — UI only, not wired up yet. */}
          <div className="mt-6 flex flex-col gap-2.5">
            {socialProviders.map(({ name, Icon }) => (
              <Button key={name} type="button" variant="outline" size="lg">
                <Icon className="size-5" />
                Continue with {name}
              </Button>
            ))}
          </div>

          <div className="my-6 flex items-center gap-3">
            <span className="h-px flex-1 bg-border" />
            <span className="text-xs uppercase tracking-wide text-muted-foreground">
              or
            </span>
            <span className="h-px flex-1 bg-border" />
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {mode === "sign-up" && (
              <div className="grid gap-1.5">
                <Label htmlFor="fullName">Full name</Label>
                <Input
                  id="fullName"
                  type="text"
                  autoComplete="name"
                  placeholder="Jane Doe"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  required
                />
              </div>
            )}

            <div className="grid gap-1.5">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="grid gap-1.5">
              <Label htmlFor="password">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete={
                    mode === "sign-in" ? "current-password" : "new-password"
                  }
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  minLength={6}
                  className="pr-10"
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setShowPassword((s) => !s)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                  className="absolute top-1/2 right-1 -translate-y-1/2 text-muted-foreground"
                >
                  {showPassword ? (
                    <EyeOff aria-hidden />
                  ) : (
                    <Eye aria-hidden />
                  )}
                </Button>
              </div>
            </div>

            {error && (
              <p className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
                {error}
              </p>
            )}

            <Button
              type="submit"
              size="lg"
              disabled={loading}
              aria-label={t.submit}
              className="mt-1"
            >
              {loading ? (
                <LoaderCircle className="animate-spin" aria-hidden />
              ) : (
                t.submit
              )}
            </Button>
          </form>
        </Card>

        <p className="mt-6 text-center text-sm text-muted-foreground">
          {t.switchText}{" "}
          <Link
            href={t.switchHref}
            className="font-medium text-foreground hover:underline"
          >
            {t.switchCta}
          </Link>
        </p>
      </div>
    </div>
  );
}
