import { Card } from "@/components/ui/Card";
import { cn } from "@/lib/utils";
import type { DashboardStat, StatAccent } from "@/types/dashboard";

const statAccentStyles: Record<StatAccent, string> = {
  violet: "bg-violet-100 text-violet-700",
  blue: "bg-blue-100 text-blue-600",
  green: "bg-emerald-100 text-emerald-600",
  amber: "bg-amber-100 text-amber-600",
};

export function StatCard({ accent, detail, icon: Icon, title, value }: DashboardStat) {
  return (
    <Card className="flex min-h-28 items-center gap-3 rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)] sm:p-5">
      <div
        className={cn(
          "flex h-13 w-13 shrink-0 items-center justify-center rounded-full sm:h-16 sm:w-16",
          statAccentStyles[accent],
        )}
      >
        <Icon className="h-7 w-7 sm:h-8 sm:w-8" aria-hidden="true" />
      </div>
      <div className="min-w-0">
        <p className="truncate text-xs font-semibold text-slate-950 sm:text-sm">
          {title}
        </p>
        <p className="mt-1 truncate text-xl font-bold tracking-[-0.01em] text-slate-950 sm:text-2xl">
          {value}
        </p>
        <p className="mt-1 truncate text-xs text-slate-500 sm:text-sm">{detail}</p>
      </div>
    </Card>
  );
}
