import type { LucideIcon } from "lucide-react";

export type DateRange = {
  dateFrom: string;
  dateTo: string;
};

export type StatAccent = "violet" | "blue" | "green" | "amber";

export type DashboardStat = {
  title: string;
  value: string;
  detail: string;
  icon: LucideIcon;
  accent: StatAccent;
};
