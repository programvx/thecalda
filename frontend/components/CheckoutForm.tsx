"use client";

import { useState, useTransition } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Textarea } from "@/components/ui/textarea";
import { useCart } from "@/components/CartProvider";
import { checkoutCart } from "@/lib/actions/cart";
import type { CheckoutAddress } from "@/lib/types";

// Mirrors the seeded public.payment_methods rows.
const PAYMENT_METHODS = [
  { code: "card", label: "Credit / Debit Card" },
  { code: "paypal", label: "PayPal" },
  { code: "bank_transfer", label: "Bank Transfer" },
  { code: "cash_on_delivery", label: "Cash on Delivery" },
];

type AddressForm = {
  firstName: string;
  lastName: string;
  addressLine1: string;
  addressLine2: string;
  city: string;
  postalCode: string;
  country: string;
  phone: string;
};

const emptyAddress: AddressForm = {
  firstName: "",
  lastName: "",
  addressLine1: "",
  addressLine2: "",
  city: "",
  postalCode: "",
  country: "",
  phone: "",
};

/** toCheckoutAddress maps the form state to the checkout request shape,
 *  turning blank optional fields into null. */
function toCheckoutAddress(a: AddressForm): CheckoutAddress {
  return {
    firstName: a.firstName,
    lastName: a.lastName,
    phone: a.phone.trim() || null,
    addressLine1: a.addressLine1,
    addressLine2: a.addressLine2.trim() || null,
    city: a.city,
    postalCode: a.postalCode,
    country: a.country,
  };
}

/** Field is a labelled text input. */
function Field({
  label,
  id,
  value,
  onChange,
  type = "text",
  required = false,
  autoComplete,
}: {
  label: string;
  id: string;
  value: string;
  onChange: (value: string) => void;
  type?: string;
  required?: boolean;
  autoComplete?: string;
}) {
  return (
    <div className="grid gap-1.5">
      <Label htmlFor={id}>{label}</Label>
      <Input
        id={id}
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        autoComplete={autoComplete}
      />
    </div>
  );
}

/** AddressFields renders the inputs for one postal address. */
function AddressFields({
  idPrefix,
  values,
  onChange,
}: {
  idPrefix: string;
  values: AddressForm;
  onChange: (field: keyof AddressForm, value: string) => void;
}) {
  return (
    <div className="grid gap-3">
      <div className="grid grid-cols-2 gap-3">
        <Field
          label="First name"
          id={`${idPrefix}-first-name`}
          value={values.firstName}
          onChange={(v) => onChange("firstName", v)}
          required
          autoComplete="given-name"
        />
        <Field
          label="Last name"
          id={`${idPrefix}-last-name`}
          value={values.lastName}
          onChange={(v) => onChange("lastName", v)}
          required
          autoComplete="family-name"
        />
      </div>
      <Field
        label="Address"
        id={`${idPrefix}-address1`}
        value={values.addressLine1}
        onChange={(v) => onChange("addressLine1", v)}
        required
        autoComplete="address-line1"
      />
      <Field
        label="Apartment, suite, etc. (optional)"
        id={`${idPrefix}-address2`}
        value={values.addressLine2}
        onChange={(v) => onChange("addressLine2", v)}
        autoComplete="address-line2"
      />
      <div className="grid grid-cols-2 gap-3">
        <Field
          label="City"
          id={`${idPrefix}-city`}
          value={values.city}
          onChange={(v) => onChange("city", v)}
          required
          autoComplete="address-level2"
        />
        <Field
          label="Postal code"
          id={`${idPrefix}-postal-code`}
          value={values.postalCode}
          onChange={(v) => onChange("postalCode", v)}
          required
          autoComplete="postal-code"
        />
      </div>
      <Field
        label="Country"
        id={`${idPrefix}-country`}
        value={values.country}
        onChange={(v) => onChange("country", v)}
        required
        autoComplete="country-name"
      />
      <Field
        label="Phone (optional)"
        id={`${idPrefix}-phone`}
        type="tel"
        value={values.phone}
        onChange={(v) => onChange("phone", v)}
        autoComplete="tel"
      />
    </div>
  );
}

/** Section is one titled card of the checkout form. */
function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>{children}</CardContent>
    </Card>
  );
}

/**
 * CheckoutForm is the checkout page: a Shopify-style two-column layout with
 * the contact / address / payment / note form on the left and the cart
 * summary on the right, each block a card. Submitting places the order and
 * redirects to the success page.
 */
export function CheckoutForm() {
  const { cart, clearCart } = useCart();
  const router = useRouter();
  const items = cart?.items ?? [];

  const [email, setEmail] = useState("");
  const [shipping, setShipping] = useState<AddressForm>(emptyAddress);
  const [billingSame, setBillingSame] = useState(true);
  const [billing, setBilling] = useState<AddressForm>(emptyAddress);
  const [paymentMethod, setPaymentMethod] = useState("card");
  const [note, setNote] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [placed, setPlaced] = useState(false);
  const [submitting, startSubmit] = useTransition();

  function updateShipping(field: keyof AddressForm, value: string) {
    setShipping((s) => ({ ...s, [field]: value }));
  }
  function updateBilling(field: keyof AddressForm, value: string) {
    setBilling((b) => ({ ...b, [field]: value }));
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!cart) {
      return;
    }
    setError(null);
    startSubmit(async () => {
      const result = await checkoutCart(cart.uid, {
        email,
        note: note.trim() || null,
        paymentMethodCode: paymentMethod,
        shippingAddress: toCheckoutAddress(shipping),
        billingAddress: billingSame ? null : toCheckoutAddress(billing),
      });
      if (result.error) {
        setError(result.error);
        return;
      }
      // Hold a "redirecting" view so clearing the cart doesn't briefly flash
      // the empty-cart state before the success page loads.
      setPlaced(true);
      clearCart();
      router.push(
        `/checkout/success?order=${encodeURIComponent(result.orderNumber ?? "")}`,
      );
    });
  }

  if (placed) {
    return (
      <div className="mx-auto w-full max-w-5xl px-6 py-24 text-center text-sm text-muted-foreground">
        Order placed — redirecting…
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="mx-auto w-full max-w-5xl px-6 py-16">
        <h1 className="text-3xl font-semibold tracking-tight">Checkout</h1>
        <p className="mt-3 text-muted-foreground">
          Your cart is empty — add some items before checking out.
        </p>
        <Button asChild className="mt-6">
          <Link href="/">Browse the store</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-6 py-16">
      <h1 className="text-3xl font-semibold tracking-tight">Checkout</h1>

      <div className="mt-8 grid gap-8 lg:grid-cols-[1.6fr_1fr]">
        {/* Left — contact, addresses, payment, note */}
        <form onSubmit={handleSubmit} className="grid gap-6">
          <Section title="Contact">
            <Field
              label="Email"
              id="contact-email"
              type="email"
              value={email}
              onChange={setEmail}
              required
              autoComplete="email"
            />
          </Section>

          <Section title="Shipping address">
            <AddressFields
              idPrefix="shipping"
              values={shipping}
              onChange={updateShipping}
            />
          </Section>

          <Section title="Billing address">
            <Label
              htmlFor="billing-same"
              className="flex cursor-pointer items-center gap-2 font-normal"
            >
              <Checkbox
                id="billing-same"
                checked={billingSame}
                onCheckedChange={(checked) => setBillingSame(checked === true)}
              />
              Same as shipping address
            </Label>
            {!billingSame && (
              <div className="mt-4">
                <AddressFields
                  idPrefix="billing"
                  values={billing}
                  onChange={updateBilling}
                />
              </div>
            )}
          </Section>

          <Section title="Payment method">
            <RadioGroup
              value={paymentMethod}
              onValueChange={setPaymentMethod}
              className="grid gap-2"
            >
              {PAYMENT_METHODS.map((method) => (
                <Label
                  key={method.code}
                  htmlFor={`pay-${method.code}`}
                  className={`flex cursor-pointer items-center gap-3 rounded-lg border px-3 py-2.5 text-sm font-normal ${
                    paymentMethod === method.code
                      ? "border-primary"
                      : "border-border"
                  }`}
                >
                  <RadioGroupItem value={method.code} id={`pay-${method.code}`} />
                  {method.label}
                </Label>
              ))}
            </RadioGroup>
          </Section>

          <Section title="Order note">
            <Textarea
              id="order-note"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              rows={3}
              placeholder="Anything we should know about your order? (optional)"
            />
          </Section>

          {error && (
            <p className="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive">
              {error}
            </p>
          )}

          <Button
            type="submit"
            size="lg"
            disabled={submitting}
            className="h-12 text-base"
          >
            {submitting ? "Placing order…" : "Place order"}
          </Button>
        </form>

        {/* Right — order summary */}
        <Card className="h-fit lg:sticky lg:top-6">
          <CardHeader>
            <CardTitle className="text-base">Order summary</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="divide-y divide-border">
              {items.map((line) => {
                const image = line.item?.medias[0];
                return (
                  <li key={line.uid} className="flex items-start gap-3 py-3">
                    <div className="relative size-14 shrink-0 overflow-hidden rounded-lg bg-muted">
                      {image && (
                        <Image
                          src={image.url}
                          alt={image.alt ?? line.itemName}
                          fill
                          sizes="56px"
                          className="object-cover"
                        />
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium">
                        {line.itemName}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Qty {line.quantity}
                      </p>
                    </div>
                    <span className="text-sm font-semibold">
                      €{line.lineTotal.toFixed(2)}
                    </span>
                  </li>
                );
              })}
            </ul>

            <div className="mt-4 flex items-center justify-between border-t border-border pt-4">
              <span className="text-sm text-muted-foreground">Total</span>
              <span className="text-lg font-semibold">
                €{(cart?.totalAmount ?? 0).toFixed(2)}
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
