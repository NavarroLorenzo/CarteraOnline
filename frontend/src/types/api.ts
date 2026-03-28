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
  description: string;
  transfer_id?: string;
  created_at: string;
};

export type TransactionSummary = {
  income_total: number;
  expense_total: number;
  net_balance: number;
  transactions_count: number;
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
