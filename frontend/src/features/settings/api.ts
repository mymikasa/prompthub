import { apiClient } from '@/lib/api';
import type { ProviderConfig, SaveProviderConfigRequest } from './types';

export function listProviders() {
  return apiClient<ProviderConfig[]>('/settings/providers');
}

export function saveProviderConfig(data: SaveProviderConfigRequest) {
  return apiClient<ProviderConfig>('/settings/providers', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function deleteProviderConfig(id: number) {
  return apiClient<void>(`/settings/providers/${id}`, {
    method: 'DELETE',
  });
}
