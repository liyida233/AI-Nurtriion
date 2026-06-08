const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api";

export type ApiResult<T> = {
  data?: T;
  error?: string;
};

export async function api<T>(
  path: string,
  options: RequestInit = {},
  token?: string
): Promise<ApiResult<T>> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers
    }
  });

  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    return { error: payload.error ?? "Request failed" };
  }
  return { data: payload.data as T };
}

export async function downloadFile(path: string, token: string): Promise<ApiResult<Blob>> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {})
    }
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => ({}));
    return { error: payload.error ?? "Download failed" };
  }

  return { data: await response.blob() };
}
