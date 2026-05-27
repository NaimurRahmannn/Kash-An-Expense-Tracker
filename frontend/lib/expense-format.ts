import type { Expense, ExpenseCategory } from "@/types/expense";

export const expenseCategories: ExpenseCategory[] = [
  "Food",
  "Transport",
  "Housing",
  "Entertainment",
  "Shopping",
  "Healthcare",
  "Education",
  "Utilities",
  "Other",
];

const categoryBadgeStyles: Record<ExpenseCategory, string> = {
  Food: "bg-violet-50 text-violet-700 ring-violet-200/80",
  Transport: "bg-blue-50 text-blue-700 ring-blue-200/80",
  Housing: "bg-emerald-50 text-emerald-700 ring-emerald-200/80",
  Entertainment: "bg-amber-50 text-amber-700 ring-amber-200/80",
  Shopping: "bg-orange-50 text-orange-700 ring-orange-200/80",
  Healthcare: "bg-rose-50 text-rose-700 ring-rose-200/80",
  Education: "bg-cyan-50 text-cyan-700 ring-cyan-200/80",
  Utilities: "bg-teal-50 text-teal-700 ring-teal-200/80",
  Other: "bg-slate-100 text-slate-700 ring-slate-200/80",
};

export function formatCurrency(amount: number) {
  const hasDecimal = !Number.isInteger(amount);
  const formatter = new Intl.NumberFormat("en-BD", {
    minimumFractionDigits: hasDecimal ? 2 : 0,
    maximumFractionDigits: 2,
  });

  return `\u09F3${formatter.format(amount)}`;
}

export function getCategoryBadgeStyle(category: ExpenseCategory) {
  return categoryBadgeStyles[category];
}

export function expenseMatchesSearch(expense: Expense, searchTerm: string) {
  const query = searchTerm.trim().toLowerCase();

  if (!query) {
    return true;
  }

  return [expense.title, expense.note, expense.category].some((value) =>
    value.toLowerCase().includes(query),
  );
}
