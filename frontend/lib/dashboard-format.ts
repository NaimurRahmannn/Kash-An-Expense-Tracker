import type { CategorySummary } from "@/types/summary";

const moneyFormatter = new Intl.NumberFormat("en-BD", {
  maximumFractionDigits: 0,
});

const decimalMoneyFormatter = new Intl.NumberFormat("en-BD", {
  maximumFractionDigits: 2,
  minimumFractionDigits: 0,
});

function parseAPIDate(dateString: string) {
  return new Date(`${dateString}T00:00:00`);
}

export function formatMoney(value: number, withDecimals = false) {
  const formatter = withDecimals ? decimalMoneyFormatter : moneyFormatter;

  return `\u09F3${formatter.format(value)}`;
}

export function getErrorMessage(error: unknown) {
  return error instanceof Error && error.message
    ? error.message
    : "Unable to connect to server. Please try again.";
}

export function getInclusiveDays(dateFrom: string, dateTo: string) {
  const start = parseAPIDate(dateFrom).getTime();
  const end = parseAPIDate(dateTo).getTime();
  const dayMs = 1000 * 60 * 60 * 24;

  return Math.max(1, Math.floor((end - start) / dayMs) + 1);
}

export function getCategoryPercent(categoryTotal: number, totalAmount: number) {
  return totalAmount > 0 ? (categoryTotal / totalAmount) * 100 : 0;
}

export function getHighestCategory(categories: CategorySummary[]) {
  if (categories.length === 0) {
    return null;
  }

  return [...categories].sort((a, b) => b.total - a.total)[0];
}

export function sortCategoriesByTotal(categories: CategorySummary[]) {
  return [...categories].sort((a, b) => b.total - a.total);
}
