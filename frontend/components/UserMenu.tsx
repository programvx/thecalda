"use client";

import Link from "next/link";
import { ChevronDown, LogOut, Package, UserRound } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

/** Authenticated user menu — a trigger button with a dropdown holding the
 *  account link, the signed-in user, and sign-out. */
export function UserMenu({ name }: { name: string }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="gap-2 pl-1">
          <span className="flex size-6 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <UserRound className="size-3.5" aria-hidden />
          </span>
          <span className="hidden max-w-[11rem] truncate sm:inline">
            {name}
          </span>
          <ChevronDown
            className="text-muted-foreground transition-transform group-data-[state=open]/button:rotate-180"
            aria-hidden
          />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel className="font-normal">
          <span className="block text-xs text-muted-foreground">
            Signed in as
          </span>
          <span className="block truncate text-sm font-medium">{name}</span>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Link href="/account">
            <UserRound />
            Account
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <Link href="/admin">
            <Package />
            Orders
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <form action="/auth/signout" method="post">
          <DropdownMenuItem asChild>
            <button type="submit" className="w-full">
              <LogOut />
              Sign out
            </button>
          </DropdownMenuItem>
        </form>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
