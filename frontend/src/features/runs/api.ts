import { apiClient } from '@/lib/api';
import type { RunListResponse } from './types';

export function getAllRuns(page = 1, pageSize = 20, status?: string, model?: string, startDate?: string, endDate?: string) {
  const params = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize),
  });
  if (status) params.set('status', status);
  if (model) params.set('model', model);
  if (startDate) params.set('startDate', startDate);
  if (endDate) params.set('endDate', endDate);
  return apiClient<RunListResponse>(`/runs?${params.toString()}`);
}
