export type ExpenseCategory =
  | "Food"
  | "Transport"
  | "Housing"
  | "Entertainment"
  | "Shopping"
  | "Healthcare"
  | "Education"
  | "Utilities"
  | "Other";

export type Expense = {
  id: number;
  title: string;
  amount: number;
  category: ExpenseCategory;
  note: string;
  expense_date: string;
};

export type ExpenseInput = {
  title: string;
  amount: number;
  category: ExpenseCategory;
  note: string;
  expense_date: string;
};

export type ExpenseSortBy = "amount" | "expense_date";

export type SortOrder = "asc" | "desc";

export type ExpenseListQuery = {
  page?: number;
  limit?: number;
  category?: ExpenseCategory | "";
  dateFrom?: string;
  dateTo?: string;
  sortBy?: ExpenseSortBy | "";
  sortOrder?: SortOrder;
};
