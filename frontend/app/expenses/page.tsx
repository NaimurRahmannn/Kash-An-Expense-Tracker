"use client";

import Link from "next/link";
import { useState } from "react";
import { CirclePlus, PlusCircle, ReceiptText } from "lucide-react";
import { ExpenseFilters } from "@/components/expenses/ExpenseFilters";
import { ExpensePagination } from "@/components/expenses/ExpensePagination";
import { ExpenseTable } from "@/components/expenses/ExpenseTable";
import {
  ExpensesError,
  ExpensesLoading,
} from "@/components/expenses/ExpensesState";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Card } from "@/components/ui/Card";
import { ConfirmDialog } from "@/components/ui/ConfirmDialog";
import { EmptyState } from "@/components/ui/EmptyState";
import { useExpensesList } from "@/hooks/useExpensesList";
import { deleteExpense } from "@/lib/expenses";
import type { Expense } from "@/types/expense";

function getDeleteErrorMessage(error: unknown) {
  const message = error instanceof Error ? error.message : "";

  if (!message || message === "Failed to fetch" || message === "Request failed") {
    return "Unable to delete expense. Please try again.";
  }

  return message;
}

export default function ExpensesPage() {
  const expenses = useExpensesList();
  const [selectedExpense, setSelectedExpense] = useState<Expense | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const shouldShowPagination =
    expenses.canGoPrevious || expenses.expenses.length > 0;

  function clearFeedback() {
    setDeleteError(null);
    setSuccessMessage(null);
  }

  function handleDeleteClick(expense: Expense) {
    clearFeedback();
    setSelectedExpense(expense);
  }

  function handleCancelDelete() {
    if (isDeleting) {
      return;
    }

    setSelectedExpense(null);
  }

  async function handleConfirmDelete() {
    if (!selectedExpense) {
      return;
    }

    setIsDeleting(true);
    setDeleteError(null);
    setSuccessMessage(null);

    try {
      const response = await deleteExpense(selectedExpense.id);

      if (!response.success) {
        throw new Error(response.message || "Unable to delete expense");
      }

      setSelectedExpense(null);
      setSuccessMessage("Expense deleted successfully");
      await expenses.refreshCurrentPage();
    } catch (caughtError) {
      setSelectedExpense(null);
      setDeleteError(getDeleteErrorMessage(caughtError));
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <AppShell title="Expenses">
      <PageHeader
        title="Expenses"
        description="View and manage all your expenses."
        action={
          <Link
            href="/expenses/new"
            className="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
          >
            <PlusCircle className="h-4 w-4" aria-hidden="true" />
            Add Expense
          </Link>
        }
      />

      <div className="space-y-5">
        {successMessage ? (
          <div className="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
            {successMessage}
          </div>
        ) : null}

        {deleteError ? (
          <div className="rounded-lg border border-rose-200 bg-rose-50 px-4 py-3 text-sm font-semibold text-rose-700">
            {deleteError}
          </div>
        ) : null}

        <ExpenseFilters
          category={expenses.category}
          dateFrom={expenses.dateFrom}
          dateTo={expenses.dateTo}
          filterError={expenses.filterError}
          isLoading={expenses.isLoading}
          searchTerm={expenses.searchTerm}
          sortBy={expenses.sortBy}
          sortOrder={expenses.sortOrder}
          onApply={(event) => {
            clearFeedback();
            expenses.handleApplyFilters(event);
          }}
          onCategoryChange={expenses.setCategory}
          onDateFromChange={expenses.setDateFrom}
          onDateToChange={expenses.setDateTo}
          onReset={() => {
            clearFeedback();
            expenses.handleResetFilters();
          }}
          onSearchChange={expenses.setSearchTerm}
          onSortByChange={expenses.setSortBy}
          onSortOrderChange={expenses.setSortOrder}
        />

        {expenses.isLoading ? (
          <ExpensesLoading />
        ) : expenses.error ? (
          <ExpensesError error={expenses.error} onRetry={expenses.retry} />
        ) : (
          <Card className="p-4 sm:p-5">
            {expenses.visibleExpenses.length === 0 ? (
              <EmptyState
                icon={<ReceiptText className="h-5 w-5" aria-hidden="true" />}
                title="No expenses found"
                description="Try changing filters or add your first expense."
                action={
                  <Link
                    href="/expenses/new"
                    className="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 text-sm font-semibold text-white shadow-sm shadow-violet-200 transition hover:bg-violet-700"
                  >
                    <CirclePlus className="h-4 w-4" aria-hidden="true" />
                    Add Expense
                  </Link>
                }
              />
            ) : (
              <ExpenseTable
                expenses={expenses.visibleExpenses}
                onDelete={handleDeleteClick}
              />
            )}

            {shouldShowPagination ? (
              <ExpensePagination
                canGoNext={expenses.canGoNext}
                canGoPrevious={expenses.canGoPrevious}
                isLoading={expenses.isLoading}
                itemCount={expenses.visibleExpenses.length}
                limit={expenses.limit}
                page={expenses.page}
                onNext={expenses.handleNextPage}
                onPrevious={expenses.handlePreviousPage}
              />
            ) : null}
          </Card>
        )}
      </div>

      <ConfirmDialog
        open={selectedExpense !== null}
        title="Delete expense?"
        description={
          selectedExpense
            ? `Are you sure you want to delete "${selectedExpense.title}"? This action cannot be undone.`
            : ""
        }
        confirmText="Delete"
        cancelText="Cancel"
        isLoading={isDeleting}
        onConfirm={handleConfirmDelete}
        onCancel={handleCancelDelete}
      />
    </AppShell>
  );
}
