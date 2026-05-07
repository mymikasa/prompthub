export type ApiError = {
  code: number;
  message: string;
  details?: unknown;
};

const API_BASE = '/api';

async function parseError(response: Response): Promise<ApiError> {
  try {
    const body = await response.json();
    return {
      code: body.code || 0,
      message: body.message || '请求失败',
      details: body.details,
    };
  } catch {
    return {
      code: response.status * 100,
      message: `请求失败 (${response.status})`,
    };
  }
}

export async function apiClient<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',
  });

  if (!res.ok) {
    const err = await parseError(res);
    throw err;
  }

  if (res.status === 204) return undefined as T;

  const body = await res.json();
  return body.data as T;
}
