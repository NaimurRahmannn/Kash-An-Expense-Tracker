import Link from "next/link";
import {
  ArrowRight,
  BarChart3,
  CalendarDays,
  Mic,
  PlusCircle,
  ReceiptText,
  Tag,
  WalletCards,
} from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";

const summaryCards = [
  {
    title: "Total Expenses",
    value: "BDT 42,580",
    detail: "Static monthly placeholder",
    icon: WalletCards,
    color: "from-violet-600 to-indigo-500",
  },
  {
    title: "Total Transactions",
    value: "128",
    detail: "Across all categories",
    icon: ReceiptText,
    color: "from-cyan-500 to-blue-500",
  },
  {
    title: "Daily Average",
    value: "BDT 1,419",
    detail: "Based on sample data",
    icon: CalendarDays,
    color: "from-emerald-500 to-teal-500",
  },
  {
    title: "Highest Category",
    value: "Food",
    detail: "BDT 12,460 sample total",
    icon: Tag,
    color: "from-amber-500 to-orange-500",
  },
];

const recentExpenses = [
  ["Lunch at Pizza Hut", "Food", "BDT 560"],
  ["Metro card recharge", "Transport", "BDT 1,200"],
  ["Monthly internet bill", "Utilities", "BDT 1,050"],
  ["Notebook and pens", "Education", "BDT 420"],
];

export default function DashboardPage() {
  return (
    <AppShell>
      <PageHeader
        title="Dashboard"
        description="Static dashboard foundation. Summary API data will be connected in a later frontend part."
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {summaryCards.map((item) => {
          const Icon = item.icon;

          return (
            <Card
              key={item.title}
              className="overflow-hidden bg-gradient-to-br from-white to-slate-50 p-5"
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-medium text-slate-500">{item.title}</p>
                  <p className="mt-3 text-2xl font-bold text-slate-950">{item.value}</p>
                </div>
                <div
                  className={`flex h-11 w-11 items-center justify-center rounded-lg bg-gradient-to-br ${item.color} text-white shadow-sm`}
                >
                  <Icon className="h-5 w-5" aria-hidden="true" />
                </div>
              </div>
              <p className="mt-4 text-xs font-medium text-slate-500">{item.detail}</p>
            </Card>
          );
        })}
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1.3fr_0.7fr]">
        <Card className="p-5">
          <div className="flex items-center justify-between gap-4">
            <div>
              <h2 className="text-base font-bold text-slate-950">Expense Overview</h2>
              <p className="mt-1 text-sm text-slate-500">Chart placeholder for summary API.</p>
            </div>
            <Badge>Coming later</Badge>
          </div>
          <div className="mt-6 grid h-64 grid-cols-6 items-end gap-3 rounded-lg bg-slate-50 p-4">
            {[42, 68, 55, 80, 61, 74].map((height, index) => (
              <div key={height} className="flex h-full flex-col justify-end gap-2">
                <div
                  className="rounded-t-lg bg-gradient-to-t from-violet-600 to-indigo-400"
                  style={{ height: `${height}%` }}
                />
                <span className="text-center text-xs font-semibold text-slate-500">
                  {["Jan", "Feb", "Mar", "Apr", "May", "Jun"][index]}
                </span>
              </div>
            ))}
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5 text-violet-600" aria-hidden="true" />
            <h2 className="text-base font-bold text-slate-950">Expenses by Category</h2>
          </div>
          <div className="mt-6 space-y-4">
            {[
              ["Food", "42%", "bg-violet-500"],
              ["Transport", "24%", "bg-cyan-500"],
              ["Utilities", "18%", "bg-emerald-500"],
              ["Shopping", "16%", "bg-amber-500"],
            ].map(([label, value, color]) => (
              <div key={label}>
                <div className="mb-2 flex items-center justify-between text-sm">
                  <span className="font-medium text-slate-700">{label}</span>
                  <span className="font-semibold text-slate-950">{value}</span>
                </div>
                <div className="h-2 rounded-full bg-slate-100">
                  <div className={`h-2 rounded-full ${color}`} style={{ width: value }} />
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_0.8fr]">
        <Card className="p-5">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-bold text-slate-950">Recent Expenses</h2>
            <Link
              href="/expenses"
              className="inline-flex items-center gap-1 text-sm font-semibold text-violet-700 hover:text-violet-800"
            >
              View all
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
          </div>
          <div className="divide-y divide-slate-100">
            {recentExpenses.map(([title, category, amount]) => (
              <div key={title} className="flex items-center justify-between gap-4 py-3">
                <div>
                  <p className="text-sm font-semibold text-slate-950">{title}</p>
                  <p className="mt-1 text-xs text-slate-500">{category}</p>
                </div>
                <p className="text-sm font-bold text-slate-950">{amount}</p>
              </div>
            ))}
          </div>
        </Card>

        <Card className="p-5">
          <h2 className="text-base font-bold text-slate-950">Quick Actions</h2>
          <div className="mt-4 grid gap-3 sm:grid-cols-2">
            <Link
              href="/expenses/new"
              className="flex items-center justify-between rounded-lg border border-violet-100 bg-violet-50 p-4 text-violet-700 transition hover:border-violet-200 hover:bg-violet-100"
            >
              <span className="flex items-center gap-3 text-sm font-semibold">
                <PlusCircle className="h-5 w-5" aria-hidden="true" />
                Add Expense
              </span>
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
            <Link
              href="/voice"
              className="flex items-center justify-between rounded-lg border border-cyan-100 bg-cyan-50 p-4 text-cyan-700 transition hover:border-cyan-200 hover:bg-cyan-100"
            >
              <span className="flex items-center gap-3 text-sm font-semibold">
                <Mic className="h-5 w-5" aria-hidden="true" />
                Voice Input
              </span>
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
          </div>
        </Card>
      </div>
    </AppShell>
  );
}
