import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";

export function DashboardLoading() {
  return (
    <div className="space-y-6">
      <Card className="flex min-h-28 items-center gap-4 rounded-[14px] p-5">
        <Loader2 className="h-5 w-5 animate-spin text-violet-600" aria-hidden="true" />
        <p className="text-sm font-semibold text-slate-700">Loading dashboard...</p>
      </Card>
      <div className="grid grid-cols-2 gap-3 sm:gap-6 xl:grid-cols-4">
        {["one", "two", "three", "four"].map((item) => (
          <Card key={item} className="h-28 animate-pulse rounded-[14px] bg-slate-100" />
        ))}
      </div>
      <div className="grid gap-6 xl:grid-cols-[1fr_0.7fr]">
        <Card className="h-80 animate-pulse rounded-[14px] bg-slate-100" />
        <Card className="h-80 animate-pulse rounded-[14px] bg-slate-100" />
      </div>
    </div>
  );
}

export function DashboardError({
  error,
  onRetry,
}: {
  error: string;
  onRetry: () => void;
}) {
  return (
    <Card className="rounded-[14px] border-rose-200 bg-rose-50 p-6">
      <h2 className="text-lg font-bold text-rose-900">
        Unable to load dashboard data
      </h2>
      <p className="mt-2 text-sm leading-6 text-rose-700">{error}</p>
      <Button className="mt-5" onClick={onRetry}>
        Retry
      </Button>
    </Card>
  );
}
