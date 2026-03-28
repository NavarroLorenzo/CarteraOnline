import { apiRequest } from "./client";
import type { Account } from "../types/api";

export type CreateAccountPayload = {
  name: string;
  type: string;
  initial_amount: number;
};

export type UpdateAccountPayload = {
  name: string;
  type: string;
};

export const accountsApi = {
  list() {
    return apiRequest<Account[]>("/accounts");
  },

  create(payload: CreateAccountPayload) {
    return apiRequest<Account>("/accounts", {
      method: "POST",
      body: payload,
    });
  },

  update(id: number, payload: UpdateAccountPayload) {
    return apiRequest<Account>(`/accounts/${id}`, {
      method: "PUT",
      body: payload,
    });
  },

  delete(id: number) {
    return apiRequest<{ message: string; warning?: string }>(`/accounts/${id}`, {
      method: "DELETE",
    });
  },
};
