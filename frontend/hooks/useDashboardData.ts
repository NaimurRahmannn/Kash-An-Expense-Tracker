"use client";

import { type FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { ReceiptText, Star, TrendingUp, WalletCards } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { getCurrentMonthRange } from "@/lib/date";
import {
  formatMoney,
  getErrorMessage,
  getHighestCategory,
  getInclusiveDays,
} from "@/lib/dashboard-format";
import { getExpenseSummary, getRecentExpenses } from "@/lib/expenses";
import type { DateRange, DashboardStat } from "@/types/dashboard";
import type { Expense } from "@/types/expense";
import type { ExpenseSummary } from "@/types/summary";

function validateDateRange(dateRange: DateRange) {
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

export function useDashboardData() {
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

  const loadDashboard = useCallback(async () => {
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
  }, [appliedRange.dateFrom, appliedRange.dateTo]);

  const loadSummary = useCallback(async (range: DateRange) => {
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
  }, []);

  useEffect(() => {
    if (!isAuthLoading && isAuthenticated) {
      const timeoutId = window.setTimeout(() => {
        void loadDashboard();
      }, 0);

      return () => window.clearTimeout(timeoutId);
    }
  }, [isAuthLoading, isAuthenticated, loadDashboard]);

  function handleApplyDateRange(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationError = validateDateRange(dateRange);
    setDateError(validationError);

    if (validationError) {
      return;
    }

    void loadSummary(dateRange);
  }

  const highestCategory = useMemo(
    () => getHighestCategory(summary?.by_category ?? []),
    [summary],
  );

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

  const statCards: DashboardStat[] = [
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

  return {
    appliedRange,
    dateError,
    dateRange,
    error,
    isDashboardLoading,
    isSummaryLoading,
    recentExpenses,
    statCards,
    summary,
    totalAmount,
    handleApplyDateRange,
    loadDashboard,
    setDateFrom: (dateFrom: string) =>
      setDateRange((range) => ({ ...range, dateFrom })),
    setDateTo: (dateTo: string) =>
      setDateRange((range) => ({ ...range, dateTo })),
  };
}
