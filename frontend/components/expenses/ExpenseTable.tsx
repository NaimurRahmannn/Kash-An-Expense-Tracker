import Link from "next/link";
import { Eye, Pencil, Trash2 } from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import {
  formatCurrency,
  getCategoryBadgeStyle,
} from "@/lib/expense-format";
import { formatDisplayDate } from "@/lib/date";
import type { Expense } from "@/types/expense";

type ExpenseTableProps = {
  expenses: Expense[];
  onDeletePlaceholder: () => void;
};

export function ExpenseTable({
  expenses,
  onDeletePlaceholder,
}: ExpenseTableProps) {
  return (
    <div className="overflow-hidden rounded-lg border border-slate-200">
      <div className="overflow-x-auto">
        <table className="min-w-230 divide-y divide-slate-200 text-sm">
          <thead className="bg-slate-50">
            <tr>
              {["Date", "Title", "Category", "Note", "Amount", "Actions"].map(
                (column) => (
                  <th
                    key={column}
                    scope="col"
                    className="px-4 py-3 text-left text-xs font-bold uppercase text-slate-500"
                  >
                    {column}
                  </th>
                ),
              )}
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 bg-white">
            {expenses.map((expense) => (
              <tr key={expense.id} className="transition hover:bg-slate-50">
                <td className="whitespace-nowrap px-4 py-4 font-medium text-slate-600">
                  {formatDisplayDate(expense.expense_date)}
                </td>
                <td className="min-w-52 px-4 py-4">
                  <p className="font-semibold text-slate-950">{expense.title}</p>
                </td>
                <td className="whitespace-nowrap px-4 py-4">
                  <Badge className={getCategoryBadgeStyle(expense.category)}>
                    {expense.category}
                  </Badge>
                </td>
                <td className="max-w-72 px-4 py-4 text-slate-600">
                  <p className="line-clamp-2">
                    {expense.note || "No note added"}
                  </p>
                </td>
                <td className="whitespace-nowrap px-4 py-4 font-bold text-slate-950">
                  {formatCurrency(expense.amount)}
                </td>
                <td className="whitespace-nowrap px-4 py-4">
                  <div className="flex items-center gap-1.5">
                    <button
                      type="button"
                      disabled
                      title="View page will be added later"
                      className="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 transition disabled:cursor-not-allowed disabled:opacity-60"
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
                      onClick={onDeletePlaceholder}
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
  );
}
