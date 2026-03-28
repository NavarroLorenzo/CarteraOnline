import type { User } from "../types/api";

const AUTH_STORAGE_KEY = "cartera-online.auth";

export type StoredAuth = {
  token: string;
  user: User | null;
};

export function loadStoredAuth(): StoredAuth | null {
  const raw = localStorage.getItem(AUTH_STORAGE_KEY);
  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw) as StoredAuth;
  } catch {
    localStorage.removeItem(AUTH_STORAGE_KEY);
    return null;
  }
}

export function saveStoredAuth(auth: StoredAuth): void {
  localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(auth));
}

export function clearStoredAuth(): void {
  localStorage.removeItem(AUTH_STORAGE_KEY);
}

export function getStoredToken(): string | null {
  return loadStoredAuth()?.token ?? null;
}
