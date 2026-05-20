"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { ChevronDown, LogOut, UserRound } from "lucide-react";

/** Authenticated user menu — a trigger button with a dropdown holding the
 *  account link, the signed-in user, and sign-out. */
export function UserMenu({ name }: { name: string }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;

    function handlePointer(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    function handleKey(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }

    document.addEventListener("mousedown", handlePointer);
    document.addEventListener("keydown", handleKey);
    return () => {
      document.removeEventListener("mousedown", handlePointer);
      document.removeEventListener("keydown", handleKey);
    };
  }, [open]);

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        aria-haspopup="menu"
        aria-expanded={open}
        className="flex items-center gap-2 rounded-md border border-black/15 py-1 pl-1 pr-2 hover:bg-black/5 dark:border-white/20 dark:hover:bg-white/10"
      >
        <span className="flex size-6 items-center justify-center rounded-full bg-foreground text-background">
          <UserRound className="size-3.5" aria-hidden />
        </span>
        <span className="hidden max-w-[11rem] truncate sm:inline">{name}</span>
        <ChevronDown
          className={`size-4 text-foreground/50 transition-transform ${
            open ? "rotate-180" : ""
          }`}
          aria-hidden
        />
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-56 overflow-hidden rounded-lg border border-black/10 bg-white shadow-lg dark:border-white/10 dark:bg-neutral-900">
          <div className="border-b border-black/10 px-3 py-2.5 dark:border-white/10">
            <p className="text-xs text-foreground/50">Signed in as</p>
            <p className="truncate text-sm font-medium">{name}</p>
          </div>

          <Link
            href="/account"
            onClick={() => setOpen(false)}
            className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-black/5 dark:hover:bg-white/10"
          >
            <UserRound className="size-4 text-foreground/60" aria-hidden />
            Account
          </Link>

          <form
            action="/auth/signout"
            method="post"
            className="border-t border-black/10 dark:border-white/10"
          >
            <button
              type="submit"
              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm hover:bg-black/5 dark:hover:bg-white/10"
            >
              <LogOut className="size-4 text-foreground/60" aria-hidden />
              Sign out
            </button>
          </form>
        </div>
      )}
    </div>
  );
}
