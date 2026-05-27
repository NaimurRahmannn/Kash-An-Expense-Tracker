import { apiRequest } from "@/lib/api";
import { getStoredUserID } from "@/lib/auth-storage";
import type { ApiResponse } from "@/types/api";
import type { Expense } from "@/types/expense";
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
