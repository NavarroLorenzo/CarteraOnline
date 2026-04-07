export type User = {
  id: number;
  email: string;
  username: string;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export type Account = {
  id: number;
  name: string;
  type: string;
  is_active: boolean;
  created_at: string;
};

export type TransactionType = "income" | "expense";

export type Transaction = {
  id: number;
  title: string;
  amount: number;
  type: TransactionType;
  account_id: number;
  category: string;
  category_label?: string;
  description: string;
  transfer_id?: string;
  created_at: string;
};

export type TransactionCategory = {
  key: string;
  label: string;
  allowed_types: TransactionType[];
};

export type TransactionSummary = {
  income_total: number;
  expense_total: number;
  net_balance: number;
  transactions_count: number;
};

export type DashboardSummary = {
  income_total: number;
  expense_total: number;
  net_balance: number;
  transactions_count: number;
};

export type DashboardComparison = {
  current_month: DashboardSummary;
  previous_month: DashboardSummary;
  income_change_pct: number;
  expense_change_pct: number;
};

export type DashboardTrendPoint = {
  label: string;
  start_date: string;
  end_date: string;
  income_total: number;
  expense_total: number;
};

export type DashboardCategorySummary = {
  key: string;
  label: string;
  amount: number;
  percentage: number;
  transactions_count: number;
};

export type DashboardTransactionItem = {
  id: number;
  title: string;
  amount: number;
  type: TransactionType;
  account_id: number;
  account_name: string;
  category: string;
  category_label: string;
  description: string;
  created_at: string;
};

export type DashboardAnalytics = {
  period_summary: DashboardSummary;
  day_summary: DashboardSummary;
  month_summary: DashboardSummary;
  comparison: DashboardComparison;
  trend_interval: "hourly" | "daily" | "weekly" | "monthly";
  trend: DashboardTrendPoint[];
  expense_categories: DashboardCategorySummary[];
  top_expenses: DashboardTransactionItem[];
  recent_transactions: DashboardTransactionItem[];
};

export type DashboardCategoryDetail = {
  category: DashboardCategorySummary;
  trend_interval: "hourly" | "daily" | "weekly" | "monthly";
  trend: DashboardTrendPoint[];
  transactions: DashboardTransactionItem[];
};

export type AccountBalance = {
  id: number;
  name: string;
  balance: number;
};

export type BalanceResponse = {
  balance: number;
};

export type BalanceByAccountResponse = {
  accounts: AccountBalance[];
  total: number;
};

export type ApiErrorResponse = {
  error: string;
  details?: string;
};
