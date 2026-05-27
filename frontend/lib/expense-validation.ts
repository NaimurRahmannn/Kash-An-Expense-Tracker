import { formatDateForAPI } from "@/lib/date";
import type { Expense, ExpenseCategory, ExpenseInput } from "@/types/expense";

export type ExpenseFormValues = {
  title: string;
  amount: string;
  category: ExpenseCategory | "";
  expense_date: string;
  note: string;
};

export type ExpenseFormErrors = Partial<
  Record<keyof ExpenseFormValues, string>
>;

const defaultCategory: ExpenseCategory = "Food";

export function getDefaultExpenseFormValues(): ExpenseFormValues {
  return {
    title: "",
    amount: "",
    category: defaultCategory,
    expense_date: formatDateForAPI(new Date()),
    note: "",
  };
}

export function expenseToFormValues(expense: Expense): ExpenseFormValues {
  return {
    title: expense.title,
    amount: String(expense.amount),
    category: expense.category,
    expense_date: expense.expense_date,
    note: expense.note,
  };
}

function isValidAPIDate(dateString: string) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(dateString)) {
    return false;
  }

  const [year, month, day] = dateString.split("-").map(Number);
  const parsedDate = new Date(year, month - 1, day);

  return (
    parsedDate.getFullYear() === year &&
    parsedDate.getMonth() === month - 1 &&
    parsedDate.getDate() === day
  );
}

export function validateExpenseForm(values: ExpenseFormValues) {
  const errors: ExpenseFormErrors = {};

  if (!values.title.trim()) {
    errors.title = "Title is required";
  }

  if (!values.amount.trim()) {
    errors.amount = "Amount is required";
  } else {
    const amount = Number(values.amount);

    if (!Number.isFinite(amount) || amount <= 0) {
      errors.amount = "Amount must be positive";
    }
  }

  if (!values.category) {
    errors.category = "Category is required";
  }

  if (!values.expense_date) {
    errors.expense_date = "Expense date is required";
  } else if (!isValidAPIDate(values.expense_date)) {
    errors.expense_date = "Invalid expense date format";
  }

  return errors;
}

export function hasExpenseFormErrors(errors: ExpenseFormErrors) {
  return Object.keys(errors).length > 0;
}

export function toExpenseInput(values: ExpenseFormValues): ExpenseInput {
  return {
    title: values.title.trim(),
    amount: Number(values.amount),
    category: values.category as ExpenseCategory,
    note: values.note.trim(),
    expense_date: values.expense_date,
  };
}
