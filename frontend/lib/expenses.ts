import { apiRequest } from "@/lib/api";
import { getStoredUserID } from "@/lib/auth-storage";
import type { ApiResponse } from "@/types/api";
import type { Expense, ExpenseInput, ExpenseListQuery } from "@/types/expense";
import type { ExpenseSummary } from "@/types/summary";

function getUserHeaders() {
  const userID = getStoredUserID();

  if (!userID) {
    throw new Error("User is not authenticated");
  }

  return {
    "X-User-ID": String(userID),
  };
}

export function getExpenseSummary(params: {
  dateFrom: string;
  dateTo: string;
}): Promise<ApiResponse<ExpenseSummary>> {
  const searchParams = new URLSearchParams({
    date_from: params.dateFrom,
    date_to: params.dateTo,
  });

  return apiRequest<ApiResponse<ExpenseSummary>>(
    `/expenses/summary?${searchParams.toString()}`,
    {
      headers: getUserHeaders(),
    },
  );
}

export function getRecentExpenses(
  limit = 5,
): Promise<ApiResponse<Expense[]>> {
  const searchParams = new URLSearchParams({
    page: "1",
    limit: String(limit),
    sort_by: "expense_date",
    sort_order: "desc",
  });

  return apiRequest<ApiResponse<Expense[]>>(
    `/expenses?${searchParams.toString()}`,
    {
      headers: getUserHeaders(),
    },
  );
}

export function getExpenses(
  query: ExpenseListQuery,
): Promise<ApiResponse<Expense[]>> {
  const searchParams = new URLSearchParams();

  if (query.page) {
    searchParams.set("page", String(query.page));
  }

  if (query.limit) {
    searchParams.set("limit", String(query.limit));
  }

  if (query.category) {
    searchParams.set("category", query.category);
  }

  if (query.dateFrom) {
    searchParams.set("date_from", query.dateFrom);
  }

  if (query.dateTo) {
    searchParams.set("date_to", query.dateTo);
  }

  if (query.sortBy) {
    searchParams.set("sort_by", query.sortBy);
  }

  if (query.sortOrder) {
    searchParams.set("sort_order", query.sortOrder);
  }

  const queryString = searchParams.toString();

  return apiRequest<ApiResponse<Expense[]>>(
    `/expenses${queryString ? `?${queryString}` : ""}`,
    {
      headers: getUserHeaders(),
    },
  );
}

export function createExpense(
  input: ExpenseInput,
): Promise<ApiResponse<Expense>> {
  return apiRequest<ApiResponse<Expense>>("/expenses", {
    method: "POST",
    headers: getUserHeaders(),
    body: JSON.stringify(input),
  });
}
