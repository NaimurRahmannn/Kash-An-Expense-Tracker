"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  BarChart3,
  Mic,
  PlusCircle,
  ReceiptText,
  Settings,
  WalletCards,
} from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import { cn } from "@/lib/utils";

const navItems = [
  {
    label: "Dashboard",
    href: "/dashboard",
    icon: BarChart3,
  },
  {
    label: "Expenses",
    href: "/expenses",
    icon: ReceiptText,
  },
  {
    label: "Add Expense",
    href: "/expenses/new",
    icon: PlusCircle,
  },
  {
    label: "Voice Input",
    href: "/voice",
    icon: Mic,
    badge: "New",
  },
  {
    label: "Settings",
    href: "/settings",
    icon: Settings,
  },
];

function isActivePath(pathname: string, href: string) {
  if (href === "/expenses") {
    return (
      pathname === "/expenses" ||
      (pathname.startsWith("/expenses/") && !pathname.startsWith("/expenses/new"))
    );
  }

  return pathname === href;
}

export function Sidebar() {
  const pathname = usePathname();

  return (
    <>
      <aside className="fixed inset-y-0 left-0 z-40 hidden w-72 flex-col border-r border-slate-200 bg-white px-5 py-6 shadow-sm shadow-slate-200/60 lg:flex">
        <Link href="/dashboard" className="flex items-center gap-3">
          <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 via-indigo-600 to-fuchsia-500 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-5 w-5" aria-hidden="true" />
          </div>
          <div>
            <p className="text-base font-bold text-slate-950">Expense Tracker</p>
            <p className="text-xs font-medium text-slate-500">
              Track smart. Spend better.
            </p>
          </div>
        </Link>

        <nav className="mt-9 space-y-1.5" aria-label="Main navigation">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActivePath(pathname, item.href);

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-semibold transition",
                  active
                    ? "bg-violet-50 text-violet-700 shadow-sm"
                    : "text-slate-600 hover:bg-slate-50 hover:text-slate-950",
                )}
              >
                <span className="flex items-center gap-3">
                  <Icon className="h-4 w-4" aria-hidden="true" />
                  {item.label}
                </span>
                {item.badge ? <Badge>{item.badge}</Badge> : null}
              </Link>
            );
          })}
        </nav>

        <Link
          href="/profile"
          className={cn(
            "mt-auto flex items-center gap-3 rounded-lg border border-slate-200 bg-slate-50 p-3 transition hover:border-violet-200 hover:bg-violet-50",
            pathname === "/profile" ? "border-violet-200 bg-violet-50" : "",
          )}
        >
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-violet-600 text-sm font-bold text-white">
            JD
          </div>
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold text-slate-950">John Doe</p>
            <p className="truncate text-xs text-slate-500">john.doe@example.com</p>
          </div>
        </Link>
      </aside>

      <nav
        className="fixed inset-x-3 bottom-3 z-50 grid grid-cols-5 gap-1 rounded-lg border border-slate-200 bg-white/95 p-1 shadow-lg shadow-slate-200/80 backdrop-blur lg:hidden"
        aria-label="Mobile navigation"
      >
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = isActivePath(pathname, item.href);

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "relative flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg px-1 text-[11px] font-semibold transition",
                active ? "bg-violet-50 text-violet-700" : "text-slate-500",
              )}
            >
              <Icon className="h-4 w-4" aria-hidden="true" />
              <span className="max-w-full truncate">{item.label}</span>
              {item.badge ? (
                <span className="absolute right-1 top-1 h-1.5 w-1.5 rounded-full bg-violet-500" />
              ) : null}
            </Link>
          );
        })}
      </nav>
    </>
  );
}
