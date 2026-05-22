import { redirect } from "next/navigation";
import { createClient } from "@/lib/supabase/server";
import { apiFetch } from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import { OrdersTable } from "@/components/OrdersTable";
import type { Order } from "@/lib/types";

export default async function AdminPage() {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (!user) {
    redirect("/auth/sign-in");
  }

  const { data: orders } = await apiFetch<Order[]>(
    "/api/orders?type=order&pageSize=100",
  );

  return (
    <div className="mx-auto w-full max-w-5xl px-6 py-12">
      <h1 className="text-2xl font-semibold tracking-tight">Orders</h1>
      <p className="mt-1 text-sm text-muted-foreground">
        Review the placed orders on your account.
      </p>

      <Card className="mt-6 gap-0 overflow-hidden py-0">
        <CardContent className="p-0">
          <OrdersTable orders={orders ?? []} />
        </CardContent>
      </Card>
    </div>
  );
}
