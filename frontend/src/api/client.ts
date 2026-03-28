import { API_BASE_URL } from "../lib/config";
import { getStoredToken } from "../lib/storage";
import type { ApiErrorResponse } from "../types/api";

type Primitive = string | number | boolean;

type RequestOptions = {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  query?: Record<string, Primitive | null | undefined>;
  auth?: boolean;
  headers?: HeadersInit;
};

export class ApiError extends Error {
  status: number;
  payload?: ApiErrorResponse | string;

  constructor(message: string, status: number, payload?: ApiErrorResponse | string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

function joinUrl(base: string, path: string): string {
  const normalizedBase = base.endsWith("/") ? base.slice(0, -1) : base;
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return `${normalizedBase}${normalizedPath}`;
}

function buildUrl(path: string, query?: RequestOptions["query"]): URL {
  const target = joinUrl(API_BASE_URL, path);
  const url =
    target.startsWith("http://") || target.startsWith("https://")
      ? new URL(target)
      : new URL(target, window.location.origin);

  if (query) {
    Object.entries(query).forEach(([key, value]) => {
      if (value === undefined || value === null || value === "") {
        return;
      }

      url.searchParams.set(key, String(value));
    });
  }

  return url;
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, query, auth = true, headers } = options;
  const requestHeaders = new Headers(headers);

  if (body !== undefined) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (auth) {
    const token = getStoredToken();
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }
  }

  const response = await fetch(buildUrl(path, query), {
    method,
    headers: requestHeaders,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const isJson = response.headers.get("content-type")?.includes("application/json");
  const payload = isJson ? ((await response.json()) as ApiErrorResponse | T) : await response.text();

  if (!response.ok) {
    const message =
      typeof payload === "object" && payload && "error" in payload
        ? payload.error
        : `La petición falló con estado ${response.status}`;

    throw new ApiError(message, response.status, payload as ApiErrorResponse | string);
  }

  return payload as T;
}
