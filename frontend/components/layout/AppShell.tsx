import type { ReactNode } from "react";
import { Sidebar } from "@/components/layout/Sidebar";
import { Topbar } from "@/components/layout/Topbar";

type AppShellProps = {
  children: ReactNode;
  title?: string;
};

export function AppShell({ children, title }: AppShellProps) {
  return (
    <div className="min-h-screen bg-[#fbfbfd] text-slate-950">
      <Sidebar />
      <div className="lg:pl-73">
        <Topbar title={title} />
        <main className="w-full px-4 py-5 pb-34 sm:px-7 lg:px-8 lg:pb-8">
          {children}
        </main>
      </div>
    </div>
  );
}
