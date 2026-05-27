import Link from "next/link";
import { ArrowRight, CirclePlus, FileText, Mic } from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";

export function QuickActions() {
  return (
    <Card className="rounded-[14px] p-5 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
      <h2 className="text-lg font-bold text-slate-950">Quick Actions</h2>
      <div className="mt-4 grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
        <Link
          href="/expenses/new"
          className="flex h-16 items-center justify-between rounded-[14px] border border-violet-200 bg-violet-50 px-4 text-slate-950 transition hover:border-violet-300 hover:bg-violet-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-violet-600 text-white shadow-md shadow-violet-200">
              <CirclePlus className="h-6 w-6" aria-hidden="true" />
            </span>
            Add Expense
          </span>
          <ArrowRight className="h-5 w-5 text-slate-950" aria-hidden="true" />
        </Link>

        <Link
          href="/expenses"
          className="flex h-16 items-center justify-between rounded-[14px] border border-blue-200 bg-blue-50 px-4 text-slate-950 transition hover:border-blue-300 hover:bg-blue-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-100 text-blue-600">
              <FileText className="h-6 w-6" aria-hidden="true" />
            </span>
            View All Expenses
          </span>
          <ArrowRight className="h-5 w-5 text-slate-950" aria-hidden="true" />
        </Link>

        <Link
          href="/voice"
          className="flex h-16 items-center justify-between rounded-[14px] border border-emerald-200 bg-emerald-50 px-4 text-slate-950 transition hover:border-emerald-300 hover:bg-emerald-100"
        >
          <span className="flex items-center gap-3 text-sm font-bold">
            <span className="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
              <Mic className="h-6 w-6" aria-hidden="true" />
            </span>
            Voice Input
          </span>
          <Badge className="bg-violet-100 text-violet-700 ring-0">NEW</Badge>
        </Link>
      </div>
    </Card>
  );
}
