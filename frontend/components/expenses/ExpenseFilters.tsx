import { RotateCcw, Search, SlidersHorizontal } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { expenseCategories } from "@/lib/expense-format";
import type {
  ExpenseCategory,
  ExpenseSortBy,
  SortOrder,
} from "@/types/expense";
import type { FormEvent } from "react";

const selectClass =
  "h-10 w-full rounded-lg border border-slate-200 bg-white px-3 text-sm font-medium text-slate-700 outline-none transition focus:border-violet-400 focus:ring-4 focus:ring-violet-500/15";

type ExpenseFiltersProps = {
  category: ExpenseCategory | "";
  dateFrom: string;
  dateTo: string;
  filterError: string;
  isLoading: boolean;
  searchTerm: string;
  sortBy: ExpenseSortBy | "";
  sortOrder: SortOrder;
  onApply: (event: FormEvent<HTMLFormElement>) => void;
  onCategoryChange: (category: ExpenseCategory | "") => void;
  onDateFromChange: (date: string) => void;
  onDateToChange: (date: string) => void;
  onReset: () => void;
  onSearchChange: (searchTerm: string) => void;
  onSortByChange: (sortBy: ExpenseSortBy | "") => void;
  onSortOrderChange: (sortOrder: SortOrder) => void;
};

export function ExpenseFilters({
  category,
  dateFrom,
  dateTo,
  filterError,
  isLoading,
  searchTerm,
  sortBy,
  sortOrder,
  onApply,
  onCategoryChange,
  onDateFromChange,
  onDateToChange,
  onReset,
  onSearchChange,
  onSortByChange,
  onSortOrderChange,
}: ExpenseFiltersProps) {
  return (
    <Card className="p-4 sm:p-5">
      <form className="space-y-4" onSubmit={onApply}>
        <div className="grid gap-3 lg:grid-cols-[1.2fr_repeat(5,minmax(0,0.8fr))]">
          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Search current page
            </span>
            <div className="relative">
              <Search
                className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400"
                aria-hidden="true"
              />
              <Input
                className="pl-9"
                placeholder="Title, note, or category"
                value={searchTerm}
                onChange={(event) => onSearchChange(event.target.value)}
              />
            </div>
          </label>

          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Date from
            </span>
            <Input
              type="date"
              value={dateFrom}
              onChange={(event) => onDateFromChange(event.target.value)}
            />
          </label>

          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Date to
            </span>
            <Input
              type="date"
              value={dateTo}
              onChange={(event) => onDateToChange(event.target.value)}
            />
          </label>

          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Category
            </span>
            <select
              className={selectClass}
              value={category}
              onChange={(event) =>
                onCategoryChange(event.target.value as ExpenseCategory | "")
              }
            >
              <option value="">All categories</option>
              {expenseCategories.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </select>
          </label>

          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Sort by
            </span>
            <select
              className={selectClass}
              value={sortBy}
              onChange={(event) =>
                onSortByChange(event.target.value as ExpenseSortBy | "")
              }
            >
              <option value="">No sort</option>
              <option value="expense_date">Expense date</option>
              <option value="amount">Amount</option>
            </select>
          </label>

          <label className="block">
            <span className="mb-2 block text-xs font-bold uppercase text-slate-500">
              Order
            </span>
            <select
              className={selectClass}
              value={sortOrder}
              onChange={(event) => onSortOrderChange(event.target.value as SortOrder)}
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </label>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <p className="min-h-5 text-sm font-medium text-rose-600">
            {filterError}
          </p>
          <div className="flex flex-col gap-2 sm:flex-row">
            <Button type="button" variant="secondary" onClick={onReset} disabled={isLoading}>
              <RotateCcw className="h-4 w-4" aria-hidden="true" />
              Reset
            </Button>
            <Button type="submit" disabled={isLoading}>
              <SlidersHorizontal className="h-4 w-4" aria-hidden="true" />
              Apply Filters
            </Button>
          </div>
        </div>
      </form>
    </Card>
  );
}
