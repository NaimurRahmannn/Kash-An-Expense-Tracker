"use client";

import { type ReactNode, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Loader2, WalletCards } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { useAuth } from "@/hooks/useAuth";

type ProtectedRouteProps = {
  children: ReactNode;
};

function AuthCheckingCard() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
      <Card className="w-full max-w-sm p-6 text-center">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
          <WalletCards className="h-6 w-6" aria-hidden="true" />
        </div>
        <div className="mx-auto mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-violet-50 text-violet-600">
          <Loader2 className="h-5 w-5 animate-spin" aria-hidden="true" />
        </div>
        <p className="text-sm font-semibold text-slate-700">
          Checking authentication...
        </p>
      </Card>
    </main>
  );
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading || !isAuthenticated) {
    return <AuthCheckingCard />;
  }

  return children;
}
