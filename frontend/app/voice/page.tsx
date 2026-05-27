import { Mic } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Card } from "@/components/ui/Card";

export default function VoiceInputPage() {
  return (
    <AppShell>
      <PageHeader
        title="Voice Input"
        description="Bonus voice capture surface. Browser speech integration will be added later."
      />

      <Card className="mx-auto max-w-3xl p-8 text-center">
        <div className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-violet-50 text-violet-600 ring-8 ring-violet-100/70">
          <Mic className="h-9 w-9" aria-hidden="true" />
        </div>
        <h2 className="mt-6 text-xl font-bold text-slate-950">
          Voice input will be added as a bonus feature
        </h2>
        <p className="mx-auto mt-3 max-w-xl text-sm leading-6 text-slate-500">
          This page is only a polished placeholder for now. No microphone access or browser speech
          API is active yet.
        </p>
        <div className="mx-auto mt-8 max-w-xl rounded-lg border border-violet-100 bg-violet-50 px-5 py-4 text-left">
          <p className="text-xs font-bold uppercase text-violet-700">Example transcript</p>
          <p className="mt-2 text-sm font-medium text-slate-700">
            &quot;Had lunch at Pizza Hut for ৳560&quot;
          </p>
        </div>
      </Card>
    </AppShell>
  );
}
