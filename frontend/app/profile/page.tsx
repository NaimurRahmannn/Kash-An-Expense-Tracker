import { CalendarDays, Mail, Shield, User } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";

const profileDetails = [
  { icon: User, label: "Name", value: "John Doe" },
  { icon: Mail, label: "Email", value: "john.doe@example.com" },
  { icon: CalendarDays, label: "Member since", value: "May 2026" },
  { icon: Shield, label: "Access", value: "Frontend placeholder" },
];

export default function ProfilePage() {
  return (
    <AppShell>
      <PageHeader
        title="Profile"
        description="Static account profile placeholder for the future authenticated user."
      />

      <div className="grid gap-6 xl:grid-cols-[360px_1fr]">
        <Card className="p-6 text-center">
          <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-violet-600 text-2xl font-bold text-white">
            JD
          </div>
          <h2 className="mt-4 text-xl font-bold text-slate-950">John Doe</h2>
          <p className="mt-1 text-sm text-slate-500">john.doe@example.com</p>
          <div className="mt-4">
            <Badge>Personal account</Badge>
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="text-base font-bold text-slate-950">Account details</h2>
          <div className="mt-5 grid gap-4 md:grid-cols-2">
            {profileDetails.map(({ icon: Icon, label, value }) => (
              <div key={label} className="rounded-lg border border-slate-100 bg-slate-50 p-4">
                <div className="flex items-center gap-2 text-sm font-semibold text-slate-500">
                  <Icon className="h-4 w-4 text-violet-600" aria-hidden="true" />
                  {label}
                </div>
                <p className="mt-2 text-sm font-bold text-slate-950">{value}</p>
              </div>
            ))}
          </div>
        </Card>
      </div>
    </AppShell>
  );
}
