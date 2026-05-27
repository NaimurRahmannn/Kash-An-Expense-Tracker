import Link from "next/link";
import { LogIn, WalletCards } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
      <Card className="w-full max-w-md p-6">
        <div className="mb-6 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <h1 className="text-2xl font-bold text-slate-950">Welcome back</h1>
          <p className="mt-2 text-sm text-slate-500">
            Login UI placeholder for Expense Tracker.
          </p>
        </div>

        <form className="space-y-4">
          <div>
            <label htmlFor="email" className="text-sm font-semibold text-slate-700">
              Email
            </label>
            <Input id="email" type="email" placeholder="john.doe@example.com" className="mt-2" />
          </div>
          <div>
            <label htmlFor="password" className="text-sm font-semibold text-slate-700">
              Password
            </label>
            <Input id="password" type="password" placeholder="Enter password" className="mt-2" />
          </div>
          <Button type="button" className="w-full">
            <LogIn className="h-4 w-4" aria-hidden="true" />
            Login
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
  );
}
