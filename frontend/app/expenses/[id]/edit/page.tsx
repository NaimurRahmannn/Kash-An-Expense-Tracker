import { RotateCcw, Save, Tag } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";

const categories = [
  "Food",
  "Transport",
  "Housing",
  "Entertainment",
  "Shopping",
  "Healthcare",
  "Education",
  "Utilities",
  "Other",
];

const fieldClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-violet-400 focus:ring-4 focus:ring-violet-500/15";

export default function EditExpensePage() {
  return (
    <AppShell>
      <PageHeader
        title="Edit Expense"
        description="Static edit form placeholder. The update endpoint will be connected later."
      />

      <div className="grid gap-6 xl:grid-cols-[1fr_360px]">
        <Card className="p-5 sm:p-6">
          <form className="space-y-5">
            <div className="grid gap-5 md:grid-cols-2">
              <div>
                <label htmlFor="title" className="text-sm font-semibold text-slate-700">
                  Title
                </label>
                <Input id="title" defaultValue="Lunch at Pizza Hut" className="mt-2" />
              </div>
              <div>
                <label htmlFor="amount" className="text-sm font-semibold text-slate-700">
                  Amount
                </label>
                <Input id="amount" type="number" defaultValue="560" className="mt-2" />
              </div>
              <div>
                <label htmlFor="category" className="text-sm font-semibold text-slate-700">
                  Category
                </label>
                <select id="category" className={`${fieldClass} mt-2 h-10`} defaultValue="Food">
                  {categories.map((category) => (
                    <option key={category}>{category}</option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor="expense_date" className="text-sm font-semibold text-slate-700">
                  Expense date
                </label>
                <Input id="expense_date" type="date" defaultValue="2026-05-27" className="mt-2" />
              </div>
            </div>

            <div>
              <label htmlFor="note" className="text-sm font-semibold text-slate-700">
                Note
              </label>
              <textarea
                id="note"
                rows={5}
                defaultValue="Team lunch during the sprint planning break."
                className={`${fieldClass} mt-2 resize-none`}
              />
            </div>

            <div className="flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row">
              <Button type="button">
                <Save className="h-4 w-4" aria-hidden="true" />
                Save Changes
              </Button>
              <Button type="reset" variant="secondary">
                <RotateCcw className="h-4 w-4" aria-hidden="true" />
                Reset
              </Button>
            </div>
          </form>
        </Card>

        <div className="space-y-6">
          <Card className="p-5">
            <div className="flex items-center gap-2">
              <Tag className="h-5 w-5 text-violet-600" aria-hidden="true" />
              <h2 className="text-base font-bold text-slate-950">Expense snapshot</h2>
            </div>
            <div className="mt-4 space-y-3 text-sm text-slate-600">
              <div className="flex items-center justify-between">
                <span>Current category</span>
                <Badge>Food</Badge>
              </div>
              <div className="flex items-center justify-between">
                <span>Created from</span>
                <span className="font-semibold text-slate-950">Static sample</span>
              </div>
              <div className="flex items-center justify-between">
                <span>Status</span>
                <span className="font-semibold text-slate-950">Ready for API wiring</span>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
