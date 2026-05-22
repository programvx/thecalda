import { createClient } from "@/lib/supabase/server";
import { getCart } from "@/lib/actions/cart";
import { Header } from "@/components/Header";
import { CartProvider } from "@/components/CartProvider";
import { CartPanel } from "@/components/CartPanel";

/** Layout for the main site — pages that carry the header and cart. */
export default async function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  // The cart belongs to a signed-in user; guests start with none.
  const initialCart = user ? await getCart() : null;

  return (
    <CartProvider initialCart={initialCart} signedIn={!!user}>
      <Header />
      <main className="flex flex-1 flex-col">{children}</main>
      <CartPanel />
    </CartProvider>
  );
}
