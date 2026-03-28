import { apiRequest } from "./client";
import type { AuthResponse, User } from "../types/api";

export type LoginPayload = {
  identifier: string;
  password: string;
};

export type RegisterPayload = {
  email: string;
  username: string;
  password: string;
};

export const authApi = {
  login(payload: LoginPayload) {
    return apiRequest<AuthResponse>("/auth/login", {
      method: "POST",
      body: payload,
      auth: false,
    });
  },

  register(payload: RegisterPayload) {
    return apiRequest<AuthResponse>("/auth/register", {
      method: "POST",
      body: payload,
      auth: false,
    });
  },

  me() {
    return apiRequest<User>("/auth/me");
  },
};
