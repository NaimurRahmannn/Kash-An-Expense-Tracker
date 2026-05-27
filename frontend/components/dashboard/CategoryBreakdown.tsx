import { EmptyState } from "@/components/ui/EmptyState";
import {
  formatMoney,
  getCategoryPercent,
  sortCategoriesByTotal,
} from "@/lib/dashboard-format";
import { cn } from "@/lib/utils";
import type { CategorySummary } from "@/types/summary";

const categoryColors = [
  "bg-violet-500",
  "bg-blue-500",
  "bg-emerald-500",
  "bg-amber-500",
  "bg-orange-500",
  "bg-pink-500",
  "bg-cyan-500",
  "bg-fuchsia-500",
  "bg-slate-500",
];

type CategoryBreakdownProps = {
  categories: CategorySummary[];
  totalAmount: number;
};

export function CategoryBreakdown({
  categories,
  totalAmount,
}: CategoryBreakdownProps) {
  const sortedCategories = sortCategoriesByTotal(categories);

  if (sortedCategories.length === 0) {
    return (
      <EmptyState
        title="No category data"
        description="Add expenses to see category breakdown."
      />
    );
  }

  return (
    <div className="space-y-4">
      {sortedCategories.map((category, index) => {
        const percent = getCategoryPercent(category.total, totalAmount);

        return (
          <div key={category.category}>
            <div className="mb-2 grid grid-cols-[1fr_auto] gap-3 text-sm">
              <div className="min-w-0">
                <p className="truncate font-semibold text-slate-950">
                  {category.category}
                </p>
                <p className="text-xs text-slate-500">
                  {category.count} {category.count === 1 ? "transaction" : "transactions"}
                </p>
              </div>
              <div className="text-right">
                <p className="font-bold text-slate-950">{formatMoney(category.total)}</p>
                <p className="text-xs text-slate-500">{percent.toFixed(1)}%</p>
              </div>
            </div>
            <div className="h-2 rounded-full bg-slate-100">
              <div
                className={cn(
                  "h-2 rounded-full",
                  categoryColors[index % categoryColors.length],
                )}
                style={{ width: `${Math.min(100, percent)}%` }}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
}
