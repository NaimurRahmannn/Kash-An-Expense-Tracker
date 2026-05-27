import { Loader2 } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { formatDisplayDate } from "@/lib/date";
import type { DateRange } from "@/types/dashboard";
import type { ExpenseSummary } from "@/types/summary";

type SummaryRangeCardProps = {
  appliedRange: DateRange;
  isSummaryLoading: boolean;
  summary: ExpenseSummary | null;
};

export function SummaryRangeCard({
  appliedRange,
  isSummaryLoading,
  summary,
}: SummaryRangeCardProps) {
  return (
    <Card className="rounded-[14px] border-violet-200/80 bg-linear-to-r from-white via-violet-50/50 to-white p-5 shadow-[0_8px_24px_rgba(109,53,242,0.10)]">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-lg font-bold text-slate-950">Selected range summary</h2>
          <p className="mt-1 text-sm text-slate-500">
            {formatDisplayDate(summary?.date_from ?? appliedRange.dateFrom)} -{" "}
            {formatDisplayDate(summary?.date_to ?? appliedRange.dateTo)}
          </p>
        </div>
        {isSummaryLoading ? (
          <div className="inline-flex items-center gap-2 text-sm font-semibold text-violet-700">
            <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
            Refreshing summary...
          </div>
        ) : null}
      </div>
    </Card>
  );
}
