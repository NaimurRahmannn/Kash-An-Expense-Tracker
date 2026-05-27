"use client";

import { type FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { createExpense } from "@/lib/expenses";
import {
  getDefaultExpenseFormValues,
  hasExpenseFormErrors,
  toExpenseInput,
  validateExpenseForm,
  type ExpenseFormErrors,
  type ExpenseFormValues,
} from "@/lib/expense-validation";

function getCreateErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to create expense. Please try again.";
  }

  return message;
}

export function useCreateExpenseForm() {
  const router = useRouter();
  const [values, setValues] = useState<ExpenseFormValues>(getDefaultExpenseFormValues);
  const [fieldErrors, setFieldErrors] = useState<ExpenseFormErrors>({});
  const [formError, setFormError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

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
    setFormError("");
    setSuccessMessage("");
  }

  function resetForm() {
    setValues(getDefaultExpenseFormValues());
    setFieldErrors({});
    setFormError("");
    setSuccessMessage("");
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const validationErrors = validateExpenseForm(values);
    setFieldErrors(validationErrors);
    setFormError("");
    setSuccessMessage("");

    if (hasExpenseFormErrors(validationErrors)) {
      return;
    }

    setIsSubmitting(true);

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
