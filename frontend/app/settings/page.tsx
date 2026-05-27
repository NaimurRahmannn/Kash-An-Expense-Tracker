import { Database, Palette, Server } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { API_BASE_URL } from "@/config/env";

export default function SettingsPage() {
  return (
    <AppShell>
      <PageHeader
        title="Settings"
        description="Static configuration placeholders for the frontend foundation."
      />

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="p-5">
          <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-violet-50 text-violet-600">
            <Server className="h-5 w-5" aria-hidden="true" />
          </div>
          <h2 className="mt-4 text-base font-bold text-slate-950">API Base URL</h2>
          <p className="mt-2 break-all rounded-lg bg-slate-50 px-3 py-2 text-sm font-semibold text-slate-700">
            {API_BASE_URL}
          </p>
        </Card>

        <Card className="p-5">
          <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-cyan-50 text-cyan-600">
            <Palette className="h-5 w-5" aria-hidden="true" />
          </div>
          <h2 className="mt-4 text-base font-bold text-slate-950">Theme</h2>
          <p className="mt-2 text-sm leading-6 text-slate-500">
            Light SaaS dashboard theme placeholder.
          </p>
          <div className="mt-4">
            <Badge>Light mode</Badge>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
            <Database className="h-5 w-5" aria-hidden="true" />
          </div>
          <h2 className="mt-4 text-base font-bold text-slate-950">Data source</h2>
          <p className="mt-2 text-sm leading-6 text-slate-500">
            Pages are static in Frontend Part 1. API integration is intentionally deferred.
          </p>
        </Card>
      </div>
    </AppShell>
  );
}
