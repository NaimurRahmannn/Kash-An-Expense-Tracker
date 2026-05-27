import { API_BASE_URL } from "@/config/env";
import { getStoredUserID } from "@/lib/auth-storage";

type ApiRequestOptions = RequestInit & {
  headers?: HeadersInit;
};

function buildApiUrl(endpoint: string) {
  if (endpoint.startsWith("http")) {
    return endpoint;
  }

  return `${API_BASE_URL}${endpoint.startsWith("/") ? endpoint : `/${endpoint}`}`;
}

function getErrorMessage(data: unknown) {
  if (data && typeof data === "object" && "message" in data) {
    const message = (data as { message?: unknown }).message;

    if (typeof message === "string" && message.trim()) {
      return message;
    }
  }

  return "Request failed";
}

export async function apiRequest<T>(
  endpoint: string,
  options: ApiRequestOptions = {},
): Promise<T> {
  const headers = new Headers(options.headers);

  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(buildApiUrl(endpoint), {
    ...options,
    headers,
  });

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    throw new Error(getErrorMessage(data));
  }

  return data as T;
}

export function getAuthHeaders(): Record<string, string> {
  const userID = getStoredUserID();

  return userID ? { "X-User-ID": String(userID) } : {};
}
