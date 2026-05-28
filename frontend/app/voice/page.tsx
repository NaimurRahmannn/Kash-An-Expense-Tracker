"use client";

import Link from "next/link";
import { type FormEvent, useEffect, useRef, useState } from "react";
import {
  CheckCircle2,
  CirclePlus,
  Info,
  Loader2,
  Mic,
  RotateCcw,
  Save,
  ShieldAlert,
  Sparkles,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import {
  expenseCategories,
  getCategoryBadgeStyle,
} from "@/lib/expense-format";
import {
  hasExpenseFormErrors,
  toExpenseInput,
  validateExpenseForm,
  type ExpenseFormErrors,
  type ExpenseFormValues,
} from "@/lib/expense-validation";
import { createExpense } from "@/lib/expenses";
import {
  createSpeechRecognition,
  isSpeechRecognitionSupported,
} from "@/lib/speech";
import {
  parseExpenseFromTranscript,
  type ParsedVoiceExpense,
} from "@/lib/voice-parser";
import type { ExpenseCategory } from "@/types/expense";

const examples = [
  "Lunch 350 food today",
  "Bus fare 80 transport yesterday",
  "Medicine 500 healthcare 2025-06-10",
  "Netflix 200 entertainment today",
  "Electricity bill 1200 utilities today",
  "Health care bill 900 health care today",
  "Spent on medicine 500 healthcare today",
];

const emptyVoiceValues: ExpenseFormValues = {
  title: "",
  amount: "",
  category: "",
  expense_date: "",
  note: "",
};

const fieldClass =
  "w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-violet-400 focus:ring-4 focus:ring-violet-500/15";

function FieldError({ message }: { message?: string }) {
  return message ? (
    <p className="mt-2 text-sm font-medium text-rose-600">{message}</p>
  ) : null;
}

function getCreateErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to create expense. Please try again.";
  }

  return message;
}

function parsedToFormValues(parsed: ParsedVoiceExpense): ExpenseFormValues {
  return {
    title: parsed.title,
    amount: parsed.amount,
    category: parsed.category,
    expense_date: parsed.expense_date,
    note: parsed.note,
  };
}

export default function VoiceInputPage() {
  const router = useRouter();
  const recognitionRef = useRef<SpeechRecognition | null>(null);
  const receivedResultRef = useRef(false);
  const recognitionErroredRef = useRef(false);
  const [isSupported, setIsSupported] = useState(false);
  const [supportChecked, setSupportChecked] = useState(false);
  const [isListening, setIsListening] = useState(false);
  const [transcript, setTranscript] = useState("");
  const [parsedExpense, setParsedExpense] = useState<ParsedVoiceExpense | null>(
    null,
  );
  const [values, setValues] = useState<ExpenseFormValues>(emptyVoiceValues);
  const [fieldErrors, setFieldErrors] = useState<ExpenseFormErrors>({});
  const [captureError, setCaptureError] = useState("");
  const [saveError, setSaveError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      setIsSupported(isSpeechRecognitionSupported());
      setSupportChecked(true);
    }, 0);

    return () => {
      window.clearTimeout(timeoutId);
      recognitionRef.current?.abort();
    };
  }, []);

  function setParsedResult(nextTranscript: string) {
    const parsed = parseExpenseFromTranscript(nextTranscript);
    setTranscript(nextTranscript.trim());
    setParsedExpense(parsed);
    setValues(parsedToFormValues(parsed));
    setFieldErrors({});
    setSaveError("");
    setSuccessMessage("");
  }

  function updateField<K extends keyof ExpenseFormValues>(
    field: K,
    value: ExpenseFormValues[K],
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value,
    }));
    setFieldErrors((currentErrors) => ({
      ...currentErrors,
      [field]: "",
    }));
    setSaveError("");
    setSuccessMessage("");
  }

  function handleStartRecording() {
    setCaptureError("");
    setSaveError("");
    setSuccessMessage("");
    receivedResultRef.current = false;
    recognitionErroredRef.current = false;

    const recognition = createSpeechRecognition();

    if (!recognition) {
      setCaptureError(
        "Voice input is not supported in this browser. Please use Chrome or add expenses manually.",
      );
      return;
    }

    recognitionRef.current = recognition;

    recognition.onresult = (event) => {
      receivedResultRef.current = true;
      const nextTranscript = event.results?.[0]?.[0]?.transcript?.trim() || "";
      setParsedResult(nextTranscript);
      setIsListening(false);
    };

    recognition.onerror = () => {
      recognitionErroredRef.current = true;
      setCaptureError("Unable to capture voice input. Please try again.");
      setIsListening(false);
    };

    recognition.onend = () => {
      setIsListening(false);

      if (!receivedResultRef.current && !recognitionErroredRef.current) {
        setParsedResult("");
      }
    };

    try {
      setIsListening(true);
      recognition.start();
    } catch {
      setCaptureError("Unable to capture voice input. Please try again.");
      setIsListening(false);
    }
  }

  function handleReset() {
    recognitionRef.current?.abort();
    setIsListening(false);
    setTranscript("");
    setParsedExpense(null);
    setValues(emptyVoiceValues);
    setFieldErrors({});
    setCaptureError("");
    setSaveError("");
    setSuccessMessage("");
  }

  async function handleSave(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationErrors = validateExpenseForm(values);
    setFieldErrors(validationErrors);
    setSaveError("");
    setSuccessMessage("");

    if (hasExpenseFormErrors(validationErrors)) {
      return;
    }

    setIsSaving(true);

    try {
      const response = await createExpense(toExpenseInput(values));

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to create expense");
      }

      setSuccessMessage("Expense created successfully");
      window.setTimeout(() => {
        router.push("/expenses");
      }, 500);
    } catch (caughtError) {
      setSaveError(getCreateErrorMessage(caughtError));
    } finally {
      setIsSaving(false);
    }
  }

  return (
    <AppShell title="Voice Input">
      <PageHeader
        title="Voice Input"
        description="Add expenses faster by speaking in a simple structured format."
      />

      <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
        <div className="space-y-6">
          <Card className="p-5 sm:p-6">
            <div className="flex items-center gap-2">
              <Info className="h-5 w-5 text-violet-600" aria-hidden="true" />
              <h2 className="text-lg font-bold text-slate-950">
                Use this format
              </h2>
            </div>
            <div className="mt-4 rounded-lg border border-violet-100 bg-violet-50 px-4 py-3">
              <p className="text-base font-bold text-violet-900">
                Title Amount Category Date
              </p>
            </div>
            <div className="mt-5 space-y-2">
              {examples.map((example) => (
                <p key={example} className="text-sm font-medium text-slate-600">
                  {example}
                </p>
              ))}
            </div>
            <div className="mt-5 space-y-3 border-t border-slate-100 pt-5">
              <div>
                <p className="text-xs font-bold uppercase text-slate-500">
                  Allowed categories
                </p>
                <div className="mt-2 flex flex-wrap gap-2">
                  {expenseCategories.map((category) => (
                    <Badge
                      key={category}
                      className={getCategoryBadgeStyle(category)}
                    >
                      {category}
                    </Badge>
                  ))}
                </div>
              </div>
              <p className="text-sm font-medium text-slate-600">
                Allowed dates: today, yesterday, or YYYY-MM-DD
              </p>
            </div>
          </Card>

          {supportChecked && !isSupported ? (
            <Card className="border-amber-200 bg-amber-50 p-5">
              <div className="flex items-start gap-3">
                <ShieldAlert
                  className="mt-0.5 h-5 w-5 shrink-0 text-amber-600"
                  aria-hidden="true"
                />
                <div>
                  <h2 className="text-base font-bold text-amber-950">
                    Voice input is not supported in this browser.
                  </h2>
                  <p className="mt-2 text-sm leading-6 text-amber-800">
                    Please use Chrome or add expenses manually.
                  </p>
                  <Link
                    href="/expenses/new"
                    className="mt-4 inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
                  >
                    <CirclePlus className="h-4 w-4" aria-hidden="true" />
                    Add Expense
                  </Link>
                </div>
              </div>
            </Card>
          ) : null}

          <Card className="p-6 text-center">
            <button
              type="button"
              disabled={isListening || (supportChecked && !isSupported)}
              className="mx-auto flex h-24 w-24 items-center justify-center rounded-full bg-violet-600 text-white shadow-lg shadow-violet-200 ring-10 ring-violet-100/80 transition hover:bg-violet-700 disabled:pointer-events-none disabled:opacity-60"
              onClick={handleStartRecording}
              aria-label={isListening ? "Listening" : "Start recording"}
            >
              {isListening ? (
                <Loader2 className="h-10 w-10 animate-spin" aria-hidden="true" />
              ) : (
                <Mic className="h-10 w-10" aria-hidden="true" />
              )}
            </button>
            <h2 className="mt-6 text-xl font-bold text-slate-950">
              {isListening ? "Listening..." : "Start Recording"}
            </h2>
            <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-slate-500">
              {isListening
                ? "Listening for your expense..."
                : "Click the microphone and say your expense."}
            </p>
            {captureError ? (
              <p className="mt-4 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
                {captureError}
              </p>
            ) : null}
          </Card>

          <Card className="p-5">
            <h2 className="text-lg font-bold text-slate-950">Transcript</h2>
            <div className="mt-4 min-h-24 rounded-lg border border-slate-200 bg-slate-50 px-4 py-3 text-sm leading-6 text-slate-700">
              {transcript || "Your transcript will appear here."}
            </div>
          </Card>
        </div>

        <div className="space-y-6">
          {parsedExpense ? (
            <Card className="p-5 sm:p-6">
              <div className="flex flex-col gap-3 border-b border-slate-100 pb-5 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <Sparkles className="h-5 w-5 text-violet-600" aria-hidden="true" />
                    <h2 className="text-lg font-bold text-slate-950">
                      Review Parsed Expense
                    </h2>
                  </div>
                  <p className="mt-2 text-sm leading-6 text-slate-500">
                    Review and edit these fields before saving.
                  </p>
                </div>
                <Badge
                  className={
                    parsedExpense.confidence === "high"
                      ? "bg-emerald-50 text-emerald-700 ring-emerald-200/80"
                      : "bg-amber-50 text-amber-700 ring-amber-200/80"
                  }
                >
                  {parsedExpense.confidence === "high"
                    ? "High confidence"
                    : "Low confidence"}
                </Badge>
              </div>

              {parsedExpense.error ? (
                <div className="mt-5 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm font-semibold text-amber-800">
                  {parsedExpense.error}
                </div>
              ) : null}

              {parsedExpense.confidence === "low" ? (
                <div className="mt-3 rounded-lg border border-violet-100 bg-violet-50 px-4 py-3 text-sm font-semibold text-violet-800">
                  Please review and fix the fields before saving.
                </div>
              ) : null}

              {saveError ? (
                <div className="mt-5 rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
                  {saveError}
                </div>
              ) : null}

              {successMessage ? (
                <div className="mt-5 flex items-center gap-2 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
                  <CheckCircle2 className="h-4 w-4" aria-hidden="true" />
                  {successMessage}
                </div>
              ) : null}

              <form className="mt-5 space-y-5" onSubmit={handleSave}>
                <div className="grid gap-5 md:grid-cols-2">
                  <div>
                    <label
                      htmlFor="voice-title"
                      className="text-sm font-semibold text-slate-700"
                    >
                      Title
                    </label>
                    <Input
                      id="voice-title"
                      className="mt-2"
                      placeholder="Lunch"
                      value={values.title}
                      onChange={(event) =>
                        updateField("title", event.target.value)
                      }
                    />
                    <FieldError message={fieldErrors.title} />
                  </div>

                  <div>
                    <label
                      htmlFor="voice-amount"
                      className="text-sm font-semibold text-slate-700"
                    >
                      Amount
                    </label>
                    <Input
                      id="voice-amount"
                      className="mt-2"
                      type="number"
                      min="0.01"
                      step="0.01"
                      placeholder="350"
                      value={values.amount}
                      onChange={(event) =>
                        updateField("amount", event.target.value)
                      }
                    />
                    <FieldError message={fieldErrors.amount} />
                  </div>

                  <div>
                    <label
                      htmlFor="voice-category"
                      className="text-sm font-semibold text-slate-700"
                    >
                      Category
                    </label>
                    <select
                      id="voice-category"
                      className={`${fieldClass} mt-2 h-10`}
                      value={values.category}
                      onChange={(event) =>
                        updateField(
                          "category",
                          event.target.value as ExpenseCategory | "",
                        )
                      }
                    >
                      <option value="">Select category</option>
                      {expenseCategories.map((category) => (
                        <option key={category} value={category}>
                          {category}
                        </option>
                      ))}
                    </select>
                    <FieldError message={fieldErrors.category} />
                  </div>

                  <div>
                    <label
                      htmlFor="voice-expense-date"
                      className="text-sm font-semibold text-slate-700"
                    >
                      Expense date
                    </label>
                    <Input
                      id="voice-expense-date"
                      className="mt-2"
                      type="date"
                      value={values.expense_date}
                      onChange={(event) =>
                        updateField("expense_date", event.target.value)
                      }
                    />
                    <FieldError message={fieldErrors.expense_date} />
                  </div>
                </div>

                <div>
                  <label
                    htmlFor="voice-note"
                    className="text-sm font-semibold text-slate-700"
                  >
                    Note
                  </label>
                  <textarea
                    id="voice-note"
                    rows={4}
                    placeholder="Voice note"
                    className={`${fieldClass} mt-2 resize-none`}
                    value={values.note}
                    onChange={(event) => updateField("note", event.target.value)}
                  />
                </div>

                <div className="flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:items-center">
                  <Button type="submit" disabled={isSaving}>
                    <Save className="h-4 w-4" aria-hidden="true" />
                    {isSaving ? "Saving..." : "Save Expense"}
                  </Button>
                  <Button
                    type="button"
                    variant="secondary"
                    disabled={isSaving}
                    onClick={handleReset}
                  >
                    <RotateCcw className="h-4 w-4" aria-hidden="true" />
                    Reset
                  </Button>
                </div>
              </form>
            </Card>
          ) : (
            <Card className="flex min-h-80 items-center justify-center p-8 text-center">
              <div>
                <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-lg bg-violet-50 text-violet-600">
                  <Mic className="h-7 w-7" aria-hidden="true" />
                </div>
                <h2 className="mt-5 text-lg font-bold text-slate-950">
                  Parsed fields will appear here
                </h2>
                <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-slate-500">
                  Speak using the guided format, then review the fields before saving.
                </p>
              </div>
            </Card>
          )}
        </div>
      </div>
    </AppShell>
  );
}
