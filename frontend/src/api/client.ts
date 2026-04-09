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
  retry?: boolean;
  timeoutMs?: number;
};

export type ServerWakeState = {
  isWaking: boolean;
  activeRequests: number;
  retryCount: number;
};

const REQUEST_TIMEOUT_MS = 12_000;
const SERVER_WAKE_THRESHOLD_MS = 2_500;
const RETRY_DELAYS_MS = [2_000, 3_000, 5_000, 5_000, 5_000];

let nextWakeRequestId = 1;
const wakeListeners = new Set<() => void>();
const wakeRequests = new Map<number, number>();
let wakeSnapshot: ServerWakeState = {
  isWaking: false,
  activeRequests: 0,
  retryCount: 0,
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

export function subscribeToServerWakeState(listener: () => void): () => void {
  wakeListeners.add(listener);

  return () => {
    wakeListeners.delete(listener);
  };
}

export function getServerWakeStateSnapshot(): ServerWakeState {
  return wakeSnapshot;
}

function publishWakeState() {
  const nextSnapshot: ServerWakeState = {
    isWaking: wakeRequests.size > 0,
    activeRequests: wakeRequests.size,
    retryCount: Math.max(0, ...wakeRequests.values()),
  };

  if (
    nextSnapshot.isWaking === wakeSnapshot.isWaking &&
    nextSnapshot.activeRequests === wakeSnapshot.activeRequests &&
    nextSnapshot.retryCount === wakeSnapshot.retryCount
  ) {
    return;
  }

  wakeSnapshot = nextSnapshot;
  wakeListeners.forEach((listener) => listener());
}

function markWakeRequest(requestId: number, retryCount: number) {
  wakeRequests.set(requestId, retryCount);
  publishWakeState();
}

function clearWakeRequest(requestId: number) {
  if (wakeRequests.delete(requestId)) {
    publishWakeState();
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

async function wait(ms: number) {
  await new Promise((resolve) => window.setTimeout(resolve, ms));
}

async function fetchWithTimeout(input: URL, init: RequestInit, timeoutMs: number): Promise<Response> {
  const controller = new AbortController();
  const timeoutId = window.setTimeout(() => controller.abort(), timeoutMs);

  try {
    return await fetch(input, {
      ...init,
      signal: controller.signal,
    });
  } finally {
    window.clearTimeout(timeoutId);
  }
}

async function readPayload<T>(response: Response): Promise<ApiErrorResponse | T | string> {
  if (response.status === 204 || response.status === 205) {
    return "";
  }

  const isJson = response.headers.get("content-type")?.includes("application/json");

  if (isJson) {
    return (await response.json()) as ApiErrorResponse | T;
  }

  return await response.text();
}

function shouldRetryStatus(status: number): boolean {
  return status === 408 || status === 425 || status === 429 || status >= 500;
}

function shouldRetryError(error: unknown): boolean {
  return error instanceof TypeError || (error instanceof DOMException && error.name === "AbortError");
}

function getStatusFallbackMessage(status: number): string {
  if (status === 401) {
    return "Tu sesión venció. Iniciá sesión nuevamente para continuar.";
  }

  if (status === 403) {
    return "No tenés permisos para realizar esta acción.";
  }

  if (status === 404) {
    return "No encontramos la información solicitada.";
  }

  if (status === 409) {
    return "No pudimos completar la operación con la información actual.";
  }

  if (status >= 500) {
    return "No pudimos cargar la información en este momento.";
  }

  return "No pudimos procesar la solicitud.";
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const {
    method = "GET",
    body,
    query,
    auth = true,
    headers,
    retry = method === "GET",
    timeoutMs = REQUEST_TIMEOUT_MS,
  } = options;
  const requestHeaders = new Headers(headers);
  const requestId = nextWakeRequestId++;
  const url = buildUrl(path, query);
  const retryDelays = retry ? RETRY_DELAYS_MS : [];

  if (body !== undefined) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (auth) {
    const token = getStoredToken();
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }
  }

  for (let attempt = 0; attempt <= retryDelays.length; attempt += 1) {
    let wakeTimer: number | null = null;

    try {
      if (retry) {
        wakeTimer = window.setTimeout(() => {
          markWakeRequest(requestId, attempt + 1);
        }, SERVER_WAKE_THRESHOLD_MS);
      }

      const response = await fetchWithTimeout(
        url,
        {
          method,
          headers: requestHeaders,
          body: body !== undefined ? JSON.stringify(body) : undefined,
        },
        timeoutMs,
      );

      if (wakeTimer !== null) {
        window.clearTimeout(wakeTimer);
      }

      if (!response.ok) {
        if (retry && shouldRetryStatus(response.status) && attempt < retryDelays.length) {
          markWakeRequest(requestId, attempt + 1);
          await wait(retryDelays[attempt]);
          continue;
        }

        clearWakeRequest(requestId);

        const payload = await readPayload<T>(response);
        const message =
          typeof payload === "object" && payload && "error" in payload
            ? payload.error
            : getStatusFallbackMessage(response.status);

        throw new ApiError(message, response.status, payload as ApiErrorResponse | string);
      }

      clearWakeRequest(requestId);
      return (await readPayload<T>(response)) as T;
    } catch (error) {
      if (wakeTimer !== null) {
        window.clearTimeout(wakeTimer);
      }

      if (error instanceof ApiError) {
        throw error;
      }

      if (retry && shouldRetryError(error) && attempt < retryDelays.length) {
        markWakeRequest(requestId, attempt + 1);
        await wait(retryDelays[attempt]);
        continue;
      }

      clearWakeRequest(requestId);

      const message =
        retry && shouldRetryError(error)
          ? "El servicio se está iniciando. Vamos a mostrar la información apenas esté disponible."
          : "No pudimos conectarnos con el servicio. Intentá nuevamente en unos segundos.";

      throw new ApiError(message, 0);
    }
  }

  clearWakeRequest(requestId);
  throw new ApiError("No pudimos cargar la información en este momento.", 0);
}
