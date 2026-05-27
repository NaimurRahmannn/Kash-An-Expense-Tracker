"use client";

import { type FormEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import {
  ArrowRight,
  CalendarDays,
  CirclePlus,
  FileText,
  Loader2,
  Mic,
  ReceiptText,
  Star,
  TrendingUp,
  WalletCards,
} from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { Input } from "@/components/ui/Input";
import { useAuth } from "@/hooks/useAuth";
import {
  formatDisplayDate,
  getCurrentMonthRange,
} from "@/lib/date";
import { getExpenseSummary, getRecentExpenses } from "@/lib/expenses";
import { cn } from "@/lib/utils";
import type { Expense } from "@/types/expense";
import type { CategorySummary, ExpenseSummary } from "@/types/summary";

const moneyFormatter = new Intl.NumberFormat("en-BD", {
  maximumFractionDigits: 0,
});

const decimalMoneyFormatter = new Intl.NumberFormat("en-BD", {
  maximumFractionDigits: 2,
  minimumFractionDigits: 0,
});

const categoryColors = [
  "bg-violet-500",
  "bg-blue-500",
  "bg-emerald-500",
  "bg-amber-500",
  "bg-orange-500",
  "bg-pink-500",
  "bg-cyan-500",
  "bg-fuchsia-500",
  "bg-slate-500",
];

const statAccentStyles = {
  violet: "bg-violet-100 text-violet-700",
  blue: "bg-blue-100 text-blue-600",
  green: "bg-emerald-100 text-emerald-600",
  amber: "bg-amber-100 text-amber-600",
} as const;

type DateRange = {
  dateFrom: string;
  dateTo: string;
};

type StatCardProps = {
  title: string;
  value: string;
  detail: string;
  icon: typeof WalletCards;
  accent: keyof typeof statAccentStyles;
};

function getErrorMessage(error: unknown) {
  return error instanceof Error && error.message
    ? error.message
    : "Unable to connect to server. Please try again.";
}

function parseAPIDate(dateString: string) {
  return new Date(`${dateString}T00:00:00`);
}

function getInclusiveDays(dateFrom: string, dateTo: string) {
  const start = parseAPIDate(dateFrom).getTime();
  const end = parseAPIDate(dateTo).getTime();
  const dayMs = 1000 * 60 * 60 * 24;

  return Math.max(1, Math.floor((end - start) / dayMs) + 1);
}

function formatMoney(value: number, withDecimals = false) {
  const formatter = withDecimals ? decimalMoneyFormatter : moneyFormatter;

  return `৳${formatter.format(value)}`;
}

function StatCard({ accent, detail, icon: Icon, title, value }: StatCardProps) {
  return (
    <Card className="flex min-h-28 items-center gap-3 rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)] sm:p-5">
      <div
        className={cn(
          "flex h-13 w-13 shrink-0 items-center justify-center rounded-full sm:h-16 sm:w-16",
          statAccentStyles[accent],
        )}
      >
        <Icon className="h-7 w-7 sm:h-8 sm:w-8" aria-hidden="true" />
      </div>
      <div className="min-w-0">
        <p className="truncate text-xs font-semibold text-slate-950 sm:text-sm">
          {title}
        </p>
        <p className="mt-1 truncate text-xl font-bold tracking-[-0.01em] text-slate-950 sm:text-2xl">
          {value}
        </p>
        <p className="mt-1 truncate text-xs text-slate-500 sm:text-sm">{detail}</p>
      </div>
    </Card>
  );
}

function DashboardLoading() {
  return (
    <div className="space-y-6">
      <Card className="flex min-h-28 items-center gap-4 rounded-[14px] p-5">
        <Loader2 className="h-5 w-5 animate-spin text-violet-600" aria-hidden="true" />
        <p className="text-sm font-semibold text-slate-700">Loading dashboard...</p>
      </Card>
      <div className="grid grid-cols-2 gap-3 sm:gap-6 xl:grid-cols-4">
        {["one", "two", "three", "four"].map((item) => (
          <Card key={item} className="h-28 animate-pulse rounded-[14px] bg-slate-100" />
        ))}
      </div>
      <div className="grid gap-6 xl:grid-cols-[1fr_0.7fr]">
        <Card className="h-80 animate-pulse rounded-[14px] bg-slate-100" />
        <Card className="h-80 animate-pulse rounded-[14px] bg-slate-100" />
      </div>
    </div>
  );
}

function DashboardError({
  error,
  onRetry,
}: {
  error: string;
  onRetry: () => void;
}) {
  return (
    <Card className="rounded-[14px] border-rose-200 bg-rose-50 p-6">
      <h2 className="text-lg font-bold text-rose-900">
        Unable to load dashboard data
      </h2>
      <p className="mt-2 text-sm leading-6 text-rose-700">{error}</p>
      <Button className="mt-5" onClick={onRetry}>
        Retry
      </Button>
    </Card>
  );
}

function DateRangeCard({
  dateFrom,
  dateTo,
  isLoading,
  onApply,
  onDateFromChange,
  onDateToChange,
  validationError,
}: {
  dateFrom: string;
  dateTo: string;
  isLoading: boolean;
  onApply: (event: FormEvent<HTMLFormElement>) => void;
  onDateFromChange: (value: string) => void;
  onDateToChange: (value: string) => void;
  validationError: string;
}) {
  return (
    <Card className="rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
      <form className="grid gap-3 lg:grid-cols-[1fr_1fr_auto]" onSubmit={onApply}>
        <div>
          <label htmlFor="date-from" className="text-xs font-bold uppercase text-slate-500">
            Date from
          </label>
          <Input
            id="date-from"
            type="date"
            className="mt-2"
            value={dateFrom}
            onChange={(event) => onDateFromChange(event.target.value)}
          />
        </div>
        <div>
          <label htmlFor="date-to" className="text-xs font-bold uppercase text-slate-500">
            Date to
          </label>
          <Input
            id="date-to"
            type="date"
            className="mt-2"
            value={dateTo}
            onChange={(event) => onDateToChange(event.target.value)}
          />
        </div>
        <div className="flex items-end">
          <Button type="submit" disabled={isLoading} className="w-full lg:w-auto">
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
            ) : (
              <CalendarDays className="h-4 w-4" aria-hidden="true" />
            )}
            Apply
          </Button>
        </div>
      </form>
      {validationError ? (
        <p className="mt-3 text-sm font-medium text-rose-600">{validationError}</p>
      ) : null}
    </Card>
  );
}

function CategoryBreakdown({
  categories,
  totalAmount,
}: {
  categories: CategorySummary[];
  totalAmount: number;
}) {
  const sortedCategories = [...categories].sort((a, b) => b.total - a.total);

  if (sortedCategories.length === 0) {
    return (
      <EmptyState
        title="No category data"
        description="Add expenses to see category breakdown."
      />
    );
  }

  return (
    <div className="space-y-4">
      {sortedCategories.map((category, index) => {
        const percent = totalAmount > 0 ? (category.total / totalAmount) * 100 : 0;

        return (
          <div key={category.category}>
            <div className="mb-2 grid grid-cols-[1fr_auto] gap-3 text-sm">
              <div className="min-w-0">
                <p className="truncate font-semibold text-slate-950">
                  {category.category}
                </p>
                <p className="text-xs text-slate-500">
                  {category.count} {category.count === 1 ? "transaction" : "transactions"}
                </p>
              </div>
              <div className="text-right">
                <p className="font-bold text-slate-950">{formatMoney(category.total)}</p>
                <p className="text-xs text-slate-500">{percent.toFixed(1)}%</p>
              </div>
            </div>
            <div className="h-2 rounded-full bg-slate-100">
              <div
                className={cn(
                  "h-2 rounded-full",
                  categoryColors[index % categoryColors.length],
                )}
                style={{ width: `${Math.min(100, percent)}%` }}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
}

function RecentExpenses({ expenses }: { expenses: Expense[] }) {
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

function QuickActions() {
  return (
    <Card className="rounded-[14px] p-5 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
      <h2 className="text-lg font-bold text-slate-950">Quick Actions</h2>
      <div className="mt-4 grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
        <Link
          href="/expenses/new"
          className="flex h-16 items-center justify-between rounded-[14px] border border-violet-200 bg-violet-50 px-4 text-slate-950 transition hover:border-violet-300 hover:bg-violet-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-violet-600 text-white shadow-md shadow-violet-200">
              <CirclePlus className="h-6 w-6" aria-hidden="true" />
            </span>
            Add Expense
          </span>
          <ArrowRight className="h-5 w-5 text-slate-950" aria-hidden="true" />
        </Link>

        <Link
          href="/expenses"
          className="flex h-16 items-center justify-between rounded-[14px] border border-blue-200 bg-blue-50 px-4 text-slate-950 transition hover:border-blue-300 hover:bg-blue-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-100 text-blue-600">
              <FileText className="h-6 w-6" aria-hidden="true" />
            </span>
            View All Expenses
          </span>
          <ArrowRight className="h-5 w-5 text-slate-950" aria-hidden="true" />
        </Link>

        <Link
          href="/voice"
          className="flex h-16 items-center justify-between rounded-[14px] border border-emerald-200 bg-emerald-50 px-4 text-slate-950 transition hover:border-emerald-300 hover:bg-emerald-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
              <Mic className="h-6 w-6" aria-hidden="true" />
            </span>
            Voice Input
          </span>
          <Badge className="bg-violet-100 text-violet-700 ring-0">NEW</Badge>
        </Link>
      </div>
    </Card>
  );
}

export default function DashboardPage() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const defaultRange = useMemo(() => getCurrentMonthRange(), []);
  const [dateRange, setDateRange] = useState<DateRange>(defaultRange);
  const [appliedRange, setAppliedRange] = useState<DateRange>(defaultRange);
  const [summary, setSummary] = useState<ExpenseSummary | null>(null);
  const [recentExpenses, setRecentExpenses] = useState<Expense[]>([]);
  const [isDashboardLoading, setIsDashboardLoading] = useState(true);
  const [isSummaryLoading, setIsSummaryLoading] = useState(false);
  const [error, setError] = useState("");
  const [dateError, setDateError] = useState("");

  const loadDashboard = async () => {
    setError("");
    setIsDashboardLoading(true);

    try {
      const [summaryResponse, expensesResponse] = await Promise.all([
        getExpenseSummary({
          dateFrom: appliedRange.dateFrom,
          dateTo: appliedRange.dateTo,
        }),
        getRecentExpenses(5),
      ]);

      if (!summaryResponse.success || !summaryResponse.data) {
        throw new Error(summaryResponse.message || "Unable to load summary");
      }

      if (!expensesResponse.success) {
        throw new Error(expensesResponse.message || "Unable to load expenses");
      }

      setSummary(summaryResponse.data);
      setRecentExpenses(expensesResponse.data ?? []);
    } catch (caughtError) {
      setError(getErrorMessage(caughtError));
    } finally {
      setIsDashboardLoading(false);
    }
  };

  const loadSummary = async (range: DateRange) => {
    setError("");
    setIsSummaryLoading(true);

    try {
      const response = await getExpenseSummary({
        dateFrom: range.dateFrom,
        dateTo: range.dateTo,
      });

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to load summary");
      }

      setSummary(response.data);
      setAppliedRange(range);
    } catch (caughtError) {
      setError(getErrorMessage(caughtError));
    } finally {
      setIsSummaryLoading(false);
    }
  };

  useEffect(() => {
    if (!isAuthLoading && isAuthenticated) {
      const timeoutId = window.setTimeout(() => {
        void loadDashboard();
      }, 0);

      return () => window.clearTimeout(timeoutId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthLoading, isAuthenticated]);

  function validateRange() {
    if (!dateRange.dateFrom) {
      return "Date from is required";
    }

    if (!dateRange.dateTo) {
      return "Date to is required";
    }

    if (dateRange.dateFrom > dateRange.dateTo) {
      return "Date from cannot be after date to";
    }

    return "";
  }

  function handleApplyDateRange(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationError = validateRange();
    setDateError(validationError);

    if (validationError) {
      return;
    }

    void loadSummary(dateRange);
  }

  const highestCategory = useMemo(() => {
    if (!summary?.by_category.length) {
      return null;
    }

    return [...summary.by_category].sort((a, b) => b.total - a.total)[0];
  }, [summary]);

  const totalAmount = summary?.total_amount ?? 0;
  const totalCount = summary?.total_count ?? 0;
  const dayCount = summary
    ? getInclusiveDays(summary.date_from, summary.date_to)
    : getInclusiveDays(appliedRange.dateFrom, appliedRange.dateTo);
  const dailyAverage = totalAmount / dayCount;
  const highestCategoryPercent =
    highestCategory && totalAmount > 0
      ? (highestCategory.total / totalAmount) * 100
      : 0;

  const statCards: StatCardProps[] = [
    {
      title: "Total Expenses",
      value: formatMoney(totalAmount, true),
      detail: "In selected range",
      icon: WalletCards,
      accent: "violet",
    },
    {
      title: "Total Transactions",
      value: String(totalCount),
      detail: "In selected range",
      icon: ReceiptText,
      accent: "blue",
    },
    {
      title: "Daily Average",
      value: formatMoney(dailyAverage, true),
      detail: `${dayCount} ${dayCount === 1 ? "day" : "days"}`,
      icon: TrendingUp,
      accent: "green",
    },
    {
      title: "Highest Category",
      value: highestCategory?.category ?? "No data",
      detail: highestCategory
        ? `${formatMoney(highestCategory.total)} (${highestCategoryPercent.toFixed(1)}%)`
        : "Add expenses to compare",
      icon: Star,
      accent: "amber",
    },
  ];

  return (
    <AppShell title="Dashboard">
      <div className="space-y-5 sm:space-y-6">
        <DateRangeCard
          dateFrom={dateRange.dateFrom}
          dateTo={dateRange.dateTo}
          isLoading={isSummaryLoading}
          validationError={dateError}
          onApply={handleApplyDateRange}
          onDateFromChange={(value) => setDateRange((range) => ({ ...range, dateFrom: value }))}
          onDateToChange={(value) => setDateRange((range) => ({ ...range, dateTo: value }))}
        />

        {isDashboardLoading ? (
          <DashboardLoading />
        ) : error ? (
          <DashboardError error={error} onRetry={() => void loadDashboard()} />
        ) : (
          <>
            <Card className="rounded-[14px] border-violet-200/80 bg-linear-to-r from-white via-violet-50/50 to-white p-5 shadow-[0_8px_24px_rgba(109,53,242,0.10)]">
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h2 className="text-lg font-bold text-slate-950">
                    Selected range summary
                  </h2>
                  <p className="mt-1 text-sm text-slate-500">
                    {formatDisplayDate(summary?.date_from ?? appliedRange.dateFrom)} -{" "}
                    {formatDisplayDate(summary?.date_to ?? appliedRange.dateTo)}
                  </p>
                </div>
                {isSummaryLoading ? (
                  <div className="inline-flex items-center gap-2 text-sm font-semibold text-violet-700">
                    <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
                    Refreshing summary...
                  </div>
                ) : null}
              </div>
            </Card>

            <div className="grid grid-cols-2 gap-3 sm:gap-6 xl:grid-cols-4">
              {statCards.map((item) => (
                <StatCard key={item.title} {...item} />
              ))}
            </div>

            <div className="grid gap-6 xl:grid-cols-[1fr_0.7fr]">
              <Card className="rounded-[14px] p-5 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
                <div className="mb-5 flex items-center justify-between gap-3">
                  <div>
                    <h2 className="text-lg font-bold text-slate-950">
                      Expenses by Category
                    </h2>
                    <p className="mt-1 text-sm text-slate-500">
                      Real breakdown from the summary API
                    </p>
                  </div>
                  <Badge>{summary?.by_category.length ?? 0} categories</Badge>
                </div>
                <CategoryBreakdown
                  categories={summary?.by_category ?? []}
                  totalAmount={totalAmount}
                />
              </Card>

              <QuickActions />
            </div>

            <Card className="rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
              <div className="flex items-center justify-between px-2 pb-3">
                <div>
                  <h2 className="text-lg font-bold text-slate-950">Recent Expenses</h2>
                  <p className="mt-1 text-sm text-slate-500">
                    Latest entries from your expense list
                  </p>
                </div>
                <Link
                  href="/expenses"
                  className="text-sm font-semibold text-violet-700 hover:text-violet-800"
                >
                  View all
                </Link>
              </div>
              <RecentExpenses expenses={recentExpenses} />
            </Card>
          </>
        )}
      </div>
    </AppShell>
  );
}
