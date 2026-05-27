"use client";

import { type FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { getErrorMessage } from "@/lib/dashboard-format";
import { expenseMatchesSearch } from "@/lib/expense-format";
import { getExpenses } from "@/lib/expenses";
import type {
  Expense,
  ExpenseCategory,
  ExpenseListQuery,
  ExpenseSortBy,
  SortOrder,
} from "@/types/expense";

const DEFAULT_EXPENSE_LIMIT = 8;
const DEFAULT_SORT_BY: ExpenseSortBy = "expense_date";
const DEFAULT_SORT_ORDER: SortOrder = "desc";

function getDefaultQuery(limit: number): ExpenseListQuery {
  return {
    page: 1,
    limit,
    sortBy: DEFAULT_SORT_BY,
    sortOrder: DEFAULT_SORT_ORDER,
  };
}

function validateDateRange(dateFrom: string, dateTo: string) {
  if (dateFrom && dateTo && dateFrom > dateTo) {
    return "date_from cannot be after date_to";
  }

  return "";
}

export function useExpensesList() {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [filterError, setFilterError] = useState("");
  const [page, setPage] = useState(1);
  const [limit] = useState(DEFAULT_EXPENSE_LIMIT);
  const [searchTerm, setSearchTerm] = useState("");
  const [category, setCategory] = useState<ExpenseCategory | "">("");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [sortBy, setSortBy] = useState<ExpenseSortBy | "">(DEFAULT_SORT_BY);
  const [sortOrder, setSortOrder] = useState<SortOrder>(DEFAULT_SORT_ORDER);
  const [appliedQuery, setAppliedQuery] = useState<ExpenseListQuery>(() =>
    getDefaultQuery(DEFAULT_EXPENSE_LIMIT),
  );

  const requestExpenses = useCallback(async (query: ExpenseListQuery) => {
    setAppliedQuery(query);
    setError("");
    setIsLoading(true);

    try {
      const response = await getExpenses(query);

      if (!response.success) {
        throw new Error(response.message || "Unable to load expenses");
      }

      setExpenses(response.data ?? []);
    } catch (caughtError) {
      setExpenses([]);
      setError(getErrorMessage(caughtError));
    } finally {
      setIsLoading(false);
    }
  }, []);

  const buildQuery = useCallback(
    (pageValue: number): ExpenseListQuery => ({
      page: pageValue,
      limit,
      category,
      dateFrom,
      dateTo,
      sortBy,
      sortOrder,
    }),
    [category, dateFrom, dateTo, limit, sortBy, sortOrder],
  );

  useEffect(() => {
    if (!isAuthLoading && isAuthenticated) {
      const timeoutId = window.setTimeout(() => {
        void requestExpenses(getDefaultQuery(limit));
      }, 0);

      return () => window.clearTimeout(timeoutId);
    }
  }, [isAuthLoading, isAuthenticated, limit, requestExpenses]);

  const visibleExpenses = useMemo(
    () => expenses.filter((expense) => expenseMatchesSearch(expense, searchTerm)),
    [expenses, searchTerm],
  );

  function handleApplyFilters(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationError = validateDateRange(dateFrom, dateTo);
    setFilterError(validationError);

    if (validationError) {
      return;
    }

    const nextPage = 1;
    setPage(nextPage);
    void requestExpenses(buildQuery(nextPage));
  }

  function handleResetFilters() {
    const nextPage = 1;
    const query = getDefaultQuery(limit);

    setSearchTerm("");
    setCategory("");
    setDateFrom("");
    setDateTo("");
    setSortBy(DEFAULT_SORT_BY);
    setSortOrder(DEFAULT_SORT_ORDER);
    setFilterError("");
    setPage(nextPage);
    void requestExpenses(query);
  }

  function handlePreviousPage() {
    const nextPage = Math.max(1, page - 1);
    setPage(nextPage);
    void requestExpenses(buildQuery(nextPage));
  }

  function handleNextPage() {
    const nextPage = page + 1;
    setPage(nextPage);
    void requestExpenses(buildQuery(nextPage));
  }

  return {
    appliedQuery,
    canGoNext: expenses.length === limit,
    canGoPrevious: page > 1,
    category,
    dateFrom,
    dateTo,
    error,
    expenses,
    filterError,
    handleApplyFilters,
    handleNextPage,
    handlePreviousPage,
    handleResetFilters,
    isLoading,
    limit,
    page,
    searchTerm,
    setCategory,
    setDateFrom,
    setDateTo,
    setSearchTerm,
    setSortBy,
    setSortOrder,
    sortBy,
    sortOrder,
    visibleExpenses,
    retry: () => void requestExpenses(appliedQuery),
  };
}
