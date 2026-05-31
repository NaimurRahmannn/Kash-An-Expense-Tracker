"use client";

import { Hash, Mail, ShieldCheck, User } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { useAuth } from "@/hooks/useAuth";

function getInitials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

export default function ProfilePage() {
  const { isLoading, user } = useAuth();
  const displayName = isLoading ? "Loading..." : user?.name ?? "User";
  const displayEmail = isLoading ? "Loading..." : user?.email ?? "No email available";
  const displayUserId = isLoading ? "Loading..." : String(user?.user_id ?? "Not available");
  const accessStatus = isLoading
    ? "Checking session"
    : user
      ? "Authenticated user"
      : "Not authenticated";
  const initials = getInitials(user?.name ?? "") || "U";
  const profileDetails = [
    { icon: User, label: "Name", value: displayName },
    { icon: Mail, label: "Email", value: displayEmail },
    { icon: Hash, label: "User ID", value: displayUserId },
    { icon: ShieldCheck, label: "Access", value: accessStatus },
  ];

  return (
    <AppShell>
      <PageHeader
        title="Profile"
        description="Account details for the currently signed-in user."
      />

      <div className="grid gap-6 xl:grid-cols-[360px_1fr]">
        <Card className="p-6 text-center">
          <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-violet-600 text-2xl font-bold text-white">
            {initials}
          </div>
          <h2 className="mt-4 text-xl font-bold text-slate-950">{displayName}</h2>
          <p className="mt-1 break-all text-sm text-slate-500">{displayEmail}</p>
          <div className="mt-4">
            <Badge>{user ? "Signed in" : accessStatus}</Badge>
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
