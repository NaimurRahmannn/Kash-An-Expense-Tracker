"use client";

import Link from "next/link";
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
import { EmptyState } from "@/components/ui/EmptyState";
import { useExpensesList } from "@/hooks/useExpensesList";

export default function ExpensesPage() {
  const expenses = useExpensesList();
  const shouldShowPagination =
    expenses.canGoPrevious || expenses.expenses.length > 0;

  function handleDeletePlaceholder() {
    window.alert("Delete will be implemented in the next part.");
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
        <ExpenseFilters
          category={expenses.category}
          dateFrom={expenses.dateFrom}
          dateTo={expenses.dateTo}
          filterError={expenses.filterError}
          isLoading={expenses.isLoading}
          searchTerm={expenses.searchTerm}
          sortBy={expenses.sortBy}
          sortOrder={expenses.sortOrder}
          onApply={expenses.handleApplyFilters}
          onCategoryChange={expenses.setCategory}
          onDateFromChange={expenses.setDateFrom}
          onDateToChange={expenses.setDateTo}
          onReset={expenses.handleResetFilters}
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
                onDeletePlaceholder={handleDeletePlaceholder}
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
    </AppShell>
  );
}
