import { ShoppingBag } from "lucide-react";

/** TheCalda brand logo — icon mark plus wordmark, laid out horizontally or stacked. */
export function Logo({
  orientation = "horizontal",
}: {
  orientation?: "horizontal" | "vertical";
}) {
  return (
    <span
      className={
        orientation === "vertical"
          ? "inline-flex flex-col items-center gap-2"
          : "inline-flex items-center gap-2"
      }
    >
      <span className="flex size-8 items-center justify-center rounded-lg bg-foreground text-background">
        <ShoppingBag className="size-5" aria-hidden />
      </span>
      <span className="text-lg font-semibold tracking-tight">TheCalda</span>
    </span>
  );
}
