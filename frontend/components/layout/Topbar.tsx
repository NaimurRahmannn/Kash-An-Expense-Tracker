"use client";

import {
  CalendarDays,
  ChevronDown,
  LogOut,
  Search,
  WalletCards,
} from "lucide-react";
import { Input } from "@/components/ui/Input";
import { useAuth } from "@/hooks/useAuth";

type TopbarProps = {
  title?: string;
};

function getInitials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

export function Topbar({ title }: TopbarProps) {
  const { isLoading, logout, user } = useAuth();
  const displayName = isLoading ? "User" : user?.name ?? "User";
  const displayEmail = isLoading ? "Loading..." : user?.email ?? "";
  const initials = getInitials(displayName) || "U";

  return (
    <header className="sticky top-0 z-30 bg-[#fbfbfd]/95 backdrop-blur">
      <div className="px-4 pb-4 pt-4 sm:px-7 lg:hidden">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex h-14 w-14 items-center justify-center rounded-[14px] bg-linear-to-br from-violet-400 to-violet-700 text-white shadow-lg shadow-violet-200">
              <WalletCards className="h-7 w-7" aria-hidden="true" />
            </div>
            <p className="text-2xl font-bold tracking-[-0.01em] text-violet-700">
              Expense Tracker
            </p>
          </div>

          <button
            type="button"
            onClick={logout}
            className="inline-flex items-center gap-3 rounded-full text-slate-950"
            aria-label="Logout"
          >
            <span className="flex h-16 w-16 items-center justify-center overflow-hidden rounded-full bg-linear-to-br from-violet-200 via-amber-100 to-sky-200 text-base font-bold text-slate-800">
              {initials}
            </span>
            <LogOut className="h-5 w-5" aria-hidden="true" />
          </button>
        </div>
      </div>

      <div className="hidden flex-col gap-4 px-4 py-5 sm:px-7 lg:flex lg:flex-row lg:items-center lg:justify-between lg:px-8">
        <div className="min-w-0">
          {title ? (
            <h1 className="truncate text-3xl font-bold tracking-[-0.01em] text-slate-950">
              {title}
            </h1>
          ) : null}
        </div>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
      

          <div className="relative w-full sm:w-67.5">
            <Search
              className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400"
              aria-hidden="true"
            />
            <Input
              aria-label="Search expenses"
              placeholder="Search expenses..."
              className="h-12 rounded-2xl pl-11"
            />
          </div>

          <button
            type="button"
            className="relative inline-flex h-12 w-12 items-center justify-center rounded-2xl text-slate-700 transition hover:bg-violet-50"
            aria-label="Notifications"
          >
            <span className="relative h-5 w-5">
              <span className="absolute inset-x-1 top-0 h-4 rounded-t-full border-2 border-slate-700 border-b-0" />
              <span className="absolute bottom-0 left-1/2 h-1.5 w-1.5 -translate-x-1/2 rounded-full bg-slate-700" />
            </span>
            <span className="absolute right-3 top-2.5 h-2.5 w-2.5 rounded-full bg-violet-600 ring-2 ring-white" />
          </button>

          <div className="inline-flex h-12 items-center gap-3 rounded-2xl px-1.5 pr-2 text-sm font-semibold text-slate-800">
            <span className="flex h-11 w-11 items-center justify-center overflow-hidden rounded-full bg-linear-to-br from-violet-200 via-amber-100 to-sky-200 text-sm font-bold text-slate-800">
              {initials}
            </span>
            <span className="hidden min-w-0 sm:block">
              <span className="block max-w-34 truncate">{displayName}</span>
              <span className="block max-w-34 truncate text-xs font-medium text-slate-500">
                {displayEmail}
              </span>
            </span>
          </div>

          <button
            type="button"
            onClick={logout}
            className="inline-flex h-10 items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-3 text-sm font-semibold text-slate-600 shadow-sm transition hover:border-violet-200 hover:bg-violet-50 hover:text-violet-700"
          >
            <LogOut className="h-4 w-4" aria-hidden="true" />
            Logout
          </button>
        </div>
      </div>
    </header>
  );
}
