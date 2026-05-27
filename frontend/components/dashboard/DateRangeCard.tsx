import type { FormEvent } from "react";
import { CalendarDays, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";

type DateRangeCardProps = {
  dateFrom: string;
  dateTo: string;
  isLoading: boolean;
  validationError: string;
  onApply: (event: FormEvent<HTMLFormElement>) => void;
  onDateFromChange: (value: string) => void;
  onDateToChange: (value: string) => void;
};

export function DateRangeCard({
  dateFrom,
  dateTo,
  isLoading,
  onApply,
  onDateFromChange,
  onDateToChange,
  validationError,
}: DateRangeCardProps) {
  return (
    <Card className="rounded-[14px] p-4 shadow-[0_6px_18px_rgba(15,23,42,0.08)]">
      <form className="grid gap-3 lg:grid-cols-[1fr_1fr_auto]" onSubmit={onApply}>
        <div>
          <label htmlFor="date-from" className="text-xs font-bold uppercase text-slate-500">
            Date from
          </label>
          <Input
            id="date-from"
            type="date"
            className="mt-2"
            value={dateFrom}
            onChange={(event) => onDateFromChange(event.target.value)}
          />
        </div>
        <div>
          <label htmlFor="date-to" className="text-xs font-bold uppercase text-slate-500">
            Date to
          </label>
          <Input
            id="date-to"
            type="date"
            className="mt-2"
            value={dateTo}
            onChange={(event) => onDateToChange(event.target.value)}
          />
        </div>
        <div className="flex items-end">
          <Button type="submit" disabled={isLoading} className="w-full lg:w-auto">
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
            ) : (
              <CalendarDays className="h-4 w-4" aria-hidden="true" />
            )}
            Apply
          </Button>
        </div>
      </form>
      {validationError ? (
        <p className="mt-3 text-sm font-medium text-rose-600">{validationError}</p>
      ) : null}
    </Card>
  );
}
