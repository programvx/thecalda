import { redirect } from "next/navigation";
import { UserRound } from "lucide-react";
import { createClient } from "@/lib/supabase/server";
import { apiFetch } from "@/lib/api";
import type { User } from "@/lib/types";

export default async function AccountPage() {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (!user) {
    redirect("/auth/sign-in");
  }

  const { data: profile, error } = await apiFetch<User>("/api/me");

  return (
    <div className="mx-auto w-full max-w-2xl px-6 py-12">
      <h1 className="text-2xl font-semibold tracking-tight">Account</h1>
      <p className="mt-1 text-sm text-foreground/60">
        View your profile and account details.
      </p>

      {error && (
        <p className="mt-6 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/40 dark:text-red-300">
          Could not load your profile: {error.message} (code {error.code})
        </p>
      )}

      {profile && (
        <div className="mt-6 overflow-hidden rounded-2xl border border-black/10 bg-white shadow-sm dark:border-white/10 dark:bg-white/[0.03]">
          <div className="flex items-center gap-4 border-b border-black/10 px-6 py-5 dark:border-white/10">
            <span className="flex size-12 shrink-0 items-center justify-center rounded-full bg-foreground text-base font-semibold text-background">
              {getInitials(profile.fullName) || (
                <UserRound className="size-6" aria-hidden />
              )}
            </span>
            <div className="min-w-0">
              <p className="truncate font-semibold">
                {profile.fullName || "Unnamed user"}
              </p>
              <p className="truncate text-sm text-foreground/60">
                {profile.email}
              </p>
            </div>
          </div>

          <dl className="divide-y divide-black/5 dark:divide-white/5">
            <Row label="User ID" value={profile.uid} mono />
            <Row label="Auth user ID" value={profile.authUserId} mono />
            <Row label="Member since" value={formatDate(profile.createdAt)} />
          </dl>
        </div>
      )}
    </div>
  );
}

function Row({
  label,
  value,
  mono,
}: {
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex flex-col gap-0.5 px-6 py-3.5 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
      <dt className="text-sm text-foreground/55">{label}</dt>
      <dd
        className={
          mono
            ? "break-all font-mono text-xs text-foreground/80 sm:text-right"
            : "text-sm font-medium sm:text-right"
        }
      >
        {value}
      </dd>
    </div>
  );
}

function getInitials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}
