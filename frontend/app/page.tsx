import Link from "next/link";
import { ArrowRight, BarChart3, LogIn, ReceiptText, WalletCards } from "lucide-react";

export default function Home() {
  return (
    <main className="min-h-screen bg-slate-50">
      <section className="mx-auto grid min-h-screen max-w-6xl items-center gap-10 px-6 py-10 lg:grid-cols-[1fr_0.9fr]">
        <div>
          <div className="mb-6 flex h-12 w-12 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <h1 className="max-w-xl text-4xl font-bold leading-tight text-slate-950 sm:text-5xl">
            Expense Tracker
          </h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-slate-600 sm:text-lg">
            A personal expense tracker powered by a Go + Beego REST API.
          </p>
          <div className="mt-8 flex flex-col gap-3 sm:flex-row">
            <Link
              href="/login"
              className="inline-flex h-11 items-center justify-center gap-2 rounded-lg bg-violet-600 px-5 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
            >
              <LogIn className="h-4 w-4" aria-hidden="true" />
              Login
            </Link>
            <Link
              href="/dashboard"
              className="inline-flex h-11 items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-5 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-50"
            >
              Go to Dashboard
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
          </div>
        </div>

        <div className="rounded-lg border border-slate-200 bg-white p-5 shadow-xl shadow-slate-200/80">
          <div className="mb-5 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-slate-950">Dashboard preview</p>
              <p className="text-xs text-slate-500">Static UI foundation</p>
            </div>
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-violet-50 text-violet-700">
              <BarChart3 className="h-4 w-4" aria-hidden="true" />
            </div>
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            {[
              ["Total Expenses", "BDT 42,580"],
              ["Transactions", "128"],
              ["Daily Average", "BDT 1,419"],
              ["Top Category", "Food"],
            ].map(([label, value]) => (
              <div
                key={label}
                className="rounded-lg border border-slate-100 bg-slate-50 p-4"
              >
                <p className="text-xs font-medium text-slate-500">{label}</p>
                <p className="mt-2 text-lg font-bold text-slate-950">{value}</p>
              </div>
            ))}
          </div>

          <div className="mt-5 rounded-lg border border-slate-100">
            {["Lunch at Pizza Hut", "Metro card recharge", "Monthly internet"].map(
              (item, index) => (
                <div
                  key={item}
                  className="flex items-center justify-between border-b border-slate-100 px-4 py-3 last:border-b-0"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600">
                      <ReceiptText className="h-4 w-4" aria-hidden="true" />
                    </div>
                    <span className="text-sm font-medium text-slate-700">{item}</span>
                  </div>
                  <span className="text-sm font-semibold text-slate-950">
                    BDT {index === 0 ? "560" : index === 1 ? "1,200" : "1,050"}
                  </span>
                </div>
              ),
            )}
          </div>
        </div>
      </section>
    </main>
  );
}
