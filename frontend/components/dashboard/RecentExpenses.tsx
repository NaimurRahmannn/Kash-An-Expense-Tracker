import Link from "next/link";
import { ArrowRight, CirclePlus } from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import { EmptyState } from "@/components/ui/EmptyState";
import { formatDisplayDate } from "@/lib/date";
import { formatMoney } from "@/lib/dashboard-format";
import type { Expense } from "@/types/expense";

export function RecentExpenses({ expenses }: { expenses: Expense[] }) {
  if (expenses.length === 0) {
    return (
      <EmptyState
        title="No recent expenses"
        description="Create your first expense to see it here."
        action={
          <Link
            href="/expenses/new"
            className="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
          >
            <CirclePlus className="h-4 w-4" aria-hidden="true" />
            Add Expense
          </Link>
        }
      />
    );
  }

  return (
    <>
      <div className="space-y-0 lg:hidden">
        {expenses.map((expense) => (
          <Link
            key={`${expense.id}-mobile`}
            href="/expenses"
            className="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-3 border-b border-slate-200 px-2 py-3 last:border-b-0"
          >
            <span className="min-w-0">
              <span className="block truncate text-base font-bold text-slate-950">
                {expense.title}
              </span>
              <span className="block truncate text-sm font-medium text-slate-500">
                {formatDisplayDate(expense.expense_date)}
              </span>
              <Badge className="mt-2">{expense.category}</Badge>
            </span>
            <span className="flex items-center gap-3 text-base font-bold text-slate-950">
              {formatMoney(expense.amount)}
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </span>
          </Link>
        ))}
      </div>

      <div className="hidden overflow-hidden rounded-lg border border-slate-200 lg:block">
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-slate-50">
              <tr className="border-b border-slate-200 text-left text-slate-700">
                <th className="px-3 py-2 font-semibold">Date</th>
                <th className="px-3 py-2 font-semibold">Title</th>
                <th className="px-3 py-2 font-semibold">Category</th>
                <th className="px-3 py-2 font-semibold">Amount</th>
              </tr>
            </thead>
            <tbody>
              {expenses.map((expense) => (
                <tr
                  key={expense.id}
                  className="border-b border-slate-200 last:border-b-0"
                >
                  <td className="whitespace-nowrap px-3 py-3 text-slate-800">
                    {formatDisplayDate(expense.expense_date)}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3 font-medium text-slate-900">
                    {expense.title}
                  </td>
                  <td className="whitespace-nowrap px-3 py-3">
                    <Badge>{expense.category}</Badge>
                  </td>
                  <td className="whitespace-nowrap px-3 py-3 font-semibold text-slate-950">
                    {formatMoney(expense.amount)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}
