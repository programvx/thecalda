import Link from "next/link";
import { CircleCheck } from "lucide-react";
import { Button } from "@/components/ui/button";

export default async function CheckoutSuccessPage({
  searchParams,
}: {
  searchParams: Promise<{ order?: string }>;
}) {
  const { order } = await searchParams;

  return (
    <div className="mx-auto flex w-full max-w-xl flex-col items-center px-6 py-24 text-center">
      <span className="flex size-14 items-center justify-center rounded-full bg-primary text-primary-foreground">
        <CircleCheck className="size-7" aria-hidden />
      </span>
      <h1 className="mt-6 text-2xl font-semibold tracking-tight">
        Order placed
      </h1>
      <p className="mt-2 text-muted-foreground">
        Thank you for your order. A confirmation will follow shortly.
      </p>
      {order && (
        <p className="mt-1 text-sm text-muted-foreground">
          Order number:{" "}
          <span className="font-medium text-foreground">{order}</span>
        </p>
      )}
      <div className="mt-8 flex flex-wrap justify-center gap-3">
        <Button asChild>
          <Link href="/">Continue shopping</Link>
        </Button>
        <Button asChild variant="outline">
          <Link href="/admin">View your orders</Link>
        </Button>
      </div>
    </div>
  );
}
