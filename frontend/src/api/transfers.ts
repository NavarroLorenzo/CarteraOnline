import { apiRequest } from "./client";

export type CreateTransferPayload = {
  from_account_id: number;
  to_account_id: number;
  amount: number;
  description: string;
};

export const transfersApi = {
  create(payload: CreateTransferPayload) {
    return apiRequest<{ message: string }>("/transfers", {
      method: "POST",
      body: payload,
    });
  },
};
