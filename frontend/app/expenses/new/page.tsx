"use client";

import Link from "next/link";
import { ArrowLeft, CheckCircle2, Lightbulb, RotateCcw, Save, Tag } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useCreateExpenseForm } from "@/hooks/useCreateExpenseForm";
import {
  expenseCategories,
  getCategoryBadgeStyle,
} from "@/lib/expense-format";
import type { ExpenseCategory } from "@/types/expense";

const fieldClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-violet-400 focus:ring-4 focus:ring-violet-500/15";

function FieldError({ message }: { message?: string }) {
  return message ? (
    <p className="mt-2 text-sm font-medium text-rose-600">{message}</p>
  ) : null;
}

export default function AddExpensePage() {
  const form = useCreateExpenseForm();

  return (
    <AppShell title="Add Expense">
      <PageHeader
        title="Add Expense"
        description="Record a new expense in your tracker."
        action={
          <Link
            href="/expenses"
            className="inline-flex h-10 items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-50"
          >
            <ArrowLeft className="h-4 w-4" aria-hidden="true" />
            Back to Expenses
          </Link>
        }
      />

      <div className="grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card className="p-5 sm:p-6">
          <form className="space-y-5" onSubmit={form.handleSubmit}>
            {form.formError ? (
              <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
                {form.formError}
              </div>
            ) : null}

            {form.successMessage ? (
              <div className="flex items-center gap-2 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
                <CheckCircle2 className="h-4 w-4" aria-hidden="true" />
                {form.successMessage}
              </div>
            ) : null}

            <div className="grid gap-5 md:grid-cols-2">
              <div>
                <label htmlFor="title" className="text-sm font-semibold text-slate-700">
                  Title
                </label>
                <Input
                  id="title"
                  className="mt-2"
                  placeholder="Lunch at office"
                  value={form.values.title}
                  onChange={(event) => form.updateField("title", event.target.value)}
                />
                <FieldError message={form.fieldErrors.title} />
              </div>

              <div>
                <label htmlFor="amount" className="text-sm font-semibold text-slate-700">
                  Amount
                </label>
                <Input
                  id="amount"
                  className="mt-2"
                  type="number"
                  min="0.01"
                  step="0.01"
                  placeholder="350.50"
                  value={form.values.amount}
                  onChange={(event) => form.updateField("amount", event.target.value)}
                />
                <FieldError message={form.fieldErrors.amount} />
              </div>

              <div>
                <label htmlFor="category" className="text-sm font-semibold text-slate-700">
                  Category
                </label>
                <select
                  id="category"
                  className={`${fieldClass} mt-2 h-10`}
                  value={form.values.category}
                  onChange={(event) =>
                    form.updateField(
                      "category",
                      event.target.value as ExpenseCategory | "",
                    )
                  }
                >
                  <option value="">Select category</option>
                  {expenseCategories.map((category) => (
                    <option key={category} value={category}>
                      {category}
                    </option>
                  ))}
                </select>
                <FieldError message={form.fieldErrors.category} />
              </div>

              <div>
                <label
                  htmlFor="expense_date"
                  className="text-sm font-semibold text-slate-700"
                >
                  Expense date
                </label>
                <Input
                  id="expense_date"
                  className="mt-2"
                  type="date"
                  value={form.values.expense_date}
                  onChange={(event) =>
                    form.updateField("expense_date", event.target.value)
                  }
                />
                <FieldError message={form.fieldErrors.expense_date} />
              </div>
            </div>

            <div>
              <label htmlFor="note" className="text-sm font-semibold text-slate-700">
                Note
              </label>
              <textarea
                id="note"
                rows={5}
                placeholder="Optional notes about this expense"
                className={`${fieldClass} mt-2 resize-none`}
                value={form.values.note}
                onChange={(event) => form.updateField("note", event.target.value)}
              />
            </div>

            <div className="flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:items-center">
              <Button type="submit" disabled={form.isSubmitting}>
                <Save className="h-4 w-4" aria-hidden="true" />
                {form.isSubmitting ? "Saving..." : "Save Expense"}
              </Button>
              <Button
                type="button"
                variant="secondary"
                disabled={form.isSubmitting}
                onClick={form.resetForm}
              >
                <RotateCcw className="h-4 w-4" aria-hidden="true" />
                Reset
              </Button>
              <Link
                href="/expenses"
                className="inline-flex h-10 items-center justify-center rounded-lg px-4 text-sm font-semibold text-slate-600 transition hover:bg-violet-50 hover:text-violet-700"
              >
                Cancel
              </Link>
            </div>
          </form>
        </Card>

        <div className="space-y-6">
          <Card className="p-5">
            <div className="flex items-center gap-2">
              <Lightbulb className="h-5 w-5 text-amber-500" aria-hidden="true" />
              <h2 className="text-base font-bold text-slate-950">Quick Tips</h2>
            </div>
            <div className="mt-4 space-y-3 text-sm leading-6 text-slate-600">
              <p>Use clear titles like &quot;Lunch at office&quot;.</p>
              <p>Choose the correct category so dashboard summaries stay useful.</p>
              <p>Add a note when the expense needs extra context.</p>
            </div>
          </Card>

          <Card className="p-5">
            <div className="flex items-center gap-2">
              <Tag className="h-5 w-5 text-violet-600" aria-hidden="true" />
              <h2 className="text-base font-bold text-slate-950">
                Allowed Categories
              </h2>
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              {expenseCategories.map((category) => (
                <Badge
                  key={category}
                  className={getCategoryBadgeStyle(category)}
                >
                  {category}
                </Badge>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
