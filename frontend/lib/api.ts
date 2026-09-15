export interface ApiResponse<T = any> {
  ok: boolean;
  status: number;
  statusText: string;
  durationMs: number;
  data: T;
  rawText: string;
  endpoint: string;
  method: string;
}

export function getApiBaseUrl(): string {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem("semard_api_url");
    if (saved && saved.trim() !== "") return saved.trim().replace(/\/$/, "");
  }
  return (
    process.env.NEXT_PUBLIC_API_URL ||
    "https://semardcontrolcenter.147.5.103.87.sslip.io"
  ).replace(/\/$/, "");
}

export function setApiBaseUrl(url: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("semard_api_url", url);
  }
}

export function getAuthToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("semard_jwt");
}

export function setAuthToken(token: string) {
  if (typeof window !== "undefined") {
    localStorage.setItem("semard_jwt", token);
  }
}

export function removeAuthToken() {
  if (typeof window !== "undefined") {
    localStorage.removeItem("semard_jwt");
  }
}

export async function apiCall<T = any>(
  endpoint: string,
  options: {
    method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
    body?: any;
    requireAuth?: boolean;
    isFormData?: boolean;
  } = {}
): Promise<ApiResponse<T>> {
  const { method = "GET", body, requireAuth = false, isFormData = false } = options;
  const baseUrl = getApiBaseUrl();
  const url = endpoint.startsWith("http") ? endpoint : `${baseUrl}${endpoint}`;

  const headers: Record<string, string> = {};

  if (requireAuth) {
    const token = getAuthToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  if (!isFormData && body) {
    headers["Content-Type"] = "application/json";
  }

  const startTime = performance.now();

  try {
    const res = await fetch(url, {
      method,
      headers,
      body: isFormData ? body : body ? JSON.stringify(body) : undefined,
    });

    const durationMs = Math.round(performance.now() - startTime);
    const rawText = await res.text();

    let data: any;
    try {
      data = JSON.parse(rawText);
    } catch {
      data = rawText;
    }

    return {
      ok: res.ok,
      status: res.status,
      statusText: res.statusText,
      durationMs,
      data,
      rawText,
      endpoint,
      method,
    };
  } catch (err: any) {
    const durationMs = Math.round(performance.now() - startTime);
    return {
      ok: false,
      status: 0,
      statusText: "Network Error",
      durationMs,
      data: { error: err?.message || "Fallo de conexión" } as any,
      rawText: err?.message || "Fallo de conexión",
      endpoint,
      method,
    };
  }
}
