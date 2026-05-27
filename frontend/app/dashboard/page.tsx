"use client";

import Link from "next/link";
import { AppShell } from "@/components/layout/AppShell";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { CategoryBreakdown } from "@/components/dashboard/CategoryBreakdown";
import { DashboardError, DashboardLoading } from "@/components/dashboard/DashboardState";
import { DateRangeCard } from "@/components/dashboard/DateRangeCard";
import { QuickActions } from "@/components/dashboard/QuickActions";
import { RecentExpenses } from "@/components/dashboard/RecentExpenses";
import { StatCard } from "@/components/dashboard/StatCard";
import { SummaryRangeCard } from "@/components/dashboard/SummaryRangeCard";
import { useDashboardData } from "@/hooks/useDashboardData";

export default function DashboardPage() {
  const dashboard = useDashboardData();

  return (
    <AppShell title="Dashboard">
      <div className="space-y-5 sm:space-y-6">
        <DateRangeCard
          dateFrom={dashboard.dateRange.dateFrom}
          dateTo={dashboard.dateRange.dateTo}
          isLoading={dashboard.isSummaryLoading}
          validationError={dashboard.dateError}
          onApply={dashboard.handleApplyDateRange}
          onDateFromChange={dashboard.setDateFrom}
          onDateToChange={dashboard.setDateTo}
        />

        {dashboard.isDashboardLoading ? (
          <DashboardLoading />
        ) : dashboard.error ? (
          <DashboardError
            error={dashboard.error}
            onRetry={() => void dashboard.loadDashboard()}
          />
        ) : (
          <>
            <SummaryRangeCard
              appliedRange={dashboard.appliedRange}
              isSummaryLoading={dashboard.isSummaryLoading}
              summary={dashboard.summary}
            />

            <div className="grid grid-cols-2 gap-3 sm:gap-6 xl:grid-cols-4">
              {dashboard.statCards.map((item) => (
                <StatCard key={item.title} {...item} />
              ))}
            </div>

            <div className="grid gap-6 xl:grid-cols-[1fr_0.7fr]">
              <Card className="rounded-[14px] p-5 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
                <div className="mb-5 flex items-center justify-between gap-3">
                  <div>
                    <h2 className="text-lg font-bold text-slate-950">
                      Expenses by Category
                    </h2>
                    <p className="mt-1 text-sm text-slate-500">
                      Real breakdown from the summary API
                    </p>
                  </div>
                  <Badge>{dashboard.summary?.by_category.length ?? 0} categories</Badge>
                </div>
                <CategoryBreakdown
                  categories={dashboard.summary?.by_category ?? []}
                  totalAmount={dashboard.totalAmount}
                />
              </Card>

              <QuickActions />
            </div>

            <Card className="rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
              <div className="flex items-center justify-between px-2 pb-3">
                <div>
                  <h2 className="text-lg font-bold text-slate-950">Recent Expenses</h2>
                  <p className="mt-1 text-sm text-slate-500">
                    Latest entries from your expense list
                  </p>
                </div>
                <Link
                  href="/expenses"
                  className="text-sm font-semibold text-violet-700 hover:text-violet-800"
                >
                  View all
                </Link>
              </div>
              <RecentExpenses expenses={dashboard.recentExpenses} />
            </Card>
          </>
        )}
      </div>
    </AppShell>
  );
}
