import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";

export function ExpensesLoading() {
  return (
    <Card className="p-6">
      <div className="flex items-center gap-3">
        <Loader2 className="h-5 w-5 animate-spin text-violet-600" aria-hidden="true" />
        <p className="text-sm font-semibold text-slate-700">Loading expenses...</p>
      </div>
      <div className="mt-5 space-y-3">
        {["one", "two", "three", "four"].map((item) => (
          <div
            key={item}
            className="h-12 animate-pulse rounded-lg bg-slate-100"
          />
        ))}
      </div>
    </Card>
  );
}

export function ExpensesError({
  error,
  onRetry,
}: {
  error: string;
  onRetry: () => void;
}) {
  return (
    <Card className="border-rose-200 bg-rose-50 p-6">
      <h2 className="text-lg font-bold text-rose-900">Unable to load expenses</h2>
      <p className="mt-2 text-sm leading-6 text-rose-700">{error}</p>
      <Button className="mt-5" onClick={onRetry}>
        Retry
      </Button>
    </Card>
  );
}
