import { apiRequest } from "./client";
import type {
  BalanceByAccountResponse,
  BalanceResponse,
  Transaction,
  TransactionSummary,
  TransactionType,
} from "../types/api";

export type CreateTransactionPayload = {
  title: string;
  amount: number;
  type: TransactionType;
  account_id: number;
  category: string;
  description: string;
};

export type TransactionFilters = {
  account_id?: number;
  type?: TransactionType | "";
  category?: string;
  date_from?: string;
  date_to?: string;
};

export const transactionsApi = {
  list(filters: TransactionFilters = {}) {
    return apiRequest<Transaction[]>("/transactions", {
      query: filters,
    });
  },

  create(payload: CreateTransactionPayload) {
    return apiRequest<Transaction>("/transactions", {
      method: "POST",
      body: payload,
    });
  },

  delete(id: number) {
    return apiRequest<{ message: string }>(`/transactions/${id}`, {
      method: "DELETE",
    });
  },

  getBalance() {
    return apiRequest<BalanceResponse>("/transactions/balance");
  },

  getSummary(filters: TransactionFilters = {}) {
    return apiRequest<TransactionSummary>("/transactions/summary", {
      query: filters,
    });
  },

  getBalanceByAccount() {
    return apiRequest<BalanceByAccountResponse>("/transactions/balance-by-account");
  },
};
