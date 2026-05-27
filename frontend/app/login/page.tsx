"use client";

import { type FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Loader2, LogIn, WalletCards } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useAuth } from "@/hooks/useAuth";
import { loginUser } from "@/lib/auth";
import { saveUser } from "@/lib/auth-storage";

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

export default function LoginPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!isAuthLoading && isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [isAuthenticated, isAuthLoading, router]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");

    if (!email.trim()) {
      setError("Email is required");
      return;
    }

    if (!password) {
      setError("Password is required");
      return;
    }

    setIsLoading(true);

    try {
      const response = await loginUser({
        email: email.trim(),
        password,
      });

      if (!response.success || !response.data) {
        setError(response.message || "Invalid email or password");
        return;
      }

      saveUser(response.data);
      router.replace("/dashboard");
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
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-linear-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
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
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-linear-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <h1 className="text-2xl font-bold text-slate-950">Welcome back</h1>
          <p className="mt-2 text-sm text-slate-500">
            Sign in to manage your expenses
          </p>
        </div>

        <form className="space-y-4" onSubmit={handleSubmit}>
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
              autoComplete="current-password"
              placeholder="Enter password"
              className="mt-2"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </div>

          {error ? (
            <div className="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-medium text-rose-700">
              {error}
            </div>
          ) : null}

          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
            ) : (
              <LogIn className="h-4 w-4" aria-hidden="true" />
            )}
            {isLoading ? "Logging in..." : "Login"}
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          New here?{" "}
          <Link href="/register" className="font-semibold text-violet-700 hover:text-violet-800">
            Create an account
          </Link>
        </p>
      </Card>
    </main>
    )
  );
}
