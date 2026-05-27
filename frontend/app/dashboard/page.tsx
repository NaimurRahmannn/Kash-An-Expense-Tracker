import Link from "next/link";
import {
  ArrowRight,
  CirclePlus,
  FileText,
  Mic,
  ReceiptText,
  ShoppingBag,
  Star,
  TrendingUp,
  Utensils,
  WalletCards,
  Zap,
} from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { cn } from "@/lib/utils";

const accentStyles = {
  violet: {
    icon: "bg-violet-100 text-violet-700",
    line: "#6d35f2",
  },
  blue: {
    icon: "bg-blue-100 text-blue-600",
    line: "#2c8cff",
  },
  green: {
    icon: "bg-emerald-100 text-emerald-600",
    line: "#35b978",
  },
  amber: {
    icon: "bg-amber-100 text-amber-600",
    line: "#f59e0b",
  },
} as const;

const stats = [
  {
    title: "Total Expenses",
    value: "18,560",
    prefix: true,
    detail: "In selected range",
    icon: WalletCards,
    accent: "violet" as const,
    line: "M2 19 L12 10 L20 13 L30 4 L38 7 L47 1",
  },
  {
    title: "Total Transactions",
    value: "42",
    detail: "In selected range",
    icon: ReceiptText,
    accent: "blue" as const,
    line: "M2 17 L11 8 L20 11 L31 3 L45 13",
  },
  {
    title: "Daily Average",
    value: "2,651",
    prefix: true,
    detail: "In selected range",
    icon: TrendingUp,
    accent: "green" as const,
    line: "M2 18 L13 8 L24 12 L34 3 L45 9",
  },
  {
    title: "Highest Category",
    value: "Food",
    detail: "6,450 (34.7%)",
    prefixDetail: true,
    icon: Star,
    accent: "amber" as const,
    line: "M2 18 L12 9 L22 12 L33 3 L45 13",
  },
];

const categoryRows = [
  { label: "Food", amount: "6,450", percent: "34.7%", color: "bg-violet-500" },
  { label: "Transport", amount: "3,120", percent: "16.8%", color: "bg-blue-500" },
  { label: "Housing", amount: "2,800", percent: "15.1%", color: "bg-emerald-400" },
  { label: "Entertainment", amount: "1,850", percent: "10.0%", color: "bg-yellow-400" },
  { label: "Shopping", amount: "1,600", percent: "8.6%", color: "bg-amber-400" },
  { label: "Healthcare", amount: "1,050", percent: "5.7%", color: "bg-orange-400" },
  { label: "Education", amount: "900", percent: "4.9%", color: "bg-pink-400" },
  { label: "Utilities", amount: "550", percent: "3.0%", color: "bg-fuchsia-400" },
  { label: "Other", amount: "240", percent: "1.3%", color: "bg-slate-400" },
];

const recentRows = [
  {
    date: "May 18, 2025",
    time: "12:45 PM",
    title: "Pizza Hut",
    category: "Food",
    amount: "560",
    color: "violet",
    icon: Utensils,
    avatar: "PH",
    avatarClass: "bg-white text-red-600 ring-1 ring-slate-200",
  },
  {
    date: "May 18, 2025",
    time: "9:15 AM",
    title: "Uber Ride",
    category: "Transport",
    amount: "320",
    color: "blue",
    icon: ReceiptText,
    avatar: "Uber",
    avatarClass: "bg-black text-white",
  },
  {
    date: "May 17, 2025",
    time: "6:30 PM",
    title: "BigBasket Groceries",
    category: "Food",
    amount: "1,250",
    color: "violet",
    icon: Utensils,
    avatar: "bb",
    avatarClass: "bg-lime-100 text-lime-700 ring-1 ring-lime-200",
  },
  {
    date: "May 16, 2025",
    time: "8:20 PM",
    title: "Movie Tickets",
    category: "Entertainment",
    amount: "480",
    color: "amber",
    icon: ShoppingBag,
    avatar: "",
    avatarClass: "bg-amber-100 text-amber-600",
  },
  {
    date: "May 16, 2025",
    time: "11:05 AM",
    title: "Electricity Bill",
    category: "Utilities",
    amount: "550",
    color: "emerald",
    icon: Zap,
    avatar: "",
    avatarClass: "bg-emerald-100 text-emerald-600",
  },
];

function Waveform() {
  const bars = [
    3, 4, 5, 8, 18, 26, 34, 44, 32, 22, 12, 28, 16, 10, 6, 8, 24, 32, 20, 12, 18,
    42, 60, 36, 24, 12, 16, 38, 30, 18, 8, 16, 34, 48, 28, 14, 18, 26, 20, 16, 10,
    12, 18, 9, 6,
  ];

  return (
    <div className="flex h-14 min-w-0 flex-1 items-center justify-center gap-0.5 overflow-hidden px-1 text-violet-600 sm:h-20 sm:gap-1 sm:px-4">
      {bars.map((bar, index) => (
        <span
          key={`${bar}-${index}`}
          className="w-0.5 rounded-full bg-violet-600"
          style={{ height: `${bar}px` }}
        />
      ))}
    </div>
  );
}

function Sparkline({ path, accent }: { path: string; accent: string }) {
  return (
    <svg className="h-7 w-12 sm:h-8 sm:w-14" viewBox="0 0 49 22" fill="none" aria-hidden="true">
      <path
        d={path}
        stroke={accent}
        strokeWidth="2.4"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function StatCard({ item }: { item: (typeof stats)[number] }) {
  const Icon = item.icon;
  const styles = accentStyles[item.accent];

  return (
    <Card className="relative flex min-h-28 items-center gap-2 overflow-hidden rounded-[14px] p-3 shadow-[0_6px_18px_rgba(15,23,42,0.08)] sm:min-h-29.5 sm:justify-between sm:gap-4 sm:p-5">
      <div className="flex min-w-0 items-center gap-3 sm:gap-4">
        <div
          className={cn(
            "flex h-13 w-13 shrink-0 items-center justify-center rounded-full sm:h-16 sm:w-16",
            styles.icon,
          )}
        >
          <Icon className="h-7 w-7 sm:h-8 sm:w-8" aria-hidden="true" />
        </div>
        <div className="min-w-0">
          <p className="truncate text-xs font-semibold text-slate-950 sm:text-sm">
            {item.title}
          </p>
          <p className="mt-1 text-xl font-bold tracking-[-0.01em] text-slate-950 sm:text-2xl">
            {item.prefix ? <span>&#8377;</span> : null}
            {item.value}
          </p>
          <p className="mt-1 truncate text-xs text-slate-500 sm:text-sm">
            {item.prefixDetail ? <span>&#8377;</span> : null}
            {item.detail}
          </p>
        </div>
      </div>
      <div className="hidden sm:static sm:block">
        <Sparkline path={item.line} accent={styles.line} />
      </div>
    </Card>
  );
}

function ExpenseLineChart() {
  const labels = ["May 12", "May 13", "May 14", "May 15", "May 16", "May 17", "May 18"];
  const points = "18,112 108,84 198,122 288,42 378,126 468,92 558,134";

  return (
    <div className="mt-4 h-52 sm:mt-5 sm:h-57.5">
      <svg className="h-full w-full" viewBox="0 0 590 230" fill="none" aria-hidden="true">
        <defs>
          <linearGradient id="expenseArea" x1="0" x2="0" y1="40" y2="172">
            <stop offset="0%" stopColor="#7c3aed" stopOpacity="0.2" />
            <stop offset="100%" stopColor="#7c3aed" stopOpacity="0" />
          </linearGradient>
        </defs>
        {[172, 137, 102, 67, 32].map((y) => (
          <line key={y} x1="18" x2="578" y1={y} y2={y} stroke="#e5e7eb" />
        ))}
        {["0", "1K", "2K", "3K", "4K", "5K"].map((label, index) => (
          <text
            key={label}
            x="0"
            y={177 - index * 28}
            className="fill-slate-500 text-[12px]"
          >
            {label}
          </text>
        ))}
        <path d={`M18 112 L108 84 L198 122 L288 42 L378 126 L468 92 L558 134 L558 172 L18 172 Z`} fill="url(#expenseArea)" />
        <polyline
          points={points}
          stroke="#6d35f2"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2.4"
        />
        {points.split(" ").map((point) => {
          const [cx, cy] = point.split(",");

          return (
            <circle
              key={point}
              cx={cx}
              cy={cy}
              r="4"
              fill="white"
              stroke="#6d35f2"
              strokeWidth="2.5"
            />
          );
        })}
        {labels.map((label, index) => (
          <text
            key={label}
            x={18 + index * 90}
            y="212"
            textAnchor="middle"
            className="fill-slate-600 text-[12px]"
          >
            {label}
          </text>
        ))}
      </svg>
    </div>
  );
}

function CategoryDonut() {
  return (
    <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:gap-6">
      <div className="relative mx-auto h-32 w-32 shrink-0 rounded-full bg-[conic-gradient(#8b5cf6_0deg_125deg,#4f8ff7_125deg_186deg,#55c7a5_186deg_240deg,#f4cf59_240deg_276deg,#f8b94a_276deg_307deg,#fb8a53_307deg_328deg,#f471a5_328deg_346deg,#e879f9_346deg_357deg,#9ca3af_357deg_360deg)] sm:h-52 sm:w-52">
        <div className="absolute inset-8 flex flex-col items-center justify-center rounded-full bg-white shadow-inner sm:inset-12">
          <p className="text-sm font-bold text-slate-950 sm:text-xl">&#8377;18,560</p>
          <p className="text-xs text-slate-500 sm:text-sm">Total</p>
        </div>
      </div>

      <div className="min-w-0 flex-1 space-y-1 sm:space-y-2">
        {categoryRows.map((row) => (
          <div
            key={row.label}
            className="grid grid-cols-[1fr_auto_auto] items-center gap-2 text-xs sm:gap-4 sm:text-sm"
          >
            <div className="flex min-w-0 items-center gap-2">
              <span className={cn("h-2.5 w-2.5 rounded-full sm:h-3 sm:w-3", row.color)} />
              <span className="truncate text-slate-700">{row.label}</span>
            </div>
            <span className="font-medium text-slate-950">&#8377;{row.amount}</span>
            <span className="w-10 text-right text-slate-700 sm:w-12">{row.percent}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

function CategoryBadge({
  color,
  children,
  icon: Icon,
}: {
  color: string;
  children: string;
  icon: typeof Utensils;
}) {
  const colorClass = {
    violet: "bg-violet-100 text-violet-700",
    blue: "bg-blue-100 text-blue-700",
    amber: "bg-amber-100 text-amber-700",
    emerald: "bg-emerald-100 text-emerald-700",
  }[color];

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs font-semibold",
        colorClass,
      )}
    >
      <Icon className="h-3.5 w-3.5" aria-hidden="true" />
      {children}
    </span>
  );
}

export default function DashboardPage() {
  return (
    <AppShell title="Dashboard">
      <div className="space-y-5 sm:space-y-6">
        <Card className="overflow-hidden rounded-[14px] border-violet-200/80 bg-linear-to-r from-white via-violet-50/40 to-white p-3 shadow-[0_8px_24px_rgba(109,53,242,0.10)] sm:p-5">
          <div className="grid min-h-32 grid-cols-[64px_minmax(0,1fr)_112px] grid-rows-[auto_1fr] items-center gap-x-2 gap-y-1 sm:min-h-34 sm:grid-cols-[112px_minmax(0,1fr)_230px] sm:gap-x-5 lg:grid-cols-[130px_minmax(0,1fr)_260px]">
            <div className="row-span-2 flex items-center justify-center">
              <div className="flex h-15 w-15 shrink-0 items-center justify-center rounded-full bg-violet-100 ring-7 ring-violet-100/60 sm:h-29 sm:w-29 sm:ring-14">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-linear-to-b from-violet-500 to-violet-800 text-white shadow-xl shadow-violet-300 sm:h-21.5 sm:w-21.5">
                  <Mic className="h-7 w-7 sm:h-11 sm:w-11" aria-hidden="true" />
                </div>
              </div>
            </div>

            <div className="col-start-2 row-start-1 min-w-0 self-end">
              <h2 className="truncate text-[13px] font-bold text-slate-950 min-[390px]:text-sm sm:text-lg">
                Add expense with your voice
              </h2>
              <p className="mt-1 truncate text-[11px] font-medium text-slate-600 min-[390px]:text-xs sm:mt-2 sm:text-base">
                Tap the mic and speak your expense
              </p>
            </div>

            <div className="col-start-2 row-start-2 min-w-0 self-start">
              <Waveform />
            </div>

            <div className="col-start-3 row-span-2 row-start-1 min-w-0 self-center">
              <Badge className="gap-1.5 bg-white px-2.5 py-1 text-[11px] text-slate-900 shadow-sm ring-0 sm:gap-2 sm:px-3 sm:text-sm">
                <span className="h-2 w-2 rounded-full bg-lime-500 sm:h-2.5 sm:w-2.5" />
                Listening...
              </Badge>
              <p className="mt-3 text-[11px] font-semibold text-slate-950 sm:text-sm">
                Live transcript
              </p>
              <p className="mt-1 text-[11px] italic leading-4 text-violet-700 min-[390px]:text-xs sm:text-sm sm:leading-6">
                &quot;Had lunch at Pizza Hut for 560 rupees&quot;
              </p>
            </div>
          </div>
        </Card>

        <div className="grid grid-cols-2 gap-3 sm:gap-6 xl:grid-cols-4">
          {stats.map((item) => (
            <StatCard key={item.title} item={item} />
          ))}
        </div>

        <div className="grid grid-cols-2 gap-4 sm:gap-6 xl:grid-cols-[1fr_1.035fr]">
          <Card className="min-w-0 rounded-[14px] p-3 shadow-[0_6px_18px_rgba(15,23,42,0.08)] sm:p-5">
            <div className="flex items-center justify-between">
              <h2 className="text-sm font-bold text-slate-950 sm:text-lg">Expense Overview</h2>
              <div className="flex items-center gap-2 text-xs font-medium text-slate-700 sm:text-sm">
                <span className="h-0.5 w-4 rounded-full bg-violet-600" />
                Amount (&#8377;)
              </div>
            </div>
            <ExpenseLineChart />
          </Card>

          <Card className="min-w-0 rounded-[14px] p-3 shadow-[0_6px_18px_rgba(15,23,42,0.08)] sm:p-5">
            <h2 className="text-sm font-bold text-slate-950 sm:text-lg">
              Expenses by Category
            </h2>
            <div className="mt-4">
              <CategoryDonut />
            </div>
          </Card>
        </div>

        <div className="grid gap-5 sm:gap-6 xl:grid-cols-[1fr_0.6fr]">
          <Card className="rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
            <div className="flex items-center justify-between px-2 pb-3">
              <h2 className="text-lg font-bold text-slate-950">Recent Expenses</h2>
              <Link
                href="/expenses"
                className="text-sm font-semibold text-violet-700 hover:text-violet-800 lg:hidden"
              >
                View all
              </Link>
            </div>

            <div className="space-y-0 lg:hidden">
              {recentRows.map((row) => (
                <Link
                  key={`${row.date}-${row.title}-mobile`}
                  href="/expenses"
                  className="grid grid-cols-[48px_minmax(0,1fr)_auto_auto] items-center gap-3 border-b border-slate-200 px-2 py-3 last:border-b-0"
                >
                  <span
                    className={cn(
                      "flex h-10 w-10 items-center justify-center rounded-full text-xs font-black",
                      row.avatarClass,
                    )}
                  >
                    {row.avatar ? row.avatar : <row.icon className="h-5 w-5" aria-hidden="true" />}
                  </span>
                  <span className="min-w-0">
                    <span className="block truncate text-base font-bold text-slate-950">
                      {row.title}
                    </span>
                    <span className="block truncate text-sm font-medium text-slate-500">
                      {row.date} • {row.time}
                    </span>
                  </span>
                  <CategoryBadge color={row.color} icon={row.icon}>
                    {row.category}
                  </CategoryBadge>
                  <span className="flex items-center gap-3 text-base font-bold text-slate-950">
                    &#8377;{row.amount}
                    <ArrowRight className="h-4 w-4" aria-hidden="true" />
                  </span>
                </Link>
              ))}
            </div>

            <div className="hidden overflow-hidden rounded-lg border border-slate-200 lg:block">
              <div className="overflow-x-auto">
                <table className="min-w-full text-sm">
                  <thead className="bg-slate-50">
                    <tr className="border-b border-slate-200 text-left text-slate-700">
                      <th className="px-3 py-2 font-semibold">Date</th>
                      <th className="px-3 py-2 font-semibold">Title</th>
                      <th className="px-3 py-2 font-semibold">Category</th>
                      <th className="px-3 py-2 font-semibold">Amount</th>
                    </tr>
                  </thead>
                  <tbody>
                    {recentRows.map((row) => (
                      <tr key={`${row.date}-${row.title}`} className="border-b border-slate-200 last:border-b-0">
                        <td className="whitespace-nowrap px-3 py-2.5 text-slate-800">{row.date}</td>
                        <td className="whitespace-nowrap px-3 py-2.5 font-medium text-slate-900">
                          {row.title}
                        </td>
                        <td className="whitespace-nowrap px-3 py-2.5">
                          <CategoryBadge color={row.color} icon={row.icon}>
                            {row.category}
                          </CategoryBadge>
                        </td>
                        <td className="whitespace-nowrap px-3 py-2.5 font-semibold text-slate-950">
                          &#8377;{row.amount}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
            <div className="pt-3 text-center">
              <Link href="/expenses" className="hidden text-sm font-semibold text-violet-700 hover:text-violet-800 lg:inline">
                View all expenses
              </Link>
            </div>
          </Card>

          <Card className="hidden rounded-[14px] p-6 shadow-[0_6px_18px_rgba(15,23,42,0.08)] lg:block">
            <h2 className="text-lg font-bold text-slate-950">Quick Actions</h2>
            <div className="mt-4 space-y-4">
              <Link
                href="/expenses/new"
                className="flex h-16 items-center justify-between rounded-[14px] border border-violet-200 bg-violet-50 px-4 text-slate-950 transition hover:border-violet-300 hover:bg-violet-100"
              >
                <span className="flex items-center gap-4 text-sm font-bold">
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
                <span className="flex items-center gap-4 text-sm font-bold">
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
                <span className="flex items-center gap-4 text-sm font-bold">
                  <span className="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
                    <Mic className="h-6 w-6" aria-hidden="true" />
                  </span>
                  Voice Input
                </span>
                <span className="flex items-center gap-4">
                  <Badge className="bg-violet-100 text-violet-700 ring-0">NEW</Badge>
                  <ArrowRight className="h-5 w-5 text-slate-950" aria-hidden="true" />
                </span>
              </Link>
            </div>
          </Card>
        </div>

        <Card className="overflow-hidden rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)] lg:hidden">
          <h2 className="mb-4 text-xl font-bold text-slate-950">Quick Actions</h2>
          <div className="grid min-w-0 grid-cols-3 gap-2 min-[390px]:gap-3">
            <Link
              href="/expenses/new"
              className="flex h-24 min-w-0 flex-col items-center justify-center gap-2 overflow-hidden rounded-[14px] border border-violet-200 bg-violet-50 px-2 text-[11px] font-bold leading-tight text-slate-950 min-[390px]:text-xs"
            >
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-violet-600 text-white shadow-md shadow-violet-200 min-[390px]:h-11 min-[390px]:w-11">
                <CirclePlus className="h-6 w-6 min-[390px]:h-7 min-[390px]:w-7" aria-hidden="true" />
              </span>
              <span className="block w-full max-w-18 whitespace-normal text-center">
                Add Expense
              </span>
            </Link>
            <Link
              href="/expenses"
              className="flex h-24 min-w-0 flex-col items-center justify-center gap-2 overflow-hidden rounded-[14px] border border-blue-200 bg-blue-50 px-2 text-[11px] font-bold leading-tight text-slate-950 min-[390px]:text-xs"
            >
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-blue-100 text-blue-600 min-[390px]:h-11 min-[390px]:w-11">
                <FileText className="h-6 w-6 min-[390px]:h-7 min-[390px]:w-7" aria-hidden="true" />
              </span>
              <span className="block w-full max-w-18 whitespace-normal text-center">
                View All Expenses
              </span>
            </Link>
            <Link
              href="/voice"
              className="relative flex h-24 min-w-0 flex-col items-center justify-center gap-2 overflow-hidden rounded-[14px] border border-emerald-200 bg-emerald-50 px-2 pt-4 text-[11px] font-bold leading-tight text-slate-950 min-[390px]:text-xs"
            >
              <Badge className="absolute right-1 top-1 bg-violet-100 px-2 py-0.5 text-[10px] text-violet-700 ring-0">
                NEW
              </Badge>
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 min-[390px]:h-11 min-[390px]:w-11">
                <Mic className="h-6 w-6 min-[390px]:h-7 min-[390px]:w-7" aria-hidden="true" />
              </span>
              <span className="block w-full max-w-18 whitespace-normal text-center">
                Voice Input
              </span>
r            </Link>
          </div>
        </Card>
      </div>
    </AppShell>
  );
}
