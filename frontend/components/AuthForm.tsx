"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { Logo } from "@/components/Logo";
import { AppleIcon, FacebookIcon, GoogleIcon } from "@/components/BrandIcons";

export type AuthMode = "sign-in" | "sign-up";

const inputClass =
  "w-full rounded-lg border border-black/15 bg-transparent px-3 py-2 text-sm outline-none transition focus:border-foreground/40 focus:ring-2 focus:ring-foreground/10 dark:border-white/15";

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
    subtitle: "Get started with thecalda.",
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

    setLoading(false);

    if (error) {
      setError(error.message);
      return;
    }

    router.push("/account");
    router.refresh();
  }

  return (
    <div className="flex flex-1 items-center justify-center px-6 py-12">
      <div className="w-full max-w-md">
        <div className="rounded-2xl border border-black/10 bg-white p-8 shadow-sm dark:border-white/10 dark:bg-white/[0.03]">
          <div className="flex flex-col items-center text-center">
            <Link href="/" aria-label="thecalda home">
              <Logo orientation="vertical" />
            </Link>
            <h1 className="mt-6 text-xl font-semibold tracking-tight">
              {t.title}
            </h1>
            <p className="mt-1 text-sm text-foreground/60">{t.subtitle}</p>
          </div>

          {/* Social sign-in — UI only, not wired up yet. */}
          <div className="mt-6 flex flex-col gap-2.5">
            {socialProviders.map(({ name, Icon }) => (
              <button
                key={name}
                type="button"
                className="flex w-full items-center justify-center gap-2.5 rounded-lg border border-black/15 px-4 py-2.5 text-sm font-medium transition hover:bg-black/[0.03] dark:border-white/15 dark:hover:bg-white/5"
              >
                <Icon className="size-5" />
                Continue with {name}
              </button>
            ))}
          </div>

          <div className="my-6 flex items-center gap-3">
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
            <span className="text-xs uppercase tracking-wide text-foreground/40">
              or
            </span>
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {mode === "sign-up" && (
              <div className="flex flex-col gap-1.5">
                <label htmlFor="fullName" className="text-sm font-medium">
                  Full name
                </label>
                <input
                  id="fullName"
                  type="text"
                  autoComplete="name"
                  placeholder="Jane Doe"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  required
                  className={inputClass}
                />
              </div>
            )}

            <div className="flex flex-col gap-1.5">
              <label htmlFor="email" className="text-sm font-medium">
                Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className={inputClass}
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="password" className="text-sm font-medium">
                Password
              </label>
              <div className="relative">
                <input
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
                  className={`${inputClass} pr-10`}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((s) => !s)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                  className="absolute right-1 top-1/2 -translate-y-1/2 rounded-md p-1.5 text-foreground/50 transition hover:text-foreground"
                >
                  {showPassword ? (
                    <EyeOff className="size-4" aria-hidden />
                  ) : (
                    <Eye className="size-4" aria-hidden />
                  )}
                </button>
              </div>
            </div>

            {error && (
              <p className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300">
                {error}
              </p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="mt-1 rounded-lg bg-foreground px-4 py-2.5 text-sm font-medium text-background transition hover:opacity-90 disabled:opacity-50"
            >
              {loading ? "Please wait…" : t.submit}
            </button>
          </form>
        </div>

        <p className="mt-6 text-center text-sm text-foreground/60">
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
