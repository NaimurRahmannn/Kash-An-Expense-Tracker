"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  ChevronsLeft,
  CirclePlus,
  LayoutGrid,
  LogOut,
  Mic,
  ReceiptText,
  Settings,
  User,
  WalletCards,
} from "lucide-react";
import { Badge } from "@/components/ui/Badge";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/utils";

const navItems = [
  {
    label: "Dashboard",
    href: "/dashboard",
    icon: LayoutGrid,
  },
  {
    label: "Expenses",
    href: "/expenses",
    icon: ReceiptText,
  },
  {
    label: "Add Expense",
    href: "/expenses/new",
    icon: CirclePlus,
  },
  {
    label: "Voice Input",
    href: "/voice",
    icon: Mic,
    badge: "NEW",
  },
  {
    label: "Profile",
    href: "/profile",
    icon: User,
  },
  {
    label: "Settings",
    href: "/settings",
    icon: Settings,
  },
];

const mobileNavItems = navItems.filter((item) => item.href !== "/settings");

function getInitials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

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
  const { isLoading, logout, user } = useAuth();
  const displayName = isLoading ? "User" : user?.name ?? "User";
  const displayEmail = isLoading ? "Loading..." : user?.email ?? "";
  const initials = getInitials(displayName) || "U";

  return (
    <>
      <aside className="fixed inset-y-0 left-0 z-40 hidden w-73 flex-col border-r border-slate-200/80 bg-white px-5 py-7 lg:flex">
        <Link href="/dashboard" className="flex items-center gap-4 px-0.5">
          <div className="flex h-12 w-12 items-center justify-center rounded-[14px] bg-linear-to-br from-violet-400 to-violet-700 text-white shadow-lg shadow-violet-200">
            <WalletCards className="h-6 w-6" aria-hidden="true" />
          </div>
          <p className="text-xl font-bold text-violet-700">Expense Tracker</p>
        </Link>

        <nav className="mt-9 space-y-3" aria-label="Main navigation">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActivePath(pathname, item.href);

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex h-14 items-center justify-between rounded-lg px-4 text-[15px] font-semibold transition",
                  active
                    ? "bg-violet-50 text-violet-700 shadow-sm shadow-violet-100"
                    : "text-slate-800 hover:bg-violet-50 hover:text-violet-700",
                )}
              >
                <span className="flex items-center gap-4">
                  <Icon
                    className={cn("h-5 w-5", active ? "text-violet-700" : "text-slate-600")}
                    aria-hidden="true"
                  />
                  {item.label}
                </span>
                {item.badge ? (
                  <Badge className="bg-violet-100 px-2 py-0.5 text-[11px] text-violet-700 ring-0">
                    {item.badge}
                  </Badge>
                ) : null}
              </Link>
            );
          })}
        </nav>

        <div className="mt-auto rounded-lg border border-slate-200 bg-slate-50 p-3">
          <Link href="/profile" className="flex min-w-0 items-center gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-violet-600 text-sm font-bold text-white">
              {initials}
            </div>
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold text-slate-950">
                {displayName}
              </p>
              <p className="truncate text-xs text-slate-500">{displayEmail}</p>
            </div>
          </Link>
          <button
            type="button"
            onClick={logout}
            className="mt-3 inline-flex h-9 w-full items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white text-sm font-semibold text-slate-600 transition hover:border-violet-200 hover:bg-violet-50 hover:text-violet-700"
          >
            <LogOut className="h-4 w-4" aria-hidden="true" />
            Logout
          </button>
        </div>

        <button
          type="button"
          className="mt-4 flex h-14 items-center gap-4 border-t border-slate-200 px-4 pt-7 text-[15px] font-semibold text-slate-500 transition hover:text-violet-700"
        >
          <ChevronsLeft className="h-5 w-5" aria-hidden="true" />
          Collapse
        </button>
      </aside>

      <nav
        className="fixed inset-x-0 bottom-0 z-50 grid h-24 grid-cols-5 items-end rounded-t-4xl border border-slate-100 bg-white/95 px-5 pb-4 pt-3 shadow-[0_-10px_34px_rgba(15,23,42,0.08)] backdrop-blur lg:hidden"
        aria-label="Mobile navigation"
      >
        {mobileNavItems.map((item) => {
          const Icon = item.icon;
          const active = isActivePath(pathname, item.href);
          const isAdd = item.href === "/expenses/new";

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "relative flex min-h-16 flex-col items-center justify-end gap-1 rounded-lg px-1 text-xs font-semibold transition",
                active ? "text-violet-700" : "text-slate-500",
                isAdd ? "justify-end" : "",
              )}
            >
              {active && !isAdd ? (
                <span className="absolute -top-3 h-1 w-14 rounded-full bg-violet-600" />
              ) : null}
              <span
                className={cn(
                  "relative flex items-center justify-center",
                  isAdd
                    ? "-mt-9 mb-1 h-16 w-16 rounded-full bg-violet-600 text-white shadow-lg shadow-violet-300 ring-8 ring-violet-100"
                    : "h-7 w-7",
                )}
              >
                <Icon className={cn(isAdd ? "h-9 w-9" : "h-7 w-7")} aria-hidden="true" />
                {item.badge ? (
                  <span className="absolute -right-2 top-0 h-2.5 w-2.5 rounded-full bg-violet-600 ring-2 ring-white" />
                ) : null}
              </span>
              <span className="max-w-full truncate">
                {item.label === "Voice Input" ? "Voice Input" : item.label}
              </span>
              {item.badge ? (
                <span className="sr-only">{item.badge}</span>
              ) : null}
            </Link>
          );
        })}
      </nav>
    </>
  );
}
