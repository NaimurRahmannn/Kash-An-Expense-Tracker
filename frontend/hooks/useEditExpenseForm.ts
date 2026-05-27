"use client";

import { type FormEvent, useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { getExpenseById, updateExpense } from "@/lib/expenses";
import {
  expenseToFormValues,
  getDefaultExpenseFormValues,
  hasExpenseFormErrors,
  toExpenseInput,
  validateExpenseForm,
  type ExpenseFormErrors,
  type ExpenseFormValues,
} from "@/lib/expense-validation";
import type { Expense } from "@/types/expense";

function getLoadErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to load expense. Please try again.";
  }

  return message;
}

function getUpdateErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to update expense. Please try again.";
  }

  return message;
}

export function useEditExpenseForm(expenseId: number | null) {
  const router = useRouter();
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const [values, setValues] = useState<ExpenseFormValues>(
    getDefaultExpenseFormValues,
  );
  const [originalValues, setOriginalValues] =
    useState<ExpenseFormValues | null>(null);
  const [currentExpense, setCurrentExpense] = useState<Expense | null>(null);
  const [fieldErrors, setFieldErrors] = useState<ExpenseFormErrors>({});
  const [loadError, setLoadError] = useState("");
  const [formError, setFormError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [isLoadingExpense, setIsLoadingExpense] = useState(() => expenseId !== null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const loadExpense = useCallback(async () => {
    if (expenseId === null) {
      return;
    }

    setLoadError("");
    setIsLoadingExpense(true);

    try {
      const response = await getExpenseById(expenseId);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to load expense");
      }

      const nextValues = expenseToFormValues(response.data);
      setCurrentExpense(response.data);
      setOriginalValues(nextValues);
      setValues(nextValues);
      setFieldErrors({});
      setFormError("");
      setSuccessMessage("");
    } catch (caughtError) {
      setLoadError(getLoadErrorMessage(caughtError));
    } finally {
      setIsLoadingExpense(false);
    }
  }, [expenseId]);

  useEffect(() => {
    if (expenseId === null || isAuthLoading || !isAuthenticated) {
      return;
    }

    const timeoutId = window.setTimeout(() => {
      void loadExpense();
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, [expenseId, isAuthLoading, isAuthenticated, loadExpense]);

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
    if (!originalValues) {
      return;
    }

    setValues(originalValues);
    setFieldErrors({});
    setFormError("");
    setSuccessMessage("");
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (expenseId === null) {
      setFormError("Invalid expense ID");
      return;
    }

    const validationErrors = validateExpenseForm(values);
    setFieldErrors(validationErrors);
    setFormError("");
    setSuccessMessage("");

    if (hasExpenseFormErrors(validationErrors)) {
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await updateExpense(expenseId, toExpenseInput(values));

      if (!response.success || !response.data) {
        throw new Error(response.message || "Unable to update expense");
      }

      const nextValues = expenseToFormValues(response.data);
      setCurrentExpense(response.data);
      setOriginalValues(nextValues);
      setValues(nextValues);
      setSuccessMessage("Expense updated successfully");
      window.setTimeout(() => {
        router.push("/expenses");
      }, 500);
    } catch (caughtError) {
      setFormError(getUpdateErrorMessage(caughtError));
    } finally {
      setIsSubmitting(false);
    }
  }

  return {
    currentExpense,
    fieldErrors,
    formError,
    handleSubmit,
    isLoadingExpense,
    isSubmitting,
    loadError,
    loadExpense,
    resetForm,
    successMessage,
    updateField,
    values,
  };
}
