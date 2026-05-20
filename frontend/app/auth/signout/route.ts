import { NextResponse, type NextRequest } from "next/server";
import { createClient } from "@/lib/supabase/server";

/** Signs the user out and redirects to the login page. */
export async function POST(request: NextRequest) {
  const supabase = await createClient();
  await supabase.auth.signOut();

  // 303 so the browser follows the redirect with a GET.
  return NextResponse.redirect(new URL("/auth/sign-in", request.url), {
    status: 303,
  });
}
