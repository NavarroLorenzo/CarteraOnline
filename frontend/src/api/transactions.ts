import { apiRequest } from "./client";
import type {
  BalanceByAccountResponse,
  BalanceResponse,
  DashboardAnalytics,
  DashboardCategoryDetail,
  Transaction,
  TransactionCategory,
  TransactionSummary,
  TransactionType,
} from "../types/api";
import type { DashboardFilterPreset } from "../lib/dashboard";

export type CreateTransactionPayload = {
  title: string;
  amount: number;
  type: TransactionType;
  account_id: number;
  category: string;
  description?: string;
};

export type TransactionFilters = {
  account_id?: number;
  type?: TransactionType | "";
  category?: string;
  date_from?: string;
  date_to?: string;
};

export type DashboardPeriodQuery = {
  period: DashboardFilterPreset;
  dateFrom?: string;
  dateTo?: string;
};

export const transactionsApi = {
  list(filters: TransactionFilters = {}) {
    return apiRequest<Transaction[]>("/transactions", {
      query: filters,
    });
  },

  getCategories() {
    return apiRequest<TransactionCategory[]>("/transactions/categories");
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

  getDashboard(filters: DashboardPeriodQuery) {
    return apiRequest<DashboardAnalytics>("/transactions/dashboard", {
      query: filters,
    });
  },

  getDashboardCategoryDetail(categoryKey: string, filters: DashboardPeriodQuery) {
    return apiRequest<DashboardCategoryDetail>("/transactions/dashboard/category", {
      query: {
        ...filters,
        category_key: categoryKey,
      },
    });
  },

  getBalanceByAccount() {
    return apiRequest<BalanceByAccountResponse>("/transactions/balance-by-account");
  },
};
