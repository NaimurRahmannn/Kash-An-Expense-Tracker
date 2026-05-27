"use client";

import { type FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Loader2, UserPlus, WalletCards } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useAuth } from "@/hooks/useAuth";
import { registerUser } from "@/lib/auth";

const SERVER_ERROR_MESSAGE = "Unable to connect to server. Please try again.";

function getAuthErrorMessage(error: unknown) {
  if (error instanceof Error) {
    const message = error.message.trim();

    if (
      message &&
      message !== "Failed to fetch" &&
      message !== "NetworkError when attempting to fetch resource." &&
      message !== "Request failed"
    ) {
      return message;
    }
  }

  return SERVER_ERROR_MESSAGE;
}

export default function RegisterPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!isAuthLoading && isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [isAuthenticated, isAuthLoading, router]);

  function validateForm() {
    if (!name.trim()) {
      return "Name is required";
    }

    if (!email.trim()) {
      return "Email is required";
    }

    if (!password) {
      return "Password is required";
    }

    if (password.length < 6) {
      return "Password must be at least 6 characters";
    }

    if (confirmPassword !== password) {
      return "Passwords do not match";
    }

    return "";
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSuccessMessage("");

    const validationError = validateForm();

    if (validationError) {
      setError(validationError);
      return;
    }

    setIsLoading(true);

    try {
      const response = await registerUser({
        name: name.trim(),
        email: email.trim(),
        password,
      });

      if (!response.success) {
        setError(response.message || "Unable to create account");
        return;
      }

      setSuccessMessage(response.message || "User registered successfully");
      window.setTimeout(() => {
        router.push("/login");
      }, 1200);
    } catch (caughtError) {
      setError(getAuthErrorMessage(caughtError));
    } finally {
      setIsLoading(false);
    }
  }

  return (
    isAuthLoading || isAuthenticated ? (
      <main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
        <Card className="w-full max-w-sm p-6 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <Loader2 className="mx-auto h-5 w-5 animate-spin text-violet-600" aria-hidden="true" />
          <p className="mt-3 text-sm font-semibold text-slate-700">
            Checking authentication...
          </p>
        </Card>
      </main>
    ) : (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
      <Card className="w-full max-w-md p-6">
        <div className="mb-6 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <h1 className="text-2xl font-bold text-slate-950">Create your account</h1>
          <p className="mt-2 text-sm text-slate-500">
            Start tracking your expenses today
          </p>
        </div>

        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="name" className="text-sm font-semibold text-slate-700">
              Name
            </label>
            <Input
              id="name"
              autoComplete="name"
              placeholder="John Doe"
              className="mt-2"
              value={name}
              onChange={(event) => setName(event.target.value)}
            />
          </div>
          <div>
            <label htmlFor="email" className="text-sm font-semibold text-slate-700">
              Email
            </label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="john@example.com"
              className="mt-2"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
            />
          </div>
          <div>
            <label htmlFor="password" className="text-sm font-semibold text-slate-700">
              Password
            </label>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              placeholder="Create password"
              className="mt-2"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </div>
          <div>
            <label
              htmlFor="confirm-password"
              className="text-sm font-semibold text-slate-700"
            >
              Confirm password
            </label>
            <Input
              id="confirm-password"
              type="password"
              autoComplete="new-password"
              placeholder="Confirm password"
              className="mt-2"
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
            />
          </div>

          {error ? (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-medium text-rose-700">
              {error}
            </div>
          ) : null}

          {successMessage ? (
            <div className="rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm font-medium text-emerald-700">
              {successMessage}
            </div>
          ) : null}

          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
            ) : (
              <UserPlus className="h-4 w-4" aria-hidden="true" />
            )}
            {isLoading ? "Creating account..." : "Create Account"}
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          Already have an account?{" "}
          <Link href="/login" className="font-semibold text-violet-700 hover:text-violet-800">
            Login
          </Link>
        </p>
      </Card>
    </main>
    )
  );
}
