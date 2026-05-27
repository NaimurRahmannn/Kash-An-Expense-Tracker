"use client";

import { type FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { formatDateForAPI } from "@/lib/date";
import { createExpense } from "@/lib/expenses";
import type { ExpenseCategory, ExpenseInput } from "@/types/expense";

type CreateExpenseFormValues = {
  title: string;
  amount: string;
  category: ExpenseCategory | "";
  expense_date: string;
  note: string;
};

type CreateExpenseFieldErrors = Partial<
  Record<keyof CreateExpenseFormValues, string>
>;

const defaultCategory: ExpenseCategory = "Food";

function getDefaultValues(): CreateExpenseFormValues {
  return {
    title: "",
    amount: "",
    category: defaultCategory,
    expense_date: formatDateForAPI(new Date()),
    note: "",
  };
}

function isValidAPIDate(dateString: string) {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(dateString)) {
    return false;
  }

  const [year, month, day] = dateString.split("-").map(Number);
  const parsedDate = new Date(year, month - 1, day);

  return (
    parsedDate.getFullYear() === year &&
    parsedDate.getMonth() === month - 1 &&
    parsedDate.getDate() === day
  );
}

function validateForm(values: CreateExpenseFormValues) {
  const errors: CreateExpenseFieldErrors = {};

  if (!values.title.trim()) {
    errors.title = "Title is required";
  }

  if (!values.amount.trim()) {
    errors.amount = "Amount is required";
  } else {
    const amount = Number(values.amount);

    if (!Number.isFinite(amount) || amount <= 0) {
      errors.amount = "Amount must be positive";
    }
  }

  if (!values.category) {
    errors.category = "Category is required";
  }

  if (!values.expense_date) {
    errors.expense_date = "Expense date is required";
  } else if (!isValidAPIDate(values.expense_date)) {
    errors.expense_date = "Invalid expense date format";
  }

  return errors;
}

function getCreateErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to create expense. Please try again.";
  }

  return message;
}

export function useCreateExpenseForm() {
  const router = useRouter();
  const [values, setValues] = useState<CreateExpenseFormValues>(getDefaultValues);
  const [fieldErrors, setFieldErrors] = useState<CreateExpenseFieldErrors>({});
  const [formError, setFormError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  function updateField<K extends keyof CreateExpenseFormValues>(
    field: K,
    value: CreateExpenseFormValues[K],
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value,
    }));
    setFieldErrors((currentErrors) => ({
      ...currentErrors,
      [field]: "",
    }));
    setFormError("");
    setSuccessMessage("");
  }

  function resetForm() {
    setValues(getDefaultValues());
    setFieldErrors({});
    setFormError("");
    setSuccessMessage("");
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationErrors = validateForm(values);
    setFieldErrors(validationErrors);
    setFormError("");
    setSuccessMessage("");

    if (Object.keys(validationErrors).length > 0) {
      return;
    }

    const payload: ExpenseInput = {
      title: values.title.trim(),
      amount: Number(values.amount),
      category: values.category as ExpenseCategory,
      note: values.note.trim(),
      expense_date: values.expense_date,
    };

    setIsSubmitting(true);

    try {
      const response = await createExpense(payload);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to create expense");
      }

      setSuccessMessage("Expense created successfully");
      window.setTimeout(() => {
        router.push("/expenses");
      }, 500);
    } catch (caughtError) {
      setFormError(getCreateErrorMessage(caughtError));
    } finally {
      setIsSubmitting(false);
    }
  }

  return {
    fieldErrors,
    formError,
    handleSubmit,
    isSubmitting,
    resetForm,
    successMessage,
    updateField,
    values,
  };
}
