import type { Metadata } from "next";
import { AuthForm } from "@/components/AuthForm";

export const metadata: Metadata = {
  title: "Sign up · TheCalda",
};

export default function SignUpPage() {
  return <AuthForm mode="sign-up" />;
}
