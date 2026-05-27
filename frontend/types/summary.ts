export type CategorySummary = {
  category: string;
  total: number;
  count: number;
};

export type ExpenseSummary = {
  date_from: string;
  date_to: string;
  total_amount: number;
  total_count: number;
  by_category: CategorySummary[];
};
