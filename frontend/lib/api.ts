import { createClient } from "@/lib/supabase/server";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ApiError = {
  code: number;
  message: string;
  details?: string[];
};

/** Envelope returned by every backend endpoint. */
export type ApiResp<T> = {
  data?: T;
  error?: ApiError;
};

/**
 * Calls the Go backend from the server, forwarding the caller's Supabase
 * access token as a Bearer credential. Server-only — it reads request cookies.
 */
export async function apiFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<ApiResp<T>> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  const headers = new Headers(init?.headers);
  if (session?.access_token) {
    headers.set("Authorization", `Bearer ${session.access_token}`);
  }

  try {
    const res = await fetch(`${API_URL}${path}`, {
      ...init,
      headers,
      cache: "no-store",
    });
    // 204 No Content (e.g. DELETE) carries no body to parse.
    if (res.status === 204) {
      return {};
    }
    return (await res.json()) as ApiResp<T>;
  } catch {
    return {
      error: { code: 503, message: "Backend API is unreachable" },
    };
  }
}
