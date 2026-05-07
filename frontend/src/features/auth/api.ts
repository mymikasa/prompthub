import { apiClient } from '@/lib/api';

export type LoginParams = {
  email: string;
  password: string;
  remember?: boolean;
};

export type SignupParams = {
  name: string;
  email: string;
  password: string;
};

export type User = {
  id: number;
  name: string;
  email: string;
  workspace_id?: number;
};

export function login(params: LoginParams) {
  return apiClient<User>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(params),
  });
}

export function signup(params: SignupParams) {
  return apiClient<User>('/auth/signup', {
    method: 'POST',
    body: JSON.stringify(params),
  });
}

export function forgotPassword(email: string) {
  return apiClient<{ message: string }>('/auth/forgot-password', {
    method: 'POST',
    body: JSON.stringify({ email }),
  });
}

export function getMe() {
  return apiClient<User>('/me');
}

export function logout() {
  return apiClient<void>('/auth/logout', { method: 'POST' });
}
