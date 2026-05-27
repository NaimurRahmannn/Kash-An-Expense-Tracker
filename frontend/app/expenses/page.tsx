import Link from "next/link";
import {
  CalendarDays,
  ChevronLeft,
  ChevronRight,
  Eye,
  Pencil,
  PlusCircle,
  Trash2,
} from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";

const selectClass =
  "h-10 rounded-lg border border-slate-200 bg-white px-3 text-sm font-medium text-slate-600 outline-none transition focus:border-violet-400 focus:ring-4 focus:ring-violet-500/15";

const expenses = [
  {
    id: 101,
    date: "May 27, 2026",
    title: "Lunch at Pizza Hut",
    category: "Food",
    amount: "৳560",
  },
  {
    id: 102,
    date: "May 26, 2026",
    title: "Metro card recharge",
    category: "Transport",
    amount: "৳1,200",
  },
  {
    id: 103,
    date: "May 25, 2026",
    title: "Monthly internet bill",
    category: "Utilities",
    amount: "৳1,050",
  },
  {
    id: 104,
    date: "May 24, 2026",
    title: "Pharmacy purchase",
    category: "Healthcare",
    amount: "৳780",
  },
  {
    id: 105,
    date: "May 23, 2026",
    title: "Notebook and pens",
    category: "Education",
    amount: "৳420",
  },
];

export default function ExpensesPage() {
  return (
    <AppShell>
      <PageHeader
        title="Expenses"
        description="Static expense list foundation. Filtering, sorting, pagination, view, edit, and delete will be wired later."
        action={
          <Link
            href="/expenses/new"
            className="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
          >
            <PlusCircle className="h-4 w-4" aria-hidden="true" />
            Add Expense
          </Link>
        }
      />

      <Card className="p-5">
        <div className="grid gap-3 lg:grid-cols-[1fr_auto_auto_auto]">
          <Input placeholder="Search by title or note..." aria-label="Search expenses table" />
          <Button variant="secondary">
            <CalendarDays className="h-4 w-4" aria-hidden="true" />
            Date range
          </Button>
          <select className={selectClass} defaultValue="">
            <option value="" disabled>
              Category
            </option>
            <option>Food</option>
            <option>Transport</option>
            <option>Utilities</option>
          </select>
          <select className={selectClass} defaultValue="">
            <option value="" disabled>
              Sort by
            </option>
            <option>Newest first</option>
            <option>Amount high to low</option>
            <option>Amount low to high</option>
          </select>
        </div>

        <div className="mt-5 overflow-hidden rounded-lg border border-slate-200">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-slate-200 text-sm">
              <thead className="bg-slate-50">
                <tr>
                  {["Date", "Title", "Category", "Amount", "Actions"].map((column) => (
                    <th
                      key={column}
                      scope="col"
                      className="px-4 py-3 text-left text-xs font-bold uppercase text-slate-500"
                    >
                      {column}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {expenses.map((expense) => (
                  <tr key={expense.id} className="transition hover:bg-slate-50">
                    <td className="whitespace-nowrap px-4 py-4 font-medium text-slate-600">
                      {expense.date}
                    </td>
                    <td className="min-w-52 px-4 py-4 font-semibold text-slate-950">
                      {expense.title}
                    </td>
                    <td className="whitespace-nowrap px-4 py-4">
                      <Badge>{expense.category}</Badge>
                    </td>
                    <td className="whitespace-nowrap px-4 py-4 font-bold text-slate-950">
                      {expense.amount}
                    </td>
                    <td className="whitespace-nowrap px-4 py-4">
                      <div className="flex items-center gap-1.5">
                        <button
                          type="button"
                          className="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-500 transition hover:bg-slate-100 hover:text-slate-900"
                          aria-label={`View ${expense.title}`}
                        >
                          <Eye className="h-4 w-4" aria-hidden="true" />
                        </button>
                        <Link
                          href={`/expenses/${expense.id}/edit`}
                          className="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-500 transition hover:bg-violet-50 hover:text-violet-700"
                          aria-label={`Edit ${expense.title}`}
                        >
                          <Pencil className="h-4 w-4" aria-hidden="true" />
                        </Link>
                        <button
                          type="button"
                          className="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-500 transition hover:bg-rose-50 hover:text-rose-600"
                          aria-label={`Delete ${expense.title}`}
                        >
                          <Trash2 className="h-4 w-4" aria-hidden="true" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-sm text-slate-500">Showing 1 to 5 of 24 placeholder expenses</p>
          <div className="flex items-center gap-2">
            <Button variant="secondary" disabled>
              <ChevronLeft className="h-4 w-4" aria-hidden="true" />
            </Button>
            <span className="rounded-lg bg-violet-50 px-3 py-2 text-sm font-bold text-violet-700">
              1
            </span>
            <Button variant="secondary">
              <ChevronRight className="h-4 w-4" aria-hidden="true" />
            </Button>
          </div>
        </div>
      </Card>
    </AppShell>
  );
}
