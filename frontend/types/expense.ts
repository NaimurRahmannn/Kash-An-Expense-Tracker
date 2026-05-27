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
