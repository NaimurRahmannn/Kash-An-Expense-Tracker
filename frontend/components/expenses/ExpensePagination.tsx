import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/Button";

type ExpensePaginationProps = {
  canGoNext: boolean;
  canGoPrevious: boolean;
  isLoading: boolean;
  itemCount: number;
  limit: number;
  page: number;
  onNext: () => void;
  onPrevious: () => void;
};

export function ExpensePagination({
  canGoNext,
  canGoPrevious,
  isLoading,
  itemCount,
  limit,
  page,
  onNext,
  onPrevious,
}: ExpensePaginationProps) {
  return (
    <div className="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm text-slate-500">
        Showing {itemCount} of up to {limit} fetched expenses
      </p>
      <div className="flex items-center gap-2">
        <Button
          variant="secondary"
          onClick={onPrevious}
          disabled={!canGoPrevious || isLoading}
          aria-label="Previous page"
        >
          <ChevronLeft className="h-4 w-4" aria-hidden="true" />
        </Button>
        <span className="rounded-lg bg-violet-50 px-3 py-2 text-sm font-bold text-violet-700">
          Page {page}
        </span>
        <Button
          variant="secondary"
          onClick={onNext}
          disabled={!canGoNext || isLoading}
          aria-label="Next page"
        >
          <ChevronRight className="h-4 w-4" aria-hidden="true" />
        </Button>
      </div>
    </div>
  );
}
